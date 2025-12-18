package main

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:assets
var assets embed.FS

// IndexInjectionHandler intercepts index.html to inject startup configuration
type IndexInjectionHandler struct {
	assets fs.FS
	app    *App
	next   http.Handler
}

func (h *IndexInjectionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Debug log to confirm handler is called
	println("[Wails Handler] Request:", r.URL.Path)

	// Intercept root or index.html requests
	if r.URL.Path == "/" || r.URL.Path == "/index.html" {
		// Read original file from embed fs
		// Note: "assets" dir is the root of our embed.FS because of //go:embed all:assets
		// Wails assetserver normally expects the files to be at root of provided FS if configured so.
		// Let's try to read "assets/index.html" first.

		content, err := fs.ReadFile(h.assets, "assets/index.html")
		if err != nil {
			// Try without "assets/" prefix just in case structure differs
			content, err = fs.ReadFile(h.assets, "index.html")
			if err != nil {
				println("[Wails Handler] Error reading index.html:", err.Error())
				if h.next != nil {
					h.next.ServeHTTP(w, r)
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
				return
			}
		}

		// Inject script into <head>
		html := string(content)

		// Disable Service Worker in Wails to prevent "protocol not supported" errors
		// And set flag for debugging
		swDisable := `<script>
		window.isMiddlewareInjected = true;
		if(window.navigator) {
			window.navigator.serviceWorker = { 
				register: function() { return Promise.resolve({}); }, 
				getRegistrations: function() { return Promise.resolve([]); },
				ready: Promise.resolve({})
			};
		}
		</script>`

		injection := swDisable + h.app.GetStartupInjectionScript()

		// Look for <head> tag
		if strings.Contains(html, "<head>") {
			html = strings.Replace(html, "<head>", "<head>"+injection, 1)
		} else {
			// Fallback: prepend to body or just valid HTML
			html = injection + html
		}

		w.Header().Set("Content-Type", "text/html")
		// Disable caching for index.html to ensure injection script is always fresh
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Write([]byte(html))
		println("[Wails Handler] Injected startup script into index.html")
		return
	}

	// For all other files, return 404 to let Wails internal asset server handle it
	if h.next != nil {
		h.next.ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	// macOS Fix: When launched from Finder, LANG is often empty, causing CGO/Clipboard encoding issues.
	// Force set it to en_US.UTF-8 to ensure UTF-8 support.
	if runtime.GOOS == "darwin" && os.Getenv("LANG") == "" {
		os.Setenv("LANG", "en_US.UTF-8")
	}

	// Create an instance of the app structure
	app := NewApp()

	// Create application menu
	appMenu := menu.NewMenu()

	// File Menu
	FileMenu := appMenu.AddSubmenu("File")
	FileMenu.AddText("Open...", keys.CmdOrCtrl("o"), func(_ *menu.CallbackData) {
		app.OpenFileDialog()
	})
	FileMenu.AddText("Import from Clipboard", keys.CmdOrCtrl("i"), func(_ *menu.CallbackData) {
		app.ImportFromClipboard()
	})
	FileMenu.AddSeparator()
	FileMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		app.Quit()
	})

	// Edit Menu
	// On macOS, use the standard Edit menu. On others, build it manually.
	if runtime.GOOS == "darwin" {
		appMenu.Append(menu.EditMenu())
	} else {
		EditMenu := appMenu.AddSubmenu("Edit")
		EditMenu.AddText("Undo", keys.CmdOrCtrl("z"), func(_ *menu.CallbackData) {
			// Native webview undo
		})
		EditMenu.AddText("Redo", keys.CmdOrCtrl("y"), func(_ *menu.CallbackData) {
			// Native webview redo
		})
		EditMenu.AddSeparator()
		EditMenu.AddText("Cut", keys.CmdOrCtrl("x"), func(_ *menu.CallbackData) {
			// Native webview cut
		})
		EditMenu.AddText("Copy", keys.CmdOrCtrl("c"), func(_ *menu.CallbackData) {
			// Native webview copy
		})
		EditMenu.AddText("Paste", keys.CmdOrCtrl("v"), func(_ *menu.CallbackData) {
			app.PasteFromClipboard()
		})
	}

	// View Menu
	ViewMenu := appMenu.AddSubmenu("View")
	ViewMenu.AddText("Zoom In", keys.CmdOrCtrl("+"), func(_ *menu.CallbackData) {
		app.ZoomIn()
	})
	ViewMenu.AddText("Zoom Out", keys.CmdOrCtrl("-"), func(_ *menu.CallbackData) {
		app.ZoomOut()
	})
	ViewMenu.AddText("Reset Zoom", keys.CmdOrCtrl("0"), func(_ *menu.CallbackData) {
		app.ZoomReset()
	})
	ViewMenu.AddSeparator()
	ViewMenu.AddText("Toggle Header", nil, func(_ *menu.CallbackData) {
		app.ToggleHeader()
	})

	// Navigation Menu (Top Level Item showing Status)
	// We manually construct the item so we can reference it to change the label
	NavSubMenu := menu.NewMenu()
	NavTopLevelItem := appMenu.AddText("No Diagrams", nil, nil)
	NavTopLevelItem.SubMenu = NavSubMenu
	// Initially hidden until we have diagrams
	NavTopLevelItem.Hidden = true

	PrevItem := NavSubMenu.AddText("Previous Diagram", keys.CmdOrCtrl("["), func(_ *menu.CallbackData) {
		app.PrevPreview()
	})

	NextItem := NavSubMenu.AddText("Next Diagram", keys.CmdOrCtrl("]"), func(_ *menu.CallbackData) {
		app.NextPreview()
	})

	// Pass menu references to App
	app.SetMenuRefs(NavTopLevelItem, PrevItem, NextItem)

	// Create application with options
	err := wails.Run(&options.App{
		Title:            "Modern Mermaid",
		Width:            1200,
		Height:           800,
		WindowStartState: options.Maximised,
		AssetServer: &assetserver.Options{
			Assets: assets,
			Middleware: func(next http.Handler) http.Handler {
				return &IndexInjectionHandler{
					assets: assets,
					app:    app,
					next:   next,
				}
			},
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        app.startup,
		OnDomReady:       app.domReady,
		Menu:             appMenu,
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			// TitleBar: mac.TitleBarHiddenInset(), // Removed to use standard title bar
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
