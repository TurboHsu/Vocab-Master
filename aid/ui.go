package aid

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var ChangeTopic31IndicatorWorkingMode bool = false
var IsEnabled bool = false
var WaitSecEntry *widget.Entry
var posMap map[string][2]int = make(map[string][2]int)

func GenerateNewWindow(app *fyne.App) (window fyne.Window) {
	window = (*app).NewWindow("Aid")
	topic31ModeTrigger := widget.NewCheck("Change Topic 31 Indicator Working Mode", func(b bool) {
		ChangeTopic31IndicatorWorkingMode = b
	})
	info := widget.NewLabel("Aid.")
	enableLoopTrigger := widget.NewCheck("Enable Loop", func(b bool) {
		IsEnabled = b
		// call loop
		go eventLoop()
	})
	WaitSecEntry = widget.NewEntry()
	WaitSecEntry.SetText("1000")

	/* Bunch of buttons */
	getVerticalOptionPosBtn := widget.NewButton("Get Vertical Option Mapping", func() {
		for i := 0; i < 4; i++ {
			pos := getOnePosByHotkey("`")
			posMap["v"+strconv.Itoa(i)] = pos
		}
	})
	getHorizonalOptionPosBtn := widget.NewButton("Get Horizonal Option Mapping", func() {
		for i := 0; i < 6; i++ {
			pos := getOnePosByHotkey("`")
			posMap["h"+strconv.Itoa(i)] = pos
		}
	})
	getNextPosBtn := widget.NewButton("Get Next Pos", func() {
		pos := getOnePosByHotkey("`")
		posMap["next"] = pos
	})
	getSubmitBtn := widget.NewButton("Get Submit Pos", func() {
		pos := getOnePosByHotkey("`")
		posMap["submit"] = pos
	})

	window.SetContent(
		container.NewVBox(
			info,
			topic31ModeTrigger,
			container.NewHBox(
				widget.NewLabel("Wait Millisec:"),
				WaitSecEntry,
			),
			container.NewHBox(
				getVerticalOptionPosBtn,
				getHorizonalOptionPosBtn,
				getNextPosBtn,
				getSubmitBtn,
			),
			container.NewHBox(
				widget.NewButton("Save PosMap", func() {
					savePosPreset()
				}),
				widget.NewButton("Load PosMap", func() {
					loadPosPreset()
				}),
			),
			enableLoopTrigger,
		),
	)

	return window
}
