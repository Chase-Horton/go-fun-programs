package views

import (
	"movie-list/logger"
	"movie-list/models"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func resetSearchResults(searchResults *tview.Table) {
	searchResults.Clear()
	searchResults.SetCell(0, 0, tview.NewTableCell("Title").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	searchResults.SetCell(0, 1, tview.NewTableCell("Type").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	searchResults.SetCell(0, 2, tview.NewTableCell("Year").SetTextColor(tcell.ColorYellow).SetSelectable(false))
}
func Search(app *tview.Application, pages *tview.Pages) (*tview.Flex, tview.Primitive) {

	searchResults := tview.NewTable()
	searchResultsBox := searchResults.SetBorder(false)

	searchInput := tview.NewInputField().SetLabel("Search: ").SetFieldWidth(0).SetFieldBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	searchInput.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			return
		}

		searchResultsBox.SetTitle("Search Results for " + searchInput.GetText())
		data := models.SearchMovie(searchInput.GetText())
		if data.Response == "False" {
			return
		}

		resetSearchResults(searchResults)
		for row, movie := range data.Movies {
			row += 1
			searchResults.SetCell(row, 0, tview.NewTableCell(movie.Title).SetTextColor(tcell.ColorWhite).SetSelectable(true))
			searchResults.SetCell(row, 1, tview.NewTableCell(movie.ShowType).SetTextColor(tcell.ColorWhite).SetSelectable(true))
			searchResults.SetCell(row, 2, tview.NewTableCell(movie.Year).SetTextColor(tcell.ColorWhite).SetSelectable(true))
		}
		searchResults.SetSelectable(true, false).
			SetSelectedFunc(func(row, column int) {
				ratingViewPage, ratingViewPageFocus := RatingView(data.Movies[row-1], pages, searchInput, app)
				pages.AddAndSwitchToPage("rating", ratingViewPage, true)
				app.SetFocus(ratingViewPageFocus)
			})

		app.SetFocus(searchResults)
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchInput, 1, 0, false).
		AddItem(searchResults, 0, 1, true)

	return flex, searchInput
}
func refreshMainView() {
	movies, err := models.GetMovies()
	mainViewMovies = movies
	if err != nil {
		logger.Log.Fatalf("Error getting movies from database: %s", err.Error())
	}
	movieList.Clear()
	movieList.SetCell(0, 0, tview.NewTableCell("Type").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	movieList.SetCell(0, 1, tview.NewTableCell("Title").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	movieList.SetCell(0, 2, tview.NewTableCell("Genres").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	movieList.SetCell(0, 3, tview.NewTableCell("Year").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	movieList.SetCell(0, 4, tview.NewTableCell("Watched").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	movieList.SetCell(0, 5, tview.NewTableCell("Runtime").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	movieList.SetCell(0, 6, tview.NewTableCell("Imdb Rating").SetTextColor(tcell.ColorYellow).SetSelectable(false))
	movieList.SetCell(0, 7, tview.NewTableCell("User Score").SetTextColor(tcell.ColorYellow).SetSelectable(false))

	for row, movie := range movies {
		row += 1
		movieList.SetCell(row, 0, tview.NewTableCell(movie.ListType).SetTextColor(tcell.ColorWhite).SetSelectable(true))
		movieList.SetCell(row, 1, tview.NewTableCell(movie.Title).SetTextColor(tcell.ColorWhite).SetSelectable(true))
		movieList.SetCell(row, 2, tview.NewTableCell(movie.Genre).SetTextColor(tcell.ColorWhite).SetSelectable(true))
		movieList.SetCell(row, 3, tview.NewTableCell(movie.Year).SetTextColor(tcell.ColorWhite).SetSelectable(true).SetAlign(tview.AlignRight))
		movieList.SetCell(row, 4, tview.NewTableCell(movie.DateWatched).SetTextColor(tcell.ColorWhite).SetSelectable(true))
		movieList.SetCell(row, 5, tview.NewTableCell(movie.Runtime).SetTextColor(tcell.ColorWhite).SetSelectable(true))
		movieList.SetCell(row, 6, tview.NewTableCell(movie.ImdbRating).SetTextColor(tcell.ColorWhite).SetSelectable(true).SetAlign(tview.AlignRight))
		movieList.SetCell(row, 7, tview.NewTableCell(strconv.FormatFloat(movie.UserScore, 'f', 1, 64)).SetTextColor(tcell.ColorWhite).SetSelectable(true).SetAlign(tview.AlignRight))
	}
}
func updateDescriptionAndNotes(description, notes *tview.TextView, row int) {
	if len(mainViewMovies) == 0 {
		return
	}
	description.SetText("[yellow]Plot: [white]" + mainViewMovies[row].Plot + "\n\n")
	notes.SetText("[pink]" + mainViewMovies[row].Notes + "[white]")
}

var mainViewMovies []models.DBMovie
var movieList *tview.Table

func MainView() (*tview.Flex, tview.Primitive) {
	var movieDescription, movieNotes *tview.TextView
	movieList = tview.NewTable().SetSelectionChangedFunc(func(row, _ int) {
		if row == 0 {
			return
		}
		updateDescriptionAndNotes(movieDescription, movieNotes, row-1)
	})

	refreshMainView()
	movieList.SetSeparator(tview.Borders.Vertical).
		SetSelectable(true, false)

	movieDescription = tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true)
	movieNotes = tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true)
	updateDescriptionAndNotes(movieDescription, movieNotes, 0)
	rightFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(movieDescription, 0, 1, false).
		AddItem(movieNotes, 0, 1, false)

	flex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(movieList, 0, 3, true).
		AddItem(rightFlex, 0, 2, false)

	return flex, movieList
}
func RatingView(m models.Movie, pages *tview.Pages, searchInput *tview.InputField, app *tview.Application) (*tview.Flex, tview.Primitive) {
	movie := models.GetMovieDetails(m.ImdbId)
	movieDescription := tview.NewTextView().
		SetTextAlign(tview.AlignLeft)
	movieDescription.SetText("[yellow]Title: [white]" + movie.Title + "\n" +
		"[yellow]Year: [white]" + movie.Year + "\n" +
		"[yellow]Released: [white]" + movie.Released + "\n" +
		"[yellow]Rated: [white]" + movie.Rated + "\n" +
		"[yellow]Runtime: [white]" + movie.Runtime + "\n" +
		"[yellow]Genre: [white]" + movie.Genre + "\n" +
		"[yellow]Director: [white]" + movie.Director + "\n\n" +
		"[yellow]Plot: [white]" + movie.Plot + "\n").
		SetDynamicColors(true).
		SetBorder(true).
		SetTitle("Description").
		SetTitleColor(tcell.ColorGreen)

	form := tview.NewForm().
		AddInputField("Rating", "", 6, tview.InputFieldFloat, nil).
		AddTextArea("Notes", "", 0, 0, 0, nil).
		AddInputField("Date Watched", time.Now().Format("01/02/2006"), 10, tview.InputFieldMaxLength(10), nil).
		AddDropDown("List Type", []string{"watched", "to-watch", "dropped"}, 0, nil)
	form.AddButton("Save", func() {
		rating := form.GetFormItemByLabel("Rating").(*tview.InputField).GetText()
		notes := form.GetFormItemByLabel("Notes").(*tview.TextArea).GetText()
		dateText := form.GetFormItemByLabel("Date Watched").(*tview.InputField).GetText()
		date, err := time.Parse("01/02/2006", dateText)
		if err != nil {
			//modal
			logger.Log.Printf("Error parsing date: %s", err.Error())
			return
		}
		_, listType := form.GetFormItemByLabel("List Type").(*tview.DropDown).GetCurrentOption()
		if rating == "" {
			logger.Log.Println("Rating is required")
			return
		}
		score, err := strconv.ParseFloat(rating, 64)
		if err != nil {
			logger.Log.Fatalf("Error parsing rating: %s", err.Error())
		}
		dbMovie := movie.ToDBMovie(score)
		dbMovie.Notes = notes
		dbMovie.ListType = listType
		dbMovie.DateWatched = date.Format("2006-01-02")
		dbMovie.Watched = listType == "watched"
		dbMovie.UserScore = score
		_, err = models.AddMovie(dbMovie)
		if err != nil {
			logger.Log.Fatalf("Error adding movie to database: %s", err.Error())
		}
		//navigate back to main view
		refreshMainView()
		pages.SwitchToPage("view")
		pages.RemovePage("rating")
	}).
		AddButton("Cancel", func() {
			pages.SwitchToPage("search")
			app.SetFocus(searchInput)
			pages.RemovePage("rating")
		})
	form.SetTitleColor(tcell.ColorGreen).
		SetBorderPadding(1, 1, 1, 1)
	flex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(movieDescription, 0, 1, false).
		AddItem(form, 0, 1, false)
	return flex, form
}
