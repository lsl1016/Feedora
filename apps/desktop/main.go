package main

import (
	"embed"
	"log"

	desktopapp "github.com/lsl1016/feedora/apps/desktop/internal/app"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	application := desktopapp.New()

	err := wails.Run(&options.App{
		Title:     "Feedora Desktop",
		Width:     1440,
		Height:    900,
		MinWidth:  1180,
		MinHeight: 720,
		AssetServer: &assetserver.Options{Assets: assets},
		OnStartup: application.Startup,
		Bind:      application.Bindings(),
	})
	if err != nil {
		log.Fatalf("start Feedora Desktop failed:%v", err)
	}
}
