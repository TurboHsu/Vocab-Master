package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/andybalholm/brotli"
	"github.com/lqqyt2423/go-mitmproxy/proxy"
)

type VocabMasterHandler struct {
	proxy.BaseAddon
}

var infoLable = widget.NewLabel("")

func (c *VocabMasterHandler) Request(f *proxy.Flow) {
	if f.Request.URL.Host != "app.vocabgo.com" {
		return
	}
	if strings.Contains(f.Request.URL.Path, "/student/api/Student/ClassTask/SubmitChoseWord") {
		//Fulsh word storage and wordlist
		words = []WordInfo{}
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

		//Setup progress ui
		progressBar := widget.NewProgressBar()
		completeBox := widget.NewLabel("")
		var wordList string
		if len(dataset.CurrentTask.WordList) > 8 {
			wordList = fmt.Sprintln(dataset.CurrentTask.WordList[:8]) + "..."
		} else {
			wordList = fmt.Sprintln(dataset.CurrentTask.WordList)
		}
		infoBox := widget.NewLabel("New task detected. Gathering chosen words' info:\n" + wordList)
		window.SetContent(container.NewVBox(infoBox, progressBar, completeBox, infoLable))

		for i := 0; i < len(dataset.CurrentTask.WordList); i++ {
			//Show progress
			progressBar.SetValue(float64(i) / float64(len(dataset.CurrentTask.WordList)))
			grabWord(dataset.CurrentTask.WordList[i])
			log.Println("[I] Grabbed word list:" + dataset.CurrentTask.WordList[i])
		}
		progressBar.SetValue(1)
		completeBox.SetText("Complete!")
	}
}

