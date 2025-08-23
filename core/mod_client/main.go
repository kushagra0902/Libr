package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed frontend/dist/*
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:            "libr",
		Width:            1024,
		Height:           768,
		AssetServer:      &assetserver.Options{Assets: assets},
		OnStartup:        app.startup,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		Bind:             []interface{}{app},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
