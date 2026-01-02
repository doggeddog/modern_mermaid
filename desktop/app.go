package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var (
	AppVersion = "0.0.0"
	BuildTime  = "dev"
)

// App struct
type App struct {
	ctx          context.Context
	autoSavePath string
	configPath   string
	db           *sql.DB

	// Navigation State
	diagrams     []Diagram
	currentIndex int

	// Zoom State
	zoomLevel float64

	// UI State
	headerVisible bool

	// Persistent Frontend State
	theme      string
	background string
	font       string
	language   string
	darkMode   bool

	// Menu References
	statusItem  *menu.MenuItem
	prevItem    *menu.MenuItem
	nextItem    *menu.MenuItem
	historyMenu *menu.Menu

	// Flags
	isNewWindow bool
}

// AppConfig stores persistent configuration
type AppConfig struct {
	ZoomLevel     float64 `json:"zoomLevel"`
	HeaderVisible *bool   `json:"headerVisible"`
	Theme         string  `json:"theme"`
	Background    string  `json:"background"`
	Font          string  `json:"font"`
	Language      string  `json:"language"`
	DarkMode      *bool   `json:"darkMode"`
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

	dbPath := filepath.Join(appDir, "modern-mermaid.db")
	db, err := initDB(dbPath)
	if err != nil {
		fmt.Printf("Failed to init DB: %v\n", err)
	}

	a := &App{
		autoSavePath:  filepath.Join(appDir, "autosave.mmd"),
		configPath:    filepath.Join(appDir, "config.json"),
		db:            db,
		diagrams:      []Diagram{},
		currentIndex:  0,
		zoomLevel:     1.0,
		headerVisible: false,         // Default hidden
		theme:         "linearLight", // Default theme
		background:    "dots",
		font:          "inter",
		language:      "zh-CN", // Default language (changed from en)
		darkMode:      false,
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
				// 同步更新当前列表中的图表代码
				if len(a.diagrams) > 0 && a.currentIndex >= 0 && a.currentIndex < len(a.diagrams) {
					currentID := a.diagrams[a.currentIndex].ID
					a.dbUpdateDiagram(currentID, content)
					a.diagrams[a.currentIndex].Content = content
				} else {
					// Fallback: If no diagram loaded (weird state), insert one
					id, _ := a.dbInsertDiagram(content, "edit")
					a.diagrams = []Diagram{{ID: id, Content: content}}
					a.currentIndex = 0
				}
			}
		}
	})

	// Listen for config changes from Frontend
	runtime.EventsOn(ctx, "configChanged", func(data ...interface{}) {
		if len(data) > 0 {
			// Log for debugging
			fmt.Printf("[Wails] ConfigChanged Event: %+v\n", data[0])

			// Expecting a JSON string or map
			// Wails usually passes map[string]interface{} for JSON objects
			if configMap, ok := data[0].(map[string]interface{}); ok {
				a.updateConfigFromMap(configMap)
			} else {
				fmt.Println("[Wails] Error: Config data is not a map")
			}
		}
	})

	// 监听前端就绪事件
	runtime.EventsOn(ctx, "onAppReady", func(optionalData ...interface{}) {
		if a.isNewWindow {
			a.loadNewDiagram()
		} else {
			a.loadLatestFromDB()
		}
		a.ImportFromClipboard() // This calls rebuildHistoryMenu inside if it imports something, but maybe not if it doesn't.
		// To be safe, call it always.
		a.rebuildHistoryMenu()
		// Apply saved state
		a.applyZoom()
		a.applyHeaderVisibility()
	})

	// 初始化菜单状态
	a.updateMenuState()
}

// GetStartupInjectionScript returns the script to be injected into index.html head
// This ensures URL params and localStorage are set BEFORE React loads.
func (a *App) GetStartupInjectionScript() string {
	queryParams := []string{}
	if a.theme != "" {
		queryParams = append(queryParams, fmt.Sprintf("theme=%s", a.theme))
	}
	if a.background != "" {
		queryParams = append(queryParams, fmt.Sprintf("bg=%s", a.background))
	}
	if a.font != "" {
		queryParams = append(queryParams, fmt.Sprintf("font=%s", a.font))
	}
	if a.language != "" {
		queryParams = append(queryParams, fmt.Sprintf("lang=%s", a.language))
	}

	queryString := strings.Join(queryParams, "&")

	// Note: We use window.history.replaceState immediately.
	// We also set localStorage for darkMode.

	script := fmt.Sprintf(`
    <script>
    (function() {
        try {
            // 1. Restore Dark Mode
            const savedDark = '%v';
            if (savedDark === 'true') {
                localStorage.setItem('darkMode', 'true');
                document.documentElement.classList.add('dark');
            } else {
                localStorage.setItem('darkMode', 'false');
                document.documentElement.classList.remove('dark');
            }

            // 2. Restore URL Params
            // Only if current URL has no params (fresh start)
            if (!window.location.search && '%s' !== "") {
                 const newUrl = window.location.pathname + '?' + '%s';
                 window.history.replaceState({}, '', newUrl);
            }
        } catch(e) { console.error("Wails Init Error:", e); }
    })();
    </script>
    `, a.darkMode, queryString, queryString)

	return script
}

