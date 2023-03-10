package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/TurboHsu/Vocab-Master/answer"
	"github.com/TurboHsu/Vocab-Master/automatic"
	"github.com/andybalholm/brotli"
	"github.com/lqqyt2423/go-mitmproxy/proxy"
)

type VocabMasterHandler struct {
	proxy.BaseAddon
}

var infoLabel = widget.NewLabel("")
var progressBar *widget.ProgressBar

func (c *VocabMasterHandler) Request(f *proxy.Flow) {
	if f.Request.URL.Host != "app.vocabgo.com" {
		return
	}
	if strings.Contains(f.Request.URL.Path, "/api/Student/ClassTask/SubmitChoseWord") || strings.Contains(
		f.Request.URL.Path, "/api/Student/StudyTask/SubmitChoseWord") {
		//Adapt class task
		//Flush word storage and wordlist
		answer.WordList = []answer.WordInfo{}
		dataset.CurrentTask.WordList = []string{}

		//Put wordlist into dataset
		wordListRegex := regexp.MustCompile(`\[.*?\]`)
		wordListRaw := wordListRegex.FindString(string(f.Request.Body))
		json.Unmarshal([]byte(wordListRaw), &dataset.CurrentTask.WordList)

		//Put task info into dataset
		taskInfoRegex := regexp.MustCompile(`"word_map":{".*?"`)
		taskInfoRaw := taskInfoRegex.FindString(string(f.Request.Body))
		taskInfo := strings.Split(taskInfoRaw[13:len(taskInfoRaw)-1], ":")
		dataset.CurrentTask.TaskID = taskInfo[0]
		dataset.CurrentTask.TaskSet = taskInfo[1]

		//Update Cookie and Header
		dataset.RequestInfo.Cookies = f.Request.Raw().Cookies()
		dataset.RequestInfo.Header = f.Request.Raw().Header

		//Create a thread to crawl all words
		go func() {
			//Setup progress ui
			progressBar = widget.NewProgressBar()
			completeBox := widget.NewLabel("Gathering word info...")
			var wordList string
			if len(dataset.CurrentTask.WordList) > 8 {
				wordList = fmt.Sprintln(dataset.CurrentTask.WordList[:8]) + "..."
			} else {
				wordList = fmt.Sprintln(dataset.CurrentTask.WordList)
			}
			infoBox := widget.NewLabel("New task detected. Gathering chosen words' info:\n" + wordList)
			window.SetContent(container.NewVBox(infoBox, progressBar, completeBox, infoLabel, toggle))

			for i := 0; i < len(dataset.CurrentTask.WordList); i++ {
				//Show progress
				progressBar.SetValue(float64(i) / float64(len(dataset.CurrentTask.WordList)))
				grabWord(dataset.CurrentTask.WordList[i])
				log.Println("[I] Grabbed word list:" + dataset.CurrentTask.WordList[i])
			}
			progressBar.SetValue(1)
			completeBox.SetText("Complete!")

			// If is auto, use auto function
			if dataset.IsAuto {
				// Get task type
				var taskType string
				if strings.Contains(f.Request.URL.Path, "/api/Student/ClassTask/SubmitChoseWord") {
					taskType = "Class"
				} else {
					taskType = "Study"
				}

				automatic.FetchDataset(automatic.VocabDataset{
					CurrentTask: automatic.CurrentTask{
						WordList: dataset.CurrentTask.WordList,
						TaskSet:  dataset.CurrentTask.TaskSet,
						TaskID:   dataset.CurrentTask.TaskID,
						TaskType: taskType,
					},
					RequestInfo: automatic.RequestInfo{
						Versions: dataset.RequestInfo.Versions,
						Cookies:  dataset.RequestInfo.Cookies,
						Header:   dataset.RequestInfo.Header,
					},
				})
				automatic.StartAutomation()
			}
		}()
	}

}

