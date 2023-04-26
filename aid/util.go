package aid

import (
	"encoding/json"
	"log"
	"os"

	"github.com/go-vgo/robotgo"
)

// savePosPreset saves the posMap
func savePosPreset() {
	// Save the posMap
	byteArr, _ := json.Marshal(PosMapJSON{PosMap: posMap})
	file, err := os.OpenFile("posMap.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error when saving posMap:", err)
		return
	}
	defer file.Close()
	os.WriteFile("posMap.json", byteArr, 0644)
}

func loadPosPreset() {
	// Check whether the file exists
	if _, err := os.Stat("posMap.json"); os.IsNotExist(err) {
		log.Println("posMap.json not found.")
		return
	}
	file, err := os.OpenFile("posMap.json", os.O_RDONLY, 0644)
	if err != nil {
		log.Println("Error when loading posMap:", err)
		return
	}
	defer file.Close()
	// Read the file
	byteArr := make([]byte, 1024)
	n, err := file.Read(byteArr)
	if err != nil {
		log.Println("Error when loading posMap:", err)
		return
	}
	// Parse the file
	var posMapJSON PosMapJSON
	err = json.Unmarshal(byteArr[:n], &posMapJSON)
	if err != nil {
		log.Println("Error when loading posMap:", err)
		return
	}
	posMap = posMapJSON.PosMap
}


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

type PosMapJSON struct {
	PosMap map[string][2]int `json:"posMap"`
}