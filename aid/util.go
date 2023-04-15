package aid

import (
	"github.com/go-vgo/robotgo"
)

// getPosByHotkey get position by pressing single hotkey twice
func getPosByHotkey(key string) (ret [2][2]int) {
	index := 0
	for {
		if index == 2 {
			return
		}
		if ok := robotgo.AddEvent(key); ok {
			posX, posY := robotgo.GetMousePos()
			ret[index] = [2]int{posX, posY}
			index++
		}
	}
}

// getPosByHotkey get position by pressing single hotkey twice
func getOnePosByHotkey(key string) (ret [2]int) {
	for {
		if ok := robotgo.AddEvent(key); ok {
			posX, posY := robotgo.GetMousePos()
			ret = [2]int{posX, posY}
			return
		}
	}
}

func captureScreen(pos [2][2]int) (ret robotgo.CBitmap) {
	bitmap := robotgo.CaptureScreen(min(pos[0][0], pos[1][0]), min(pos[0][1], pos[1][1]),
		abs(pos[1][0]-pos[0][0]), abs(pos[1][1]-pos[0][1]))
	ret = robotgo.CBitmap(bitmap)
	// robotgo.FreeBitmap(bitmap)
	return
}

// abs
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// min
func min(x int, y int) int {
	if x < y {
		return x
	}
	return y
}