// domReady is called after the front-end dom has been loaded
func (a *App) domReady(ctx context.Context) {
	// Only inject the bridge and watchers here.
	// The state restoration is now handled by Index Injection in main.go

	script := `
    (function() {
        console.log("Wails Bridge Injected");

        // --- Watchers for Persistence ---

        // Watch URL changes (Theme, Bg, Font, Lang)
        // Method 1: Monkey Patch History API
        const originalPushState = history.pushState;
        history.pushState = function() {
            console.log("[Wails] history.pushState called. Args:", arguments);
            try {
                originalPushState.apply(this, arguments);
            } catch (e) {
                console.error("[Wails] pushState failed:", e);
            }
            notifyConfigChangeFromURL();
        };

        const originalReplaceState = history.replaceState;
        history.replaceState = function() {
            console.log("[Wails] history.replaceState called. Args:", arguments);
            try {
                originalReplaceState.apply(this, arguments);
            } catch (e) {
                console.error("[Wails] replaceState failed:", e);
            }
            notifyConfigChangeFromURL();
        };

        // Method 2: Polling (Fallback for reliability)
        let lastSearch = window.location.search;
        console.log("Polling started. Initial search:", lastSearch);
        console.log("Middleware Injected Check:", window.isMiddlewareInjected);
        
        setInterval(() => {
            // Check both search and href to be sure
            if (window.location.search !== lastSearch) {
                console.log("URL Changed (polling). New Search:", window.location.search);
                lastSearch = window.location.search;
                notifyConfigChangeFromURL();
            }
        }, 1000);

        function notifyConfigChangeFromURL() {
            console.log("Notifying Config Change... Current URL:", window.location.href);
            // Debounce or immediate? Immediate is fine for now.
            setTimeout(() => {
                const params = new URLSearchParams(window.location.search);
                const config = {
                    theme: params.get('theme') || '',
                    background: params.get('bg') || '',
                    font: params.get('font') || '',
                    language: params.get('lang') || ''
                };
                console.log("Sending Config Update:", config);
                window.runtime.EventsEmit("configChanged", config);
            }, 100);
        }
        
        // Watch LocalStorage (DarkMode)
        const originalSetItem = localStorage.setItem;
        localStorage.setItem = function(key, value) {
            originalSetItem.apply(this, arguments);
            if (key === 'darkMode') {
                window.runtime.EventsEmit("configChanged", { darkMode: value === 'true' });
            }
        };


        // --- Editor Interaction ---
        
        function attachListener() {
            const textarea = document.querySelector('textarea');
            if (textarea) {
                textarea.addEventListener('input', (e) => {
                    window.runtime.EventsEmit("codeChanged", e.target.value);
                });
                window.runtime.EventsEmit("onAppReady");
                return true;
            }
            return false;
        }

        const observer = new MutationObserver((mutations) => {
            if (attachListener()) {
                observer.disconnect();
            }
        });
        
        observer.observe(document.body, { childList: true, subtree: true });
        
        attachListener();

        window.runtime.EventsOn("loadFileContent", (content) => {
            const textarea = document.querySelector('textarea');
            if (textarea) {
                const nativeInputValueSetter = Object.getOwnPropertyDescriptor(window.HTMLTextAreaElement.prototype, "value").set;
                nativeInputValueSetter.call(textarea, content);
                textarea.dispatchEvent(new Event('input', { bubbles: true }));
            }
        });

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

func (a *App) updateConfigFromMap(data map[string]interface{}) {
	changed := false

	if val, ok := data["theme"].(string); ok && val != "" {
		if a.theme != val {
			a.theme = val
			changed = true
		}
	}
	if val, ok := data["background"].(string); ok && val != "" {
		if a.background != val {
			a.background = val
			changed = true
		}
	}
	if val, ok := data["font"].(string); ok && val != "" {
		if a.font != val {
			a.font = val
			changed = true
		}
	}
	if val, ok := data["language"].(string); ok && val != "" {
		if a.language != val {
			a.language = val
			changed = true
		}
	}
	if val, ok := data["darkMode"].(bool); ok {
		if a.darkMode != val {
			a.darkMode = val
			changed = true
		}
	}

	if changed {
		a.saveConfig()
	}
}

func (a *App) loadLatestFromDB() {
	d, err := a.dbGetLatestDiagram()
	if err == nil && d != nil {
		a.diagrams = []Diagram{*d}
		a.currentIndex = 0
		runtime.EventsEmit(a.ctx, "loadFileContent", d.Content)
	}
}

func (a *App) loadNewDiagram() {
	defaultContent := "graph TD\n    A[Start] --> B[End]"
	id, err := a.dbInsertDiagram(defaultContent, "new")
	if err == nil {
		a.diagrams = []Diagram{{ID: id, Content: defaultContent, Source: "new", Title: "Untitled"}}
		a.currentIndex = 0
		runtime.EventsEmit(a.ctx, "loadFileContent", defaultContent)
	}
}

func (a *App) NewWindow() {
	ex, err := os.Executable()
	if err != nil {
		fmt.Printf("Error getting executable: %v\n", err)
		return
	}
	fmt.Printf("Launching new window: %s --new-window\n", ex)
	// 在 macOS 上，如果直接运行 .app 包中的二进制文件，它可能不会作为独立的 App 实例运行，
	// 或者需要 `open -n -a "App Name" --args ...`
	// 但是 os.Executable() 指向的是 MacOS/Binary。
	// 简单尝试直接执行二进制文件。
	cmd := exec.Command(ex, "--new-window")
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Error launching new window: %v\n", err)
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
			if cfg.HeaderVisible != nil {
				a.headerVisible = *cfg.HeaderVisible
			}
			if cfg.Theme != "" {
				a.theme = cfg.Theme
			}
			if cfg.Background != "" {
				a.background = cfg.Background
			}
			if cfg.Font != "" {
				a.font = cfg.Font
			}
			if cfg.Language != "" {
				a.language = cfg.Language
			}
			if cfg.DarkMode != nil {
				a.darkMode = *cfg.DarkMode
			}
		}
	}
}

func (a *App) saveConfig() {
	cfg := AppConfig{
		ZoomLevel:     a.zoomLevel,
		HeaderVisible: &a.headerVisible,
		Theme:         a.theme,
		Background:    a.background,
		Font:          a.font,
		Language:      a.language,
		DarkMode:      &a.darkMode,
	}
	fmt.Printf("[Wails] Saving Config: %+v\n", cfg)
	data, _ := json.MarshalIndent(cfg, "", "  ")
	err := os.WriteFile(a.configPath, data, 0644)
	if err != nil {
		fmt.Printf("[Wails] Error saving config: %v\n", err)
	}
}

// SetMenuRefs stores menu item references
func (a *App) SetMenuRefs(status, prev, next *menu.MenuItem) {
	a.statusItem = status
	a.prevItem = prev
	a.nextItem = next
	// 初始化时隐藏
	if a.statusItem != nil {
		a.statusItem.Hidden = true
	}
	if a.prevItem != nil {
		a.prevItem.Hidden = true
	}
	if a.nextItem != nil {
		a.nextItem.Hidden = true
	}
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
		code := a.diagrams[a.currentIndex].Content
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

	var newDiagrams []Diagram
	for _, m := range matches {
		if len(m) > 1 {
			code := strings.TrimSpace(m[1])
			if code != "" {
				// Normalize code to fix potential syntax issues (e.g. unquoted special chars)
				code = NormalizeMermaid(code)
				id, err := a.dbInsertDiagram(code, "clipboard")
				if err == nil {
					newDiagrams = append(newDiagrams, Diagram{
						ID:      id,
						Content: code,
						Source:  "clipboard",
						Title:   extractTitle(code),
					})
				}
			}
		}
	}

	// Fallback logic:
	// If no markdown blocks found, check if the text itself looks like mermaid code.
	if len(newDiagrams) == 0 {
		if a.isMermaidCode(text) {
			// Normalize code
			text = NormalizeMermaid(text)
			id, err := a.dbInsertDiagram(text, "clipboard")
			if err == nil {
				newDiagrams = append(newDiagrams, Diagram{
					ID:      id,
					Content: text,
					Source:  "clipboard",
					Title:   extractTitle(text),
				})
			}
		} else {
			// Not a valid mermaid code, ignore.
			return
		}
	}

	a.diagrams = newDiagrams
	a.currentIndex = 0
	a.loadCurrentDiagram()
	a.updateMenuState()
	a.rebuildHistoryMenu()
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
				var newDiagrams []Diagram
				for _, m := range matches {
					if len(m) > 1 {
						code := strings.TrimSpace(m[1])
						if code != "" {
							id, err := a.dbInsertDiagram(code, "file")
							if err == nil {
								newDiagrams = append(newDiagrams, Diagram{
									ID:      id,
									Content: code,
									Source:  "file",
									Title:   extractTitle(code),
								})
							}
						}
					}
				}
				if len(newDiagrams) > 0 {
					a.diagrams = newDiagrams
					a.currentIndex = 0
					a.loadCurrentDiagram()
					a.updateMenuState()
					a.rebuildHistoryMenu()
					return
				}
			}

			// Fallback normal load
			text = NormalizeMermaid(text)
			id, err := a.dbInsertDiagram(text, "file")
			if err == nil {
				a.diagrams = []Diagram{{
					ID:      id,
					Content: text,
					Source:  "file",
					Title:   extractTitle(text),
				}}
				a.currentIndex = 0
				a.loadCurrentDiagram()
				a.updateMenuState()
				a.rebuildHistoryMenu()
			}
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

// ToggleHeader toggles the visibility of the header
func (a *App) ToggleHeader() {
	a.headerVisible = !a.headerVisible
	fmt.Printf("ToggleHeader: %v\n", a.headerVisible)
	a.applyHeaderVisibility()
	a.saveConfig()
}

func (a *App) applyHeaderVisibility() {
	display := "flex"
	if !a.headerVisible {
		display = "none"
	}
	// Use setProperty with !important to override any other styles
	// Wrap in IIFE to avoid "duplicate variable" errors
	js := fmt.Sprintf(`(function(){
		const h = document.querySelector('header');
		console.log("Wails ToggleHeader: Setting display to", '%s'); 
		if(h) h.style.setProperty('display', '%s', 'important');
	})()`, display, display)
	runtime.WindowExecJS(a.ctx, js)
}

// Quit the application
func (a *App) Quit() {
	runtime.Quit(a.ctx)
}

// ShowVersion displays the application version
func (a *App) ShowVersion() {
	message := fmt.Sprintf("Modern Mermaid Desktop\nVersion: %s\nBuild Time: %s", AppVersion, BuildTime)
	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   "About Modern Mermaid",
		Message: message,
	})
}

var mermaidKeywords = []string{
	"flowchart", "graph", "sequenceDiagram", "classDiagram",
	"stateDiagram", "stateDiagram-v2", "erDiagram", "journey",
	"gantt", "pie", "requirementDiagram", "gitGraph",
	"C4Context", "mindmap", "timeline", "sankey-beta",
	"quadrantChart", "xychart-beta", "block-beta", "packet-beta",
	"kanban", "architecture",
}

func (a *App) isMermaidCode(text string) bool {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "%%") {
			continue
		}
		// Check if line starts with any keyword
		for _, keyword := range mermaidKeywords {
			if strings.HasPrefix(line, keyword) {
				return true
			}
		}
		// Found a non-comment line that doesn't start with a keyword
		return false
	}
	return false
}

func (a *App) SetHistoryMenuRef(m *menu.Menu) {
	a.historyMenu = m
}

func (a *App) rebuildHistoryMenu() {
	if a.historyMenu == nil {
		return
	}

	diagrams, err := a.dbListDiagrams(50)
	// Clear items
	a.historyMenu.Items = nil

	if err != nil || len(diagrams) == 0 {
		a.historyMenu.AddText("(No History)", nil, nil).Disabled = true
	} else {
		for _, d := range diagrams {
			// Copy ID for closure
			id := d.ID
			// Format: YYYY-MM-DD HH:MM Title
			timeStr := d.UpdatedAt.Format("2006-01-02 15:04")
			title := d.Title

			// Truncate title if too long
			if len(title) > 50 {
				title = title[:47] + "..."
			}

			label := fmt.Sprintf("%s %s", timeStr, title)

			a.historyMenu.AddText(label, nil, func(_ *menu.CallbackData) {
				a.loadDiagramFromHistory(id)
			})
		}

		a.historyMenu.AddSeparator()
		a.historyMenu.AddText("Clear All History", nil, func(_ *menu.CallbackData) {
			a.ClearAllHistory()
		})
	}

	if a.ctx != nil {
		runtime.MenuUpdateApplicationMenu(a.ctx)
	}
}

func (a *App) ClearAllHistory() {
	err := a.dbClearHistory()
	if err != nil {
		fmt.Printf("Error clearing history: %v\n", err)
	}
	a.rebuildHistoryMenu()
}

func (a *App) loadDiagramFromHistory(id int64) {
	d, err := a.dbGetDiagram(id)
	if err == nil {
		a.diagrams = []Diagram{*d}
		a.currentIndex = 0
		a.loadCurrentDiagram()
		a.updateMenuState()
		// Update history menu order immediately?
		// Since we just loaded it, its updated_at hasn't changed unless we edit.
	}
}
