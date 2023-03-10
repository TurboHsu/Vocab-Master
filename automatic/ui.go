package automatic

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var timeDelayEntry *widget.Entry
var valid bool
var confirmed bool
var info *widget.Label
var progressBar *widget.ProgressBar
var list *widget.List
var doGetFullScore bool

func GenerateNewWindow(app *fyne.App) (window fyne.Window) {
	window = (*app).NewWindow("Automation Console")
	trigger := widget.NewCheck("Enable Automation", func(b bool) {
		Enabled = b
	})
	info = widget.NewLabel("Please enter a task page to get info.")
	progressBar = widget.NewProgressBar()
	progressBar.Hidden = true
	list = widget.NewList(
		func() int {
			return len(pendingWord)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("List of words")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(pendingWord[i])
		},
	)
	timeDelayEntry = widget.NewEntry()
	timeDelayEntry.SetText("5")
	doGetFullScoreCheck := widget.NewCheck("Obtain words without full score", func(b bool) {
		log.Println("doGetFullScore: ", b)
		doGetFullScore = b
	})

	startAutomationButton := widget.NewButton("Start Automation", func() {
		log.Println("startAutomationButton clicked, (valid, confirmed):", valid, confirmed)
		if valid && !confirmed {
			// Show menu for users to confirm
			info.SetText("Check words in the list.\nIf u r sure what r u going to do, click the button again.")
			confirmed = true
		} else if valid && confirmed {
			// Start automation
			startAutomation()
		} else {
			dialog.NewInformation("Hmm", "I have no idea what u r going to automate.\nGo ahead and open a task page!", window).Show()
		}
	})

	window.SetContent(
		container.NewHBox(
			list,
			container.NewVBox(
				widget.NewLabel("Vocab Master Automation Console"),
				trigger,
				doGetFullScoreCheck,
				container.NewHBox(
					widget.NewLabel("Time delay (seconds):"),
					timeDelayEntry,
				),
				progressBar,
				startAutomationButton,
				info,
			),
		))
	return
}