func (c *VocabMasterHandler) Response(f *proxy.Flow) {
	//Judge whether session is from vocabgo task
	if f.Request.URL.Host != "app.vocabgo.com" {
		return
	}
	//Adapt class task
	if !strings.Contains(
		f.Request.URL.Path, "/api/Student/ClassTask/SubmitAnswerAndSave") && !strings.Contains(
		f.Request.URL.Path, "/api/Student/ClassTask/StartAnswer") && !strings.Contains(
		f.Request.URL.Path, "/api/Student/StudyTask/SubmitAnswerAndSave") && !strings.Contains(
		f.Request.URL.Path, "/api/Student/StudyTask/StartAnswer") {
		return
	}

	// Automated actions should not be MITMed.
	if dataset.IsAuto {
		return
	}

	//Switch of processor
	if dataset.IsEnabled {

		//If the progress bar has hit 100%, then hide it
		if progressBar.Value == 1 {
			progressBar.Hide()
		}

		//Get decoded content
		rawByte, _ := f.Response.DecodedBody()

		//Okay! Let's decode raw json
		var vocabRawJSON VocabRawJSONStruct
		json.Unmarshal(rawByte, &vocabRawJSON)

		//Judge whether is the last task
		if vocabRawJSON.Msg != "success" {
			return
		}

		//JSON Salt
		JSONSalt := vocabRawJSON.Data[:32]
		rawJSONBase64 := vocabRawJSON.Data[32:]

		//Let's get the insider base64-encoded info
		rawDecodedString, err := base64.StdEncoding.DecodeString(rawJSONBase64)
		if err != nil {
			panic(err)
		}

		//Judge whether json contains task info
		if !strings.Contains(string(rawDecodedString), "task_id") {
			return
		}

		var vocabTask VocabTaskStruct
		json.Unmarshal(rawDecodedString, &vocabTask)

		//Switch for tasks
		switch vocabTask.TopicMode {
		//Introducing words
		case 0:
			//UI
			infoLabel.SetText("Seems like you are entering an new task!\nPlease wait until progress bar reach 100%.")
		//Choose translation of specific word from a sentence
		case 11:
			ans := answer.FindAnswer(11, answer.VocabTaskStruct(vocabTask), "")

			//UI
			if ans.Found && len(ans.Index) > 0 {
				infoLabel.SetText("Hey! The answer is tagged out.\nAnd the answer is [" + ans.Detail.Translation + "]")
			} else {
				infoLabel.SetText("Warn: Answer not found. This might be a bug.")
			}

			//Check whether index is found
			if len(ans.Index) > 0 {
				//Tag out the correct answer
				regex := regexp.MustCompile(`（.*?）`)
				newJSON := string(rawDecodedString)
				newJSON = string(regex.ReplaceAll([]byte(newJSON), []byte("")))
				newJSON = strings.Replace(newJSON, vocabTask.Options[ans.Index[0]].Content, "-> "+vocabTask.Options[ans.Index[0]].Content+" <-", 1)
				//newJSON := strings.Replace(string(rawDecodedString), vocabTask.Stem.Content, vocabTask.Stem.Content+" ["+translation+"]", 1)
				repackedBase64 := base64.StdEncoding.EncodeToString([]byte(newJSON))
				vocabRawJSON.Data = JSONSalt + repackedBase64
				body, _ := json.Marshal(vocabRawJSON)
				var b bytes.Buffer
				br := brotli.NewWriter(&b)
				br.Write(body)
				br.Flush()
				br.Close()
				f.Response.Body = b.Bytes()
			}

		//Choose word from voice
		case 22:
			ans := answer.FindAnswer(22, answer.VocabTaskStruct(vocabTask), "")

			//UI
			if ans.Found && len(ans.Index) > 0 {
				infoLabel.SetText("Hey! The answer is tagged out.\nAnd the answer is [" + vocabTask.Options[ans.Index[0]].Content + "]")
			} else {
				infoLabel.SetText("Warn: Answer not found. This might be a bug.")
			}

			if len(ans.Index) > 0 {
				//Tag out the correct answer
				regex := regexp.MustCompile(`（.*?）`)
				newJSON := string(rawDecodedString)
				newJSON = string(regex.ReplaceAll([]byte(newJSON), []byte("")))
				newJSON = strings.Replace(newJSON, vocabTask.Options[ans.Index[0]].Content, "-> "+vocabTask.Options[ans.Index[0]].Content+" <-", 1)
				repackedBase64 := base64.StdEncoding.EncodeToString([]byte(newJSON))
				vocabRawJSON.Data = JSONSalt + repackedBase64
				body, _ := json.Marshal(vocabRawJSON)
				var b bytes.Buffer
				br := brotli.NewWriter(&b)
				br.Write(body)
				br.Flush()
				br.Close()
				f.Response.Body = b.Bytes()
			}

		//Choose word pair
		case 31:
			ans := answer.FindAnswer(31, answer.VocabTaskStruct(vocabTask), "")
			var detag []string

			// Find incorrect
			for i := 0; i < len(vocabTask.Options); i++ {
				var isCorrect bool
				for _, corrIndex := range ans.Index {
					if i == corrIndex {
						isCorrect = true
						break
					}
				}
				// Not correct, append the content in detag
				if !isCorrect {
					detag = append(detag, vocabTask.Options[i].Content)
				}
			}
			infoLabel.SetText("The incorrect answer is tagged out. \n and the correct index is:\n" + fmt.Sprintln(ans.Index))

			//Show answer in the UI

			newJSON := string(rawDecodedString)
			for i := 0; i < len(detag); i++ {
				//newJSON = strings.Replace(newJSON, `"content":"`+detag[i]+`"`, `"content":"`+"NOT-["+detag[i]+"]-THIS"+`"`, 1)
				newJSON = strings.Replace(newJSON, `"content":"`+detag[i]+`"`, `"content":"错误选项"`, 1)
			}
			repackedBase64 := base64.StdEncoding.EncodeToString([]byte(newJSON))
			vocabRawJSON.Data = JSONSalt + repackedBase64
			body, _ := json.Marshal(vocabRawJSON)
			var b bytes.Buffer
			br := brotli.NewWriter(&b)
			br.Write(body)
			br.Flush()
			br.Close()
			f.Response.Body = b.Bytes()

		//Organize word pieces
		case 32:
			ans := answer.FindAnswer(32, answer.VocabTaskStruct(vocabTask), string(rawDecodedString))
			// Get some remark

			if ans.Found {
				//UI
				infoLabel.SetText("Hey! The answer is printed out.\nAnd the answer is [" + ans.Detail.Word + "]")
				//Change the hint to the correct answer
				newJSON := strings.Replace(string(rawDecodedString), ans.Detail.Raw, ans.Detail.Word, 1)
				repackedBase64 := base64.StdEncoding.EncodeToString([]byte(newJSON))
				vocabRawJSON.Data = JSONSalt + repackedBase64
				body, _ := json.Marshal(vocabRawJSON)
				var b bytes.Buffer
				br := brotli.NewWriter(&b)
				br.Write(body)
				br.Flush()
				br.Close()
				f.Response.Body = b.Bytes()
			} else {
				infoLabel.SetText("Warn: Answer not found. This might be a bug.")
			}
		//Write words from first chars
		case 51:
			ans := answer.FindAnswer(51, answer.VocabTaskStruct(vocabTask), string(rawDecodedString))

			//UI
			if ans.Found {
				infoLabel.SetText("Hey! The answer is printed out. \n And the answer is [" + ans.Detail.Word + "]")

				//Change the tip
				regexReplaceJSON := regexp.MustCompile(`"w_tip":".*?"`)
				regexGetWord := regexp.MustCompile(`{.*?}`)
				theWord := regexGetWord.FindString(ans.Detail.Word)
				newJSON := regexReplaceJSON.ReplaceAllString(string(rawDecodedString), `"w_tip":"`+theWord[1:len(theWord)-1]+`"`)

				//Change the translation
				newJSON = strings.Replace(newJSON, ans.Detail.Raw, ans.Detail.Word, 1)

				repackedBase64 := base64.StdEncoding.EncodeToString([]byte(newJSON))
				vocabRawJSON.Data = JSONSalt + repackedBase64
				body, _ := json.Marshal(vocabRawJSON)
				var b bytes.Buffer
				br := brotli.NewWriter(&b)
				br.Write(body)
				br.Flush()
				br.Close()
				f.Response.Body = b.Bytes()
			} else if ans.Detail.Uncertain {
				infoLabel.SetText("The answer cannot be find in phrases,\nbut we found one through fuzzy queries.\nIt is " + ans.Detail.Word)

				//Change the tip
				regexReplaceJSON := regexp.MustCompile(`"w_tip":".*?"`)
				newJSON := regexReplaceJSON.ReplaceAllString(string(rawDecodedString), `"w_tip":"`+ans.Detail.Word+`"`)

				//Change the translation
				regexReplaceJSON = regexp.MustCompile(`"remark":".*?"`)
				result := regexReplaceJSON.Find([]byte(newJSON))
				newJSON = regexReplaceJSON.ReplaceAllString(newJSON, string(result)[:len(string(result))-1]+" Possible answer:"+ans.Detail.Word+`"`)

				repackedBase64 := base64.StdEncoding.EncodeToString([]byte(newJSON))
				vocabRawJSON.Data = JSONSalt + repackedBase64
				body, _ := json.Marshal(vocabRawJSON)
				var b bytes.Buffer
				br := brotli.NewWriter(&b)
				br.Write(body)
				br.Flush()
				br.Close()
				f.Response.Body = b.Bytes()
			} else {
				infoLabel.SetText("Warn: Answer not found. This might be a bug.")
			}

		default:
			infoLabel.SetText("This task is not supported or this is a bug.\n")
		}
	} else {
		infoLabel.SetText("Processor is disabled.\n")
	}
}
