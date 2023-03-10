package automatic

import (
	"log"
	"strconv"
	"time"

	"github.com/TurboHsu/Vocab-Master/grab"
)

var dataset grab.VocabDataset
var Enabled bool
var pendingWord []string

func FetchDataset(data grab.VocabDataset) {
	dataset = data
}

func startAutomation() {

	// Grab all words
	info.SetText("Grabbing words...")
	progressBar.SetValue(0)
	progressBar.Show()
	for i := 0; i < len(pendingWord); i++ {
		progressBar.SetValue(float64(i) / float64(len(pendingWord)))
		// Grab word
		grab.GrabWord(pendingWord[i], (*grab.VocabDataset)(&dataset))
		info.SetText("Grabbing word: " + pendingWord[i])
		log.Println("Grabbing word: ", pendingWord[i])
	}
	progressBar.SetValue(1)
	info.SetText("Grabbing words finished.")

	var finnished bool
	timeDelay, _ := strconv.Atoi(timeDelayEntry.Text)
	for !finnished {
		// Delay for specific time
		info.SetText("Delaying for " + timeDelayEntry.Text + " seconds...")
		time.Sleep(time.Duration(timeDelay) * time.Second)

		//TODO

		// Make the first request.
		//response := doGET(fmt.Sprintf("https://app.vocabgo.com/api/Student/%s/StartAnswer?task_id=%d&",
		//	dataset.CurrentTask.TaskType))

	}

}
