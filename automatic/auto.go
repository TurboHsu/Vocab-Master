package automatic

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
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
	taskPath := "ClassTask" //TODO: This field is not implemented, may get from url from mitm hook.

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

	// Submit the words need to be done.
	// Summon the json
	var wordListSubmitJSON string
	for i := 0; i < len(pendingWord); i++ {
		wordListSubmitJSON += `"` + pendingWord[i] + `",`
	}
	jsonSubmitWord := fmt.Sprintf(`{"task_id":%s,"word_map":{"%s:%s":[%s]},"chose_err_item":2,"reset_chose_words":1,"timestamp":%d,"version":"%s","sign":"%s","app_type":%s}`,
		taskDetail.TaskID,
		dataset.CurrentTask.TaskID,
		dataset.CurrentTask.TaskSet,
		wordListSubmitJSON[:len(wordListSubmitJSON)-1], // Trim the last comma
		time.Now().UnixMilli(),
		taskDetail.Versions,
		summonSign(),
		taskDetail.AppType,
	)
	info.SetText("Submitting word list...")
	// Generate the url
	// call get
	url := fmt.Sprintf("https://app.vocabgo.com/api/Student/%s/SubmitChoseWord",
		taskPath,
	)
	resp := doPOST(url, jsonSubmitWord)
	if !strings.Contains(resp, `"code":1`) {
		info.SetText("Failed to submit word list.")
		log.Println("Failed to submit word list.")
		return
	}
	// Get last time
	lastTime := time.Now().UnixMilli()
	// Submit the first request
	info.SetText("Submitting the first request...")
	// Generate the url
	url = fmt.Sprintf(`https://app.vocabgo.com/student/api/Student/%s/StartAnswer?task_id=%s&task_type=%s&release_id=%s&opt_img_w=1947&opt_font_size=108&opt_font_c=%%23000000&it_img_w=2287&it_font_size=121&timestamp=%d&version=%s&app_type=%s`,
		taskPath,
		taskDetail.TaskID,
		taskDetail.TaskType,
		taskDetail.ReleaseID,
		lastTime,
		taskDetail.Versions,
		taskDetail.AppType,
	)
	resp = doGET(url)
	// Parse the response
	var startAnswerResponse StartAnswerResponseStruct
	json.Unmarshal([]byte(resp), &startAnswerResponse)
	if startAnswerResponse.Code != 1 {
		info.SetText("Failed to submit the first request.")
		log.Println("Failed to submit the first request.")
		return
	}
	// Get rid of salt and parse
	_, startAnswerResponseDesalted := splitSalt(startAnswerResponse.Data)
	// Decode the base64
	startAnswerResponseDesaltedDecoded, _ := base64.StdEncoding.DecodeString(startAnswerResponseDesalted)
	// Unmarshal the data using json
	var vocabTaskData VocabTaskStruct
	json.Unmarshal(startAnswerResponseDesaltedDecoded, &vocabTaskData)

	var topicCode string
	//var rawDecodedString string = string(startAnswerResponseDesaltedDecoded)
	timeDelay, _ := strconv.Atoi(timeDelayEntry.Text)
MainLoop:
	for {
		// Delay for specific time
		info.SetText("Delaying for " + timeDelayEntry.Text + " seconds...")
		time.Sleep(time.Duration(timeDelay) * time.Second)

		switch vocabTaskData.TopicMode {
		// This is letting u read through something.
		case 0:
			//Get the next TopicCode and do nothing.
			topicCode = vocabTaskData.TopicCode
			// WIP
			//case 11:
			//	ans := answer.FindAnswer(11, answer.VocabTaskStruct(vocabTaskData), "")
			//case 22:
			//	ans := answer.FindAnswer(22, answer.VocabTaskStruct(vocabTaskData), "")
			//case 31:
			//	ans := answer.FindAnswer(31, answer.VocabTaskStruct(vocabTaskData), "")
			//case 32:
			//	ans := answer.FindAnswer(32, answer.VocabTaskStruct(vocabTaskData), string(rawDecodedString))
			//case 51:
			//	ans := answer.FindAnswer(51, answer.VocabTaskStruct(vocabTaskData), string(rawDecodedString))

		}

		// Submit and save it
		var submitAnswerAndSaveData SubmitAnswerAndSaveStruct
		submitAnswerAndSaveData.TopicCode = topicCode
		submitAnswerAndSaveData.TimeSpent = int(time.Now().UnixMilli()) - int(lastTime)
		// These are static data
		submitAnswerAndSaveData.OptFontC = "#000000"
		submitAnswerAndSaveData.OptImgW = 1947
		submitAnswerAndSaveData.OptFontSize = 108
		submitAnswerAndSaveData.ItImgW = 2287
		submitAnswerAndSaveData.ItFontSize = 121
		submitAnswerAndSaveData.Timestamp = time.Now().UnixMilli()
		submitAnswerAndSaveData.Version = taskDetail.Versions
		submitAnswerAndSaveData.Sign = summonSign()
		submitAnswerAndSaveData.AppType, _ = strconv.Atoi(taskDetail.AppType)
		// Change last time
		lastTime = time.Now().UnixMilli()
		// Marshal the data
		submitAnswerAndSaveDataJSON, _ := json.Marshal(submitAnswerAndSaveData)
		// Submit the data
		url = fmt.Sprintf("https://app.vocabgo.com/student/api/Student/%s/SubmitAnswerAndSave", taskPath)
		resp = doPOST(url, string(submitAnswerAndSaveDataJSON))
		// Decode the data
		var submitAnswerAndSaveResponse StartAnswerResponseStruct
		// Parse the response
		json.Unmarshal([]byte(resp), &submitAnswerAndSaveResponse)
		// Check if the task is finished
		if submitAnswerAndSaveResponse.Code == 20001 {
			break MainLoop
		}
		// Check if failed
		if submitAnswerAndSaveResponse.Code != 1 && submitAnswerAndSaveResponse.Code != 20001 {
			info.SetText("Failed to submit the answer.")
			log.Println("Failed to submit the answer.")
			return
		}
		// Get rid of salt and parse
		_, submitAnswerAndSaveResponseDesalted := splitSalt(submitAnswerAndSaveResponse.Data)
		// Decode the base64
		submitAnswerAndSaveResponseDesaltedDecoded, _ := base64.StdEncoding.DecodeString(submitAnswerAndSaveResponseDesalted)
		// Unmarshal the data using json
		json.Unmarshal(submitAnswerAndSaveResponseDesaltedDecoded, &vocabTaskData)
		// Update the topic code
		topicCode = vocabTaskData.TopicCode
	}
	info.SetText("Automation finished.")
	progressBar.SetValue(1)
}

// The mechanism of signing is not sure, so now its some md5 from timestamp.
func summonSign() string {
	// return md5 of timestamp
	return fmt.Sprintf("%x", md5.Sum([]byte(strconv.Itoa(int(time.Now().Unix())))))
}
