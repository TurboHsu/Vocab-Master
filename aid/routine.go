package aid

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-vgo/robotgo"
)

var offset [2]int = [2]int{0,0}

func eventLoop() {
	// parse the wait sec
	sec, _ := strconv.Atoi(WaitSecEntry.Text)
	for {
		if !IsEnabled {
			break
		}
		srcBitmap := robotgo.CaptureScreen()
		x, y := robotgo.FindBitmap(robotgo.ToMMBitmapRef(bitmapIndicator), srcBitmap, 0.35)
		fmt.Println(x,y)
		if x != -1 && y != -1 {
			x += offset[0]
			y += offset[1]
			robotgo.Move(x, y)
			robotgo.Click("left")
		}
		time.Sleep(time.Millisecond * time.Duration(sec))
		x, y = robotgo.FindBitmap(robotgo.ToMMBitmapRef(bitmapNextBtn), srcBitmap, 0.2)
		fmt.Println(x, y)
		if x != -1 && y != -1 {
			x += offset[0]
			y += offset[1]
			robotgo.Move(x, y)
			robotgo.Click("left")
		}
		time.Sleep(time.Millisecond * time.Duration(sec))
		robotgo.FreeBitmap(srcBitmap)
	}
}

func adjustOffset() {
	srcBitmap := robotgo.CaptureScreen()
	pos := getPosByHotkey("`")
	x, y := robotgo.FindBitmap(robotgo.ToMMBitmapRef(bitmapIndicator), srcBitmap, 0.2)
	offset[0] = pos[0][0] - x
	offset[1] = pos[0][1] - y
	robotgo.FreeBitmap(srcBitmap)
}