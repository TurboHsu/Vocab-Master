package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/lqqyt2423/go-mitmproxy/proxy"
	log "github.com/sirupsen/logrus"
	"os"
)

var dataset VocabDataset
var words []WordInfo
var window fyne.Window

func main() {
	//Init font
	os.Setenv("FYNE_FONT", "./font/red_bean.ttf")

	//Windows Gets proxy
	originStatus, err := ReadSystemStatus()
	if err != nil {
		panic(err)
	}

	//Set proxy
	err = SetSystemProxy("localhost:38848")
	if err != nil {
		panic(err)
	}

	opts := &proxy.Options{
		Addr:              "localhost:38848",
		StreamLargeBodies: 1024 * 1024 * 5,
		CaRootPath:        "./cert",
	}

	p, err := proxy.NewProxy(opts)
	if err != nil {
		log.Fatal(err)
	}

	p.AddAddon(&VocabMasterHandler{})

	a := app.New()
	window = a.NewWindow("Vocab Master!")
	label := widget.NewLabel("Hey! Here is Vocab Master.\nJust start a class task, program will run itself ;)\n\nProject addr: github.com/TurboHsu/VocabMaster")
	window.SetContent(label)

	go p.Start()
	window.ShowAndRun()

	//Unset font
	os.Unsetenv("FYNE_FONT")
	//Unset proxy
	err = ApplyProxyStatus(originStatus)
	if err != nil {
		panic(err)
	}
}
