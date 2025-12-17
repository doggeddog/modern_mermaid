package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx          context.Context
	autoSavePath string
	configPath   string

	// Navigation State
	diagrams     []string
	currentIndex int

	// Zoom State
	zoomLevel float64

	// Menu References
	statusItem *menu.MenuItem
	prevItem   *menu.MenuItem
	nextItem   *menu.MenuItem
}

// AppConfig stores persistent configuration
type AppConfig struct {
	ZoomLevel float64 `json:"zoomLevel"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	// 获取用户配置目录
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir, _ = os.UserHomeDir()
	}
	appDir := filepath.Join(configDir, "modern-mermaid")
	_ = os.MkdirAll(appDir, 0755)

	a := &App{
		autoSavePath: filepath.Join(appDir, "autosave.mmd"),
		configPath:   filepath.Join(appDir, "config.json"),
		diagrams:     []string{},
		currentIndex: 0,
		zoomLevel:    1.0,
	}

	// Load config on init
	a.loadConfig()

	return a
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 监听前端的代码变更事件
	runtime.EventsOn(ctx, "codeChanged", func(optionalData ...interface{}) {
		if len(optionalData) > 0 {
			content, ok := optionalData[0].(string)
			if ok {
				a.saveToDisk(content)
				// 同步更新当前列表中的图表代码
				if len(a.diagrams) > 0 && a.currentIndex >= 0 && a.currentIndex < len(a.diagrams) {
					a.diagrams[a.currentIndex] = content
				}
			}
		}
	})

	// 监听前端就绪事件
	runtime.EventsOn(ctx, "onAppReady", func(optionalData ...interface{}) {
		a.loadFromDisk()
	})
	
	// 初始化菜单状态
	a.updateMenuState()

	// Apply saved zoom level
	a.applyZoom()
}

// domReady is called after the front-end dom has been loaded
func (a *App) domReady(ctx context.Context) {
	// 注入桥接脚本
	script := `
    (function() {
        console.log("Wails Bridge Injected");
        
        // 1. 监听 textarea 输入
        function attachListener() {
            const textarea = document.querySelector('textarea');
            if (textarea) {
                console.log("Textarea found, attaching listener");
                
                // 监听输入并发送给后端
                textarea.addEventListener('input', (e) => {
                    window.runtime.EventsEmit("codeChanged", e.target.value);
                });

                // 告诉后端我们准备好了
                window.runtime.EventsEmit("onAppReady");
                return true;
            }
            return false;
        }

        // 使用 MutationObserver 等待 textarea 出现
        const observer = new MutationObserver((mutations) => {
            if (attachListener()) {
                observer.disconnect();
            }
        });
        
        observer.observe(document.body, { childList: true, subtree: true });
        
        // 尝试立即绑定
        attachListener();

        // 2. 监听来自后端的加载事件
        window.runtime.EventsOn("loadFileContent", (content) => {
            console.log("Received file content");
            const textarea = document.querySelector('textarea');
            if (textarea) {
                const nativeInputValueSetter = Object.getOwnPropertyDescriptor(window.HTMLTextAreaElement.prototype, "value").set;
                nativeInputValueSetter.call(textarea, content);
                textarea.dispatchEvent(new Event('input', { bubbles: true }));
            }
        });

        // 3. 监听粘贴内容事件
        window.runtime.EventsOn("pasteContent", (content) => {
            const textarea = document.querySelector('textarea');
            if (textarea) {
                const start = textarea.selectionStart || textarea.value.length;
                const end = textarea.selectionEnd || textarea.value.length;
                const value = textarea.value;
                const newValue = value.substring(0, start) + content + value.substring(end);
                
                const nativeInputValueSetter = Object.getOwnPropertyDescriptor(window.HTMLTextAreaElement.prototype, "value").set;
                nativeInputValueSetter.call(textarea, newValue);
                textarea.dispatchEvent(new Event('input', { bubbles: true }));
            }
        });
    })();
    `
	runtime.WindowExecJS(ctx, script)
}

func (a *App) saveToDisk(content string) {
	// 简单的直接写入（操作系统通常处理得很快）
	err := os.WriteFile(a.autoSavePath, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error saving file: %v\n", err)
	}
}

func (a *App) loadFromDisk() {
	content, err := os.ReadFile(a.autoSavePath)
	if err == nil {
		runtime.EventsEmit(a.ctx, "loadFileContent", string(content))
	}
}

func (a *App) loadConfig() {
	data, err := os.ReadFile(a.configPath)
	if err == nil {
		var cfg AppConfig
		if json.Unmarshal(data, &cfg) == nil {
			if cfg.ZoomLevel > 0 {
				a.zoomLevel = cfg.ZoomLevel
			}
		}
	}
}

func (a *App) saveConfig() {
	cfg := AppConfig{ZoomLevel: a.zoomLevel}
	data, _ := json.MarshalIndent(cfg, "", "  ")
	os.WriteFile(a.configPath, data, 0644)
}

// SetMenuRefs stores menu item references
func (a *App) SetMenuRefs(status, prev, next *menu.MenuItem) {
	a.statusItem = status
	a.prevItem = prev
	a.nextItem = next
	// 初始化时隐藏
	if a.statusItem != nil { a.statusItem.Hidden = true }
	if a.prevItem != nil { a.prevItem.Hidden = true }
	if a.nextItem != nil { a.nextItem.Hidden = true }
}

func (a *App) updateMenuState() {
	if a.ctx == nil {
		return
	}

	count := len(a.diagrams)
	if count <= 1 {
		// Hide nav items if only 0 or 1 diagram
		if a.statusItem != nil {
			a.statusItem.Hidden = true
		}
		if a.prevItem != nil {
			a.prevItem.Hidden = true
		}
		if a.nextItem != nil {
			a.nextItem.Hidden = true
		}
	} else {
		// Show items
		if a.statusItem != nil {
			a.statusItem.Hidden = false
			a.statusItem.Label = fmt.Sprintf("Diagram: %d / %d", a.currentIndex+1, count)
		}
		if a.prevItem != nil {
			a.prevItem.Hidden = false
			a.prevItem.Disabled = a.currentIndex <= 0
		}
		if a.nextItem != nil {
			a.nextItem.Hidden = false
			a.nextItem.Disabled = a.currentIndex >= count-1
		}
	}
	runtime.MenuUpdateApplicationMenu(a.ctx)
}

func (a *App) loadCurrentDiagram() {
	if len(a.diagrams) > 0 && a.currentIndex >= 0 && a.currentIndex < len(a.diagrams) {
		code := a.diagrams[a.currentIndex]
		runtime.EventsEmit(a.ctx, "loadFileContent", code)
	}
}

// ImportFromClipboard reads clipboard, extracts mermaid blocks, and loads first one
func (a *App) ImportFromClipboard() {
	text, err := runtime.ClipboardGetText(a.ctx)
	if err != nil || text == "" {
		return
	}

	// Regex extract ```mermaid ... ```
	// (?s) allows . to match newlines
	re := regexp.MustCompile("(?s)```mermaid\\s*(.*?)```")
	matches := re.FindAllStringSubmatch(text, -1)

	var blocks []string
	for _, m := range matches {
		if len(m) > 1 {
			code := strings.TrimSpace(m[1])
			if code != "" {
				blocks = append(blocks, code)
			}
		}
	}

	// If no blocks found, treat the whole text as content?
	// The requirement is to extract. If none, maybe just paste as is or do nothing?
	// Let's default to: if no blocks found, assume it is NOT markdown with blocks,
	// but just code.
	if len(blocks) == 0 {
		blocks = []string{text}
	}

	a.diagrams = blocks
	a.currentIndex = 0
	a.loadCurrentDiagram()
	a.updateMenuState()
}

// NextPreview shows the next diagram
func (a *App) NextPreview() {
	if a.currentIndex < len(a.diagrams)-1 {
		a.currentIndex++
		a.loadCurrentDiagram()
		a.updateMenuState()
	}
}

// PrevPreview shows the previous diagram
func (a *App) PrevPreview() {
	if a.currentIndex > 0 {
		a.currentIndex--
		a.loadCurrentDiagram()
		a.updateMenuState()
	}
}

// OpenFileDialog prompts user to select a file and loads it
func (a *App) OpenFileDialog() {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Open Mermaid File",
		Filters: []runtime.FileFilter{
			{DisplayName: "Mermaid Files (*.mmd;*.mermaid;*.txt;*.md)", Pattern: "*.mmd;*.mermaid;*.txt;*.md"},
		},
	})

	if err == nil && selection != "" {
		content, err := os.ReadFile(selection)
		if err == nil {
			// Also support extraction from file
			text := string(content)
			// Check if file extension is .md
			if strings.HasSuffix(strings.ToLower(selection), ".md") {
				// Parse blocks
				re := regexp.MustCompile("(?s)```mermaid\\s*(.*?)```")
				matches := re.FindAllStringSubmatch(text, -1)
				var blocks []string
				for _, m := range matches {
					if len(m) > 1 {
						code := strings.TrimSpace(m[1])
						if code != "" {
							blocks = append(blocks, code)
						}
					}
				}
				if len(blocks) > 0 {
					a.diagrams = blocks
					a.currentIndex = 0
					a.loadCurrentDiagram()
					a.updateMenuState()
					return
				}
			}
			
			// Fallback normal load
			a.diagrams = []string{} // Reset multi-view
			a.updateMenuState()
			runtime.EventsEmit(a.ctx, "loadFileContent", text)
			a.saveToDisk(text)
		}
	}
}

// SaveImage prompts user to save an image
func (a *App) SaveImage(base64Data string, filename string) error {
	// Remove header if present
	if idx := strings.Index(base64Data, ","); idx != -1 {
		base64Data = base64Data[idx+1:]
	}

	// Decode
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return err
	}

	// Save Dialog
	selection, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Save Image",
		DefaultFilename: filename,
		Filters: []runtime.FileFilter{
			{DisplayName: "Image Files", Pattern: "*.png;*.jpg"},
		},
	})

	if err != nil || selection == "" {
		return nil
	}

	return os.WriteFile(selection, data, 0644)
}

// PasteFromClipboard reads clipboard and sends to frontend
func (a *App) PasteFromClipboard() {
	text, err := runtime.ClipboardGetText(a.ctx)
	if err == nil && text != "" {
		runtime.EventsEmit(a.ctx, "pasteContent", text)
	}
}

// ZoomIn increases zoom level
func (a *App) ZoomIn() {
	a.zoomLevel += 0.1
	a.applyZoom()
	a.saveConfig()
}

// ZoomOut decreases zoom level
func (a *App) ZoomOut() {
	if a.zoomLevel > 0.2 {
		a.zoomLevel -= 0.1
	}
	a.applyZoom()
	a.saveConfig()
}

// ZoomReset resets zoom level
func (a *App) ZoomReset() {
	a.zoomLevel = 1.0
	a.applyZoom()
	a.saveConfig()
}

func (a *App) applyZoom() {
	// Apply zoom to document body
	// Note: We use %f to format float, effectively creating "1.100000" string
	js := fmt.Sprintf("document.body.style.zoom = '%f'", a.zoomLevel)
	runtime.WindowExecJS(a.ctx, js)
}

// Quit the application
func (a *App) Quit() {
	runtime.Quit(a.ctx)
}
