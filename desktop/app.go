package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx          context.Context
	autoSavePath string
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

	return &App{
		autoSavePath: filepath.Join(appDir, "autosave.mmd"),
	}
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
			}
		}
	})

	// 监听前端就绪事件
	runtime.EventsOn(ctx, "onAppReady", func(optionalData ...interface{}) {
		a.loadFromDisk()
	})
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

// OpenFileDialog prompts user to select a file and loads it
func (a *App) OpenFileDialog() {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Open Mermaid File",
		Filters: []runtime.FileFilter{
			{DisplayName: "Mermaid Files (*.mmd;*.mermaid;*.txt)", Pattern: "*.mmd;*.mermaid;*.txt"},
		},
	})

	if err == nil && selection != "" {
		content, err := os.ReadFile(selection)
		if err == nil {
			runtime.EventsEmit(a.ctx, "loadFileContent", string(content))
			// 更新自动保存状态
			a.saveToDisk(string(content))
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

// Quit the application
func (a *App) Quit() {
	runtime.Quit(a.ctx)
}
