package main

import (
	"embed"
	"os"
	"runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed assets
var assets embed.FS

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
