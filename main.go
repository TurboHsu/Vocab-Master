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
	"github.com/TurboHsu/Vocab-Master/automatic"
	"github.com/TurboHsu/Vocab-Master/grab"
	"github.com/lqqyt2423/go-mitmproxy/proxy"
)

const version = "1.2.0"

// Declare some global variable
var dataset grab.VocabDataset
var window fyne.Window
var toggle *widget.Check
var auto *widget.Button
var openBtn *widget.Button
var jsHijackCheck *widget.Check

// Declare some flags
var shouldOperateProxy bool = true
var jsHijack bool = false

// Initiate the flag parser
func init() {
	flag.BoolVar(&shouldOperateProxy, "proxy", true, "Operates the system proxy when start and stop.")
}

func main() {
	flag.Parse()

	//DEBUG
	//shouldOperateProxy = false

	platform, err := GetPlatform()
	if err != nil {
		//Platform specific failed, use default
		platform = GetDefaultPlatform()
	}

	//Init font
	os.Setenv("FYNE_FONT", platform.Font)

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
		CaRootPath:        platform.CertDir,
	}

	p, err := proxy.NewProxy(opts)
	if err != nil {
		log.Fatal(err)
	}

	p.AddAddon(&VocabMasterHandler{})

	a := app.New()
	window = a.NewWindow("Vocab Master! " + version)
	label := widget.NewLabel(
		"Hey! Here is Vocab Master.\n" +
			"Install this cert as trusted root: " + opts.CaRootPath + "\n" +
			"Then start a class task, program will run itself ;)\n\n" +
			"Project addr: github.com/TurboHsu/VocabMaster\n" +
			"Font in use: " + platform.Font)
	toggle = widget.NewCheck("Enable processor", func(b bool) {
		IsEnabled = b
		fmt.Println("Processor enabler is set to ", IsEnabled)
	})
	toggle.Checked = true

	auto = widget.NewButton("Automation", func() {
		autoWindow := automatic.GenerateNewWindow(&a)
		autoWindow.Show()
	})

	jsHijackCheck = widget.NewCheck("Enable JS hijack", func(b bool) {
		jsHijack = b
		jsHijackCheck.Text = "Waiting for JS Hijacking..."
	})

	openBtn = widget.NewButton("Open certificates directory", func() {
		platform.OpenCertDir()
	})

	window.SetContent(
		container.NewVBox(label,
			container.NewHBox(toggle, jsHijackCheck),
			container.NewHBox(openBtn, auto),
		),
	)

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

func GetDefaultPlatform() Platform {
	return Platform{
		DataDir: ".",
		CertDir: "./cert",
		Font:    "./font/red_bean.ttf",
	}
}