func (c *VocabMasterHandler) Response(f *proxy.Flow) {
	//Judge whether session is from vocabgo task
	if f.Request.URL.Host != "app.vocabgo.com" {
		return
	}
	if !strings.Contains(f.Request.URL.Path, "/student/api/Student/ClassTask/SubmitAnswerAndSave") && !strings.Contains(f.Request.URL.Path, "/student/api/Student/ClassTask/StartAnswer") {
		return
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
		infoLable.SetText("Seems like you are entering an new task!\nPlease wait until progress bar reach 100%.")
	//Choose translation of specific word from a sentence
	case 11:
		var translation string
		var found bool
		stemConverted := strings.ReplaceAll(vocabTask.Stem.Content, "  ", " ")
		for i := 0; i < len(words); i++ {
			for j := 0; j < len(words[i].Content); j++ {
				for k := 0; k < len(words[i].Content[j].ExampleEnglish); k++ {
					if words[i].Content[j].ExampleEnglish[k] == stemConverted {
						translation = words[i].Content[j].Meaning
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if found {
				break
			}
		}

		var contentIndex int
		for i := 0; i < len(vocabTask.Options); i++ {
			regex := regexp.MustCompile(`（.*?）`)
			vocabTask.Options[i].Content = string(regex.ReplaceAll([]byte(vocabTask.Options[i].Content), []byte("")))
			if compareTranslation(translation, vocabTask.Options[i].Content) {
				contentIndex = i
				break
			}
		}

		//UI
		infoLable.SetText("Hey! The anwser is tagged out.")

		//Tag out the correct answer
		regex := regexp.MustCompile(`（.*?）`)
		newJSON := string(rawDecodedString)
		newJSON = string(regex.ReplaceAll([]byte(newJSON), []byte("")))
		newJSON = strings.Replace(newJSON, vocabTask.Options[contentIndex].Content, "-> "+vocabTask.Options[contentIndex].Content+" <-", 1)
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

	//Choose word from voice
	case 22:
		var contentIndex int
		var found bool
		for i := 0; i < len(words); i++ {
			if words[i].Word == vocabTask.Stem.Content {
				for j := 0; j < len(words[i].Content); j++ {
					for k := 0; k < len(vocabTask.Options); k++ {
						regex := regexp.MustCompile(`（.*?）`)
						vocabTask.Options[k].Content = string(regex.ReplaceAll([]byte(vocabTask.Options[k].Content), []byte("")))
						if compareTranslation(vocabTask.Options[k].Content, words[i].Content[j].Meaning) {
							contentIndex = k
							found = true
							break
						}
					}
					if found {
						break
					}
				}
				break
			}
		}

		//UI
		infoLable.SetText("Hey! The anwser is tagged out.")

		//Tag out the correct answer
		regex := regexp.MustCompile(`（.*?）`)
		newJSON := string(rawDecodedString)
		newJSON = string(regex.ReplaceAll([]byte(newJSON), []byte("")))
		newJSON = strings.Replace(newJSON, vocabTask.Options[contentIndex].Content, "-> "+vocabTask.Options[contentIndex].Content+" <-", 1)
		repackedBase64 := base64.StdEncoding.EncodeToString([]byte(newJSON))
		vocabRawJSON.Data = JSONSalt + repackedBase64
		body, _ := json.Marshal(vocabRawJSON)
		var b bytes.Buffer
		br := brotli.NewWriter(&b)
		br.Write(body)
		br.Flush()
		br.Close()
		f.Response.Body = b.Bytes()

	//Choose word pair
	case 31:
		var tag []string
		for i := 0; i < len(vocabTask.Stem.Remark); i++ {
			for j := 0; j < len(vocabTask.Options); j++ {
				if strings.Contains(vocabTask.Stem.Remark[i].SenMarked, vocabTask.Options[j].Content) {
					tag = append(tag, vocabTask.Options[j].Content)
				}
			}
		}

		infoLable.SetText("The anwser is:\n" + fmt.Sprintln(tag))

		//Show answer in the UI

		//Changing the value of word will make it unclickable. This method is invalid.
		//newJSON := string(rawDecodedString)
		//for i := 0; i < len(tag); i++ {
		//	newJSON = strings.Replace(newJSON, `"content":"`+tag[i]+`"`, `"content":"`+"-> "+tag[i]+" <-"+`"`, 1)
		//}
		//repackedBase64 := base64.StdEncoding.EncodeToString([]byte(newJSON))
		//vocabRawJSON.Data = JSONSalt + repackedBase64
		//body, _ := json.Marshal(vocabRawJSON)
		//var b bytes.Buffer
		//br := brotli.NewWriter(&b)
		//br.Write(body)
		//br.Flush()
		//br.Close()
		//f.Response.Body = b.Bytes()

	//Organize word pieces
	case 32:
		regexFind := regexp.MustCompile(`"remark":".*?"`)
		raw := regexFind.FindString(string(rawDecodedString))
		word := raw[10 : len(raw)-1]
		var tag string
		var found bool
		for i := 0; i < len(words); i++ {
			for j := 0; j < len(words[i].Content); j++ {
				for k := 0; k < len(words[i].Content[j].Usage); k++ {
					if strings.Contains(words[i].Content[j].Usage[k], word) {
						tag = words[i].Content[j].Usage[k]
						break
					}
				}
				if found {
					break
				}
			}
			if found {
				break
			}
		}

		//UI
		infoLable.SetText("Hey! The anwser is printed out.")

		//Change the hint to the correct answer
		newJSON := strings.Replace(string(rawDecodedString), word, tag, 1)
		repackedBase64 := base64.StdEncoding.EncodeToString([]byte(newJSON))
		vocabRawJSON.Data = JSONSalt + repackedBase64
		body, _ := json.Marshal(vocabRawJSON)
		var b bytes.Buffer
		br := brotli.NewWriter(&b)
		br.Write(body)
		br.Flush()
		br.Close()
		f.Response.Body = b.Bytes()

	//Write words from first chars
	case 51:
		regexFind := regexp.MustCompile(`"remark":".*?"`)
		raw := regexFind.FindString(string(rawDecodedString))
		word := raw[10 : len(raw)-1]
		var tag string
		var found bool
		for i := 0; i < len(words); i++ {
			for j := 0; j < len(words[i].Content); j++ {
				for k := 0; k < len(words[i].Content[j].Usage); k++ {
					if strings.Contains(words[i].Content[j].Usage[k], word) {
						tag = words[i].Content[j].Usage[k]
						break
					}
				}
				if found {
					break
				}
			}
			if found {
				break
			}
		}

		//UI
		infoLable.SetText("Hey! The anwser is printed out.")

		//Change the hint to the correct answer
		newJSON := strings.Replace(string(rawDecodedString), word, tag, 1)
		repackedBase64 := base64.StdEncoding.EncodeToString([]byte(newJSON))
		vocabRawJSON.Data = JSONSalt + repackedBase64
		body, _ := json.Marshal(vocabRawJSON)
		var b bytes.Buffer
		br := brotli.NewWriter(&b)
		br.Write(body)
		br.Flush()
		br.Close()
		f.Response.Body = b.Bytes()

	default:
		infoLable.SetText("WTF? This is not right. This might be a bug.\n")
	}
}

func compareTranslation(str1 string, str2 string) bool {
	//Compare length first
	if len(str1) != len(str2) {
		return false
	}

	//Delete the classification of current word
	str1 = strings.Split(str1, " ")[1]
	str2 = strings.Split(str2, " ")[1]

	//Split
	str1 = strings.ReplaceAll(str1, "；", "，")
	str2 = strings.ReplaceAll(str2, "；", "，")
	str1split := strings.Split(str1, "，")
	str2split := strings.Split(str2, "，")

	//Compare split length
	if len(str1split) != len(str2split) {
		return false
	}

	//Sort and compare content
	sort.Strings(str1split)
	sort.Strings(str2split)

	for i := 0; i < len(str1split); i++ {
		if str1split[i] != str2split[i] {
			return false
		}
	}

	return true
}
