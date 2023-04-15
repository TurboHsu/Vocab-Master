package aid

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
)

var ChangeTopic31IndicatorWorkingMode bool = false
var IsEnabled bool = false
var bitmapNextBtn, bitmapIndicator robotgo.CBitmap
var WaitSecEntry *widget.Entry
var posMap map[string][2]int

func GenerateNewWindow(app *fyne.App) (window fyne.Window) {
	window = (*app).NewWindow("Aid")
	topic31ModeTrigger := widget.NewCheck("Change Topic 31 Indicator Working Mode", func(b bool) {
		ChangeTopic31IndicatorWorkingMode = b
	})
	info := widget.NewLabel("Aid not started.")
	captureNextBtn := widget.NewButton("Capture Next", func() {
		pos := getPosByHotkey("`")
		bitmapNextBtn = captureScreen(pos)
	}) 
	captureIndicatorBtn := widget.NewButton("Capture Indicator", func() {
		pos := getPosByHotkey("`")
		bitmapIndicator = captureScreen(pos)
	})
	enableLoopTrigger := widget.NewCheck("Enable Loop", func(b bool) {
		IsEnabled = b
		// call loop
		go eventLoop()
	})
	WaitSecEntry = widget.NewEntry()
	WaitSecEntry.SetText("1000")

	adjustBtn := widget.NewButton("Adjust Offset", func() {
		adjustOffset()
	})

	/**/





	window.SetContent(
		container.NewVBox(
			info,
			topic31ModeTrigger,
			container.NewHBox(
				widget.NewLabel("Wait Millisec:"),
				WaitSecEntry,
			),
			container.NewHBox(
				widget.NewLabel("Dual Monitor X axis Offset:"),
				adjustBtn,
			),
			container.NewHBox(
				captureNextBtn,
				captureIndicatorBtn,
			),
			enableLoopTrigger,
		),
	)

	return window
}
