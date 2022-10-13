package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/lqqyt2423/go-mitmproxy/proxy"
)

type VocabMasterHandler struct {
	proxy.BaseAddon
}

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

		fmt.Println("Get!")
		for i := 0; i < len(dataset.CurrentTask.WordList); i++ {
			grabWord(dataset.CurrentTask.WordList[i])
			log.Println("[I] Grabbed word list:" + dataset.CurrentTask.WordList[i])
		}
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

	//TODO:
	/*
		2. Invoke tips into vocab
		3. Improve some experience
	*/

	//Get decoded content
	rawByte, _ := f.Response.DecodedBody()

	//Okay! Let's decode raw json
	var vocabRawJSON VocabRawJSONStruct
	json.Unmarshal(rawByte, &vocabRawJSON)

	//Judge whether is the last task
	if vocabRawJSON.Msg != "success" {
		return
	}

	//Let's get the insider base64-encoded info
	rawDecodedString, err := base64.StdEncoding.DecodeString(vocabRawJSON.Data[32:])
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
	var displayText string
	var done bool = false
	switch vocabTask.TopicMode {
	case 0:
		var expression string
		for i := 0; i < len(vocabTask.Options); i++ {
			expression += vocabTask.Options[i].Content + ", "
			done = true
		}
		wordCache = append(wordCache, []string{vocabTask.Stem.Content, expression})
		displayText = "Task mode 1, gathering word list: " + vocabTask.Stem.Content
	case 11:
		regexFind := regexp.MustCompile("{.*}")
		word := regexFind.FindString(vocabTask.Stem.Content)
		displayText = "Task mode 2, word:[" + word[1:len(word)-1] + "]\n"
		for i := 0; i < len(wordCache); i++ {
			if strings.Contains(strings.ToLower(word[1:len(word)-1]), wordCache[i][0]) {
				displayText += wordCache[i][1]
				done = true
			}
		}
	case 22:
		displayText = "Task mode 3, word:[" + vocabTask.Stem.Content + "]\n"
		for i := 0; i < len(wordCache); i++ {
			if wordCache[i][0] == vocabTask.Stem.Content {
				displayText += wordCache[i][1]
				done = true
				break
			}
		}
	case 31:
		displayText = "Task mode 4.\n"
		for i := 0; i < len(vocabTask.Stem.Remark); i++ {
			displayText += vocabTask.Stem.Remark[i].SenMarked + " " + vocabTask.Stem.Remark[i].SenCN + "\n"
			done = true
		}
	case 32:
		displayText = "Task mode 5. "
		regexFind := regexp.MustCompile(`"remark":".*?"`)
		raw := regexFind.FindString(string(rawDecodedString))
		word := raw[10 : len(raw)-1]
		displayText += word + "\n"
		for i := 0; i < len(wordCache); i++ {
			var opt bool = false
			for j := 0; j < len(word); j += 3 {
				if strings.Contains(wordCache[i][1], string([]byte{word[j], word[j+1], word[j+2]})) {
					opt = true
				}
			}
			if opt {
				displayText += wordCache[i][0] + " " + wordCache[i][1] + "\n"
				done = true
			}
		}
	case 51:
		displayText = "Task mode 6.\n"
		for i := 0; i < len(wordCache); i++ {
			if string(wordCache[i][0][:len(vocabTask.WTip)]) == vocabTask.WTip {
				displayText += wordCache[i][0] + " " + wordCache[i][1] + "\n"
				done = true
			}

			//if strings.Contains(wordCache[i][0], vocabTask.WTip) {
			//}
		}
	default:
		displayText = "WTF? This is not right. This might be a bug.\n"
		displayText += string(rawDecodedString)
	}
	if !done {
		displayText += "\nFinding word failed. Displaying the full vocabulary:\n"
		for i := 0; i < len(wordCache); i++ {
			displayText += wordCache[i][0] + " " + wordCache[i][1] + "\n"
		}
	}
	textBox.SetText(displayText)

}
