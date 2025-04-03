package main

import (
	"log"
	"movie-list/views"

	"github.com/gdamore/tcell/v2"
	"github.com/joho/godotenv"
	"github.com/rivo/tview"
)

var app = tview.NewApplication()
var pages = tview.NewPages()

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	viewPage, viewPageFocus := views.MainView()
	pages.AddPage("view", viewPage, true, true)
	searchPage, searchPageFocus := views.Search(app)
	pages.AddPage("search", searchPage, true, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlS:
			pages.SwitchToPage("search")
			app.SetFocus(searchPageFocus)
			return nil
		case tcell.KeyCtrlV:
			pages.SwitchToPage("view")
			app.SetFocus(viewPageFocus)
			return nil
		}
		return event
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		panic(err)
	}
}
