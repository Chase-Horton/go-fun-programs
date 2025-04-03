package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Search(app *tview.Application) (*tview.Flex, tview.Primitive) {

	searchResults := tview.NewTable()
	searchResultsBox := searchResults.SetBorder(true)

	searchInput := tview.NewInputField().SetLabel("Search: ").SetFieldWidth(0).SetFieldBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	searchInput.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			searchResultsBox.SetTitle("Search Results for " + searchInput.GetText())
			app.SetFocus(searchResults)
		default:
			return
		}
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchInput, 1, 0, false).
		AddItem(searchResults, 0, 1, true)

	return flex, searchInput
}

func MainView() (*tview.Flex, tview.Primitive) {
	// Create a new table to display the movie list
	movieList := tview.NewTable().
		SetBorder(true).
		SetTitle("Movie List")

	movieDescription := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetBorder(true).
		SetTitle("Movie Description")

	flex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(movieList, 0, 1, true).
		AddItem(movieDescription, 0, 1, false)

	return flex, movieList
}
