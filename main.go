package main

import (
	"fmt"
	"os"
	"runtime"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Trisia/gosysproxy"
	"github.com/lqqyt2423/go-mitmproxy/proxy"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/registry"
)

var textBox = widget.NewLabel("Vocab Master Started.\nWaiting for task begin..")
var wordCache [][]string
var dataset VocabDataset
var words []WordInfo

func main() {
	//Init font
	os.Setenv("FYNE_FONT", "./font/red_bean.ttf")

	//Windows Gets proxy
	var proxyEnableRaw uint64
	var proxyServerRaw string
	if runtime.GOOS == "windows" {
		proxyReg, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.QUERY_VALUE)
		if err != nil {
			panic(err)
		}
		defer proxyReg.Close()
		proxyEnableRaw, _, err = proxyReg.GetIntegerValue("ProxyEnable")
		if err != nil {
			panic(err)
		}
		proxyServerRaw, _, err = proxyReg.GetStringValue("ProxyServer")
		if err != nil {
			panic(err)
		}
	}

	//Set proxy
	err := gosysproxy.SetGlobalProxy("localhost:38848")
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
	w := a.NewWindow("Vocab Master!")
	clearButton := widget.NewButton("Clear vocabulary cache", func() {
		wordCache = [][]string{}
		textBox.SetText("Cleared!")
	})
	testButton := widget.NewButton("test", func() {
		fmt.Println(words)
	})
	w.SetContent(container.NewVBox(textBox, clearButton, testButton))

	go p.Start()
	w.ShowAndRun()

	//Unset font
	os.Unsetenv("FYNE_FONT")
	//Unset proxy
	if runtime.GOOS == "windows" {
		if proxyEnableRaw == 1 {
			err := gosysproxy.SetGlobalProxy(proxyServerRaw)
			if err != nil {
				panic(err)
			}
		} else {
			err := gosysproxy.Off()
			if err != nil {
				panic(err)
			}
		}
	}
	err = gosysproxy.Off()
	if err != nil {
		panic(err)
	}
}
