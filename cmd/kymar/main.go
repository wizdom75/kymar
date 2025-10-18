package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"

	"github.com/pn/kymar/internal/db"
	"github.com/pn/kymar/internal/ui"
)

func main() {
	a := app.New()
	a.SetIcon(nil) // You can add a custom icon later

	// Set dark theme
	a.Settings().SetTheme(&ui.CustomDarkTheme{})

	w := a.NewWindow("Kymar - Database Client Pro")
	w.Resize(fyne.NewSize(1600, 900)) // Larger initial size
	w.CenterOnScreen()

	// Connection handler - declare as var first to allow recursive reference
	var handleConnection func(params db.ConnParams)
	handleConnection = func(params db.ConnParams) {
		// Try to connect
		dbh, closer, err := db.Connect(params)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		// Connection successful, show main interface
		ui.ShowMainInterface(w, dbh, closer, params, func() {
			// onDisconnect callback
			ui.ShowLoginScreen(w, handleConnection)
		})
	}

	// Show login screen first
	ui.ShowLoginScreen(w, handleConnection)

	w.ShowAndRun()
}
