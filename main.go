package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/lqqyt2423/go-mitmproxy/proxy"
)

// Declear some global variable
var dataset VocabDataset = VocabDataset{IsEnabled: true}
var words []WordInfo
var window fyne.Window
var toggler *widget.Check

// Declear some flags
var shouldOperateProxy bool = true

// Initiate the flag parser
func init() {
	flag.BoolVar(&shouldOperateProxy, "proxy", true, "Operates the system proxy when start and stop.")
}

func main() {
	flag.Parse()

	//Init font
	os.Setenv("FYNE_FONT", "./font/red_bean.ttf")

	var originStatus []ProxyState
	if shouldOperateProxy {
		log.Println("Processing proxy settings...")

		//Get proxy
		var err error
		originStatus, err = ReadSystemStatus()
		if err != nil {
			panic(err)
		}

		//Set proxy
		err = SetSystemProxy("localhost:38848")
		if err != nil {
			panic(err)
		}
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
	toggler = widget.NewCheck("Enable processor", func(b bool) {
		dataset.IsEnabled = b
		fmt.Println("Processor enabler is set to ", dataset.IsEnabled)
	})
	toggler.Checked = true
	window.SetContent(container.NewVBox(label, toggler))

	go p.Start()
	window.ShowAndRun()

	//Unset font
	os.Unsetenv("FYNE_FONT")

	if shouldOperateProxy {
		//Unset proxy
		err = ApplyProxyStatus(originStatus)
		if err != nil {
			panic(err)
		}
	}
}
