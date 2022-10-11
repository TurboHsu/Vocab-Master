package main

import (
	"encoding/base64"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/lqqyt2423/go-mitmproxy/proxy"
)

var wordCache [][]string

type VocabMasterHandler struct {
	proxy.BaseAddon
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

type VocabRawJSONStruct struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
	Jv   string `json:"jv"`
	Cv   string `json:"cv"`
}

type VocabTaskStruct struct {
	TaskID    int `json:"task_id"`
	TaskType  int `json:"task_type"`
	TopicMode int `json:"topic_mode"`
	Stem      struct {
		Content string `json:"content"`
		Remark  []struct {
			SenMarked string `json:"sen_marked"`
			SenCN     string `json:"sen_cn"`
			Relation  string `json:"relation"`
		} `json:"remark"`
		PhUsURL string      `json:"ph_us_url"`
		PhEnURL string      `json:"ph_en_url"`
		AuAddr  interface{} `json:"au_addr"`
	} `json:"stem"`
	Options []struct {
		Content    string      `json:"content"`
		Remark     interface{} `json:"remark"`
		Answer     interface{} `json:"answer"`
		AnswerTag  int         `json:"answer_tag"`
		CheckCode  interface{} `json:"check_code"`
		SubOptions interface{} `json:"sub_options"`
		PhInfo     interface{} `json:"ph_info"`
	} `json:"options"`
	SoundMark    string        `json:"sound_mark"`
	PhEn         string        `json:"ph_en"`
	PhUs         string        `json:"ph_us"`
	AnswerNum    int           `json:"answer_num"`
	ChanceNum    int           `json:"chance_num"`
	TopicDoneNum int           `json:"topic_done_num"`
	TopicTotal   int           `json:"topic_total"`
	WLens        []interface{} `json:"w_lens"`
	WLen         int           `json:"w_len"`
	WTip         string        `json:"w_tip"`
	Tips         string        `json:"tips"`
	WordType     int           `json:"word_type"`
	EnableI      int           `json:"enable_i"`
	EnableII     int           `json:"enable_i_i"`
	EnableIO     int           `json:"enable_i_o"`
	TopicCode    string        `json:"topic_code"`
	AnswerState  int           `json:"answer_state"`
}
