package grab

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/TurboHsu/Vocab-Master/answer"
)

var DatasetValid bool
var FetchIdentity bool
var IsDatabaseLoaded bool
var Dataset VocabDataset

func GrabWord(word string, dataset *VocabDataset, delay int) {
	time.Sleep(time.Duration(delay) * time.Millisecond)
	url := fmt.Sprintf("https://app.vocabgo.com/student/api/Student/Course/StudyWordInfo?course_id=%s&list_id=%s&word=%s&timestamp=%d&version=%s&app_type=1",
		dataset.CurrentTask.TaskID,
		dataset.CurrentTask.TaskSet,
		word,
		time.Now().UnixMilli(),
		dataset.RequestInfo.Versions)
	req, _ := http.NewRequest("GET", url, nil)
	//Adds headers
	req.Header = dataset.RequestInfo.Header
	//Adds cookies
	for i := 0; i < len(dataset.RequestInfo.Cookies); i++ {
		req.AddCookie(dataset.RequestInfo.Cookies[i])
	}
	//Do request
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("[E]" + err.Error())
	}
	defer response.Body.Close()
	read, _, _ := switchContentEncoding(response)
	raw, _ := io.ReadAll(read)
	var rawData VocabRawJSONStruct
	json.Unmarshal(raw, &rawData)
	if len(rawData.Data) < 32 {
		log.Println("[E] Data length is too short, ERR!")
		return
	}
	rawJSON, _ := base64.StdEncoding.DecodeString(rawData.Data[32:])
	var wordInfoRaw WordInfoJSON
	json.Unmarshal(rawJSON, &wordInfoRaw)

	//Add to word storage
	//Due to Issue #10, we need to check the json data version, the data form is different in different versions
	var wordInfo answer.WordInfo
	wordInfo.Word = word
	if wordInfoRaw.Version == "1" { //Old version like CET4
		for i := 0; i < len(wordInfoRaw.Options); i++ {
			var content answer.WordInfoContent
			regex := regexp.MustCompile(`（.*?）`)
			content.Meaning = string(regex.ReplaceAll([]byte(wordInfoRaw.Options[i].Content.Mean), []byte("")))
			content.Meaning = strings.ReplaceAll(content.Meaning, "\n", "")
			//append usage
			for j := 0; j < len(wordInfoRaw.Options[i].Content.Usage); j++ {
				content.Usage = append(content.Usage, wordInfoRaw.Options[i].Content.Usage[j])
			}
			//append example's english expression
			for j := 0; j < len(wordInfoRaw.Options[i].Content.Example); j++ {
				content.ExampleEnglish = append(content.ExampleEnglish, wordInfoRaw.Options[i].Content.Example[j].SenContent)
			}
			wordInfo.Content = append(wordInfo.Content, content)
		}
	} else if wordInfoRaw.Version == "2" { //New version like JJ_2
		for i := 0; i < len(wordInfoRaw.Means); i++ {
			var content answer.WordInfoContent
			//Getting rid of shit and get meaning
			regex := regexp.MustCompile(`（.*?）`)
			content.Meaning = string(regex.ReplaceAll([]byte(
				fmt.Sprintf("%s %s", wordInfoRaw.Means[i].Mean[0], wordInfoRaw.Means[i].Mean[1]),
			), []byte("")))
			content.Meaning = strings.ReplaceAll(content.Meaning, "\n", "")

			//Append usage
			for j := 0; j < len(wordInfoRaw.Means[i].Usages); j++ {
				//Usage is an array. IDK why some can be NULL :(
				for k := 0; k < len(wordInfoRaw.Means[i].Usages[j].Phrases); k++ { //Here we only collect phrases, which is quite enough
					content.Usage = append(content.Usage, wordInfoRaw.Means[i].Usages[j].Phrases[k])
				}
				for k := 0; k < len(wordInfoRaw.Means[i].Usages[j].Examples); k++ {
					content.ExampleEnglish = append(content.ExampleEnglish, wordInfoRaw.Means[i].Usages[j].Examples[k].SenContent)
				}
			}
			//Append content
			wordInfo.Content = append(wordInfo.Content, content)
		}
	} else {
		log.Printf("[E] Error when grabbing words: version %s unsupported!\n", wordInfoRaw.Version)
	}
	answer.WordList = append(answer.WordList, wordInfo)
}
