package automatic

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/TurboHsu/Vocab-Master/answer"
	"github.com/TurboHsu/Vocab-Master/grab"
)

/*
	DECEPRECATED CODEBASE

	IKD WHY IT DOES NOT WORK
*/

var dataset grab.VocabDataset
var Enabled bool
var pendingWord []string

func FetchDataset(data grab.VocabDataset) {
	dataset = data
}

func UpdateIdentity(data grab.VocabDataset) {
	dataset.RequestInfo = data.RequestInfo
}

func startAutomation() {
	taskPath := "ClassTask" //TODO: This field is not implemented, may get from url from mitm hook.

	// Size of each group
	groupSize := 5

	// Calculate the number of groups needed
	numGroups := (len(pendingWord) + groupSize - 1) / groupSize

	// Create a slice to hold the groups
	groups := make([][]string, numGroups)

	// Iterate over the groups and assign the elements from the original slice
	for i := range groups {
		start := i * groupSize
		end := start + groupSize

		if end > len(pendingWord) {
			end = len(pendingWord)
		}

		groups[i] = pendingWord[start:end]
	}

	for _, wordGroup := range groups {

		// Grab all words
		info.SetText("Grabbing words...")
		progressBar.SetValue(0)
		progressBar.Show()
		for i := 0; i < len(wordGroup); i++ {
			progressBar.SetValue(float64(i) / float64(len(wordGroup)))
			// Grab word
			grab.GrabWord(wordGroup[i], (*grab.VocabDataset)(&dataset), 50)
			info.SetText("Grabbing word: " + wordGroup[i])
			log.Println("Grabbing word: ", wordGroup[i])
		}
		progressBar.SetValue(1)
		info.SetText("Grabbing words finished.")

		// Submit the words need to be done.
		// Summon the json
		var wordListSubmitJSON string
		for i := 0; i < len(wordGroup); i++ {
			wordListSubmitJSON += `"` + wordGroup[i] + `",`
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
		url := fmt.Sprintf("https://app.vocabgo.com/student/api/Student/%s/SubmitChoseWord",
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

		progressBar.SetValue(0.5)

		var topicCode string
		var rawDecodedString string = string(startAnswerResponseDesaltedDecoded)
		timeDelay, _ := strconv.Atoi(timeDelayEntry.Text)
	MainLoop:
		for {
			// Delay for specific time
			info.SetText("Delaying for " + timeDelayEntry.Text + " seconds...")
			time.Sleep(time.Duration(timeDelay) * time.Second)

			switch vocabTaskData.TopicMode {
			// This is letting u read through something.
			case 0:
				// Get the next TopicCode and do nothing.
				log.Println("Doing Topic mode 0.")
				topicCode = vocabTaskData.TopicCode
			// Choose some translation for the word.
			case 11:
				log.Println("Doing Topic mode 11.")
				ans := answer.FindAnswer(11, answer.VocabTaskStruct(vocabTaskData), "")
				// Submit the ans or guess one
				if !ans.Found {
					ans.Index = append(ans.Index, rand.Intn(4))
				}
				// Summon JSON
				verifyJSON := fmt.Sprintf(`{"answer":%d,"topic_code":"%s","timestamp":%d,"version":"%s","sign":"%s","app_type":%s}`,
					ans.Index[0],
					topicCode,
					time.Now().UnixMilli(),
					taskDetail.Versions,
					summonSign(),
					taskDetail.AppType,
				)

				// Verify the answer
				resp = doPOST(fmt.Sprintf("https://app.vocabgo.com/student/api/Student/%s/VerifyAnswer", taskPath), verifyJSON)
				if !strings.Contains(resp, `"code":1`) {
					info.SetText("Failed to verify the answer.")
					log.Println("Failed to verify the answer, topic mode 11.")
					break MainLoop
				}
				// Parse the JSON
				var verifyAnswerResponse StartAnswerResponseStruct
				json.Unmarshal([]byte(resp), &verifyAnswerResponse)
				// Get rid of salt and parse
				_, verifyAnswerResponseDesalted := splitSalt(verifyAnswerResponse.Data)
				// Decode the base64
				verifyAnswerResponseDesaltedDecoded, _ := base64.StdEncoding.DecodeString(verifyAnswerResponseDesalted)
				// Unmarshal the data using json
				var verifyAnswerResponseDesaltedDecodedData VerifyAnswerResponseStruct
				json.Unmarshal(verifyAnswerResponseDesaltedDecoded, &verifyAnswerResponseDesaltedDecodedData)
				// Update topicCode
				topicCode = vocabTaskData.TopicCode
				// Check whether the answer is correct
				if verifyAnswerResponseDesaltedDecodedData.AnswerResult != 1 {
					// Its incorrect. Lets guess another one.
					ans.Index[0] = rand.Intn(4)
					// Summon JSON
					verifyJSON := fmt.Sprintf(`{"answer":%d,"topic_code":"%s","timestamp":%d,"version":"%s","sign":"%s","app_type":%s}`,
						ans.Index[0],
						topicCode,
						time.Now().UnixMilli(),
						taskDetail.Versions,
						summonSign(),
						taskDetail.AppType,
					)
					resp = doPOST(fmt.Sprintf("https://app.vocabgo.com/student/api/Student/%s/VerifyAnswer", taskPath), verifyJSON)
					if !strings.Contains(resp, `"code":1`) {
						info.SetText("Failed to verify the answer.")
						log.Println("Failed to verify the answer after retrying, topic mode 11.")
						break MainLoop
					}
					// Parse the JSON
					var verifyAnswerResponse StartAnswerResponseStruct
					json.Unmarshal([]byte(resp), &verifyAnswerResponse)
					// Get rid of salt and parse
					_, verifyAnswerResponseDesalted := splitSalt(verifyAnswerResponse.Data)
					// Decode the base64
					verifyAnswerResponseDesaltedDecoded, _ := base64.StdEncoding.DecodeString(verifyAnswerResponseDesalted)
					// Unmarshal the data using json
					var verifyAnswerResponseDesaltedDecodedData VerifyAnswerResponseStruct
					json.Unmarshal(verifyAnswerResponseDesaltedDecoded, &verifyAnswerResponseDesaltedDecodedData)
					// Update topicCode
					topicCode = vocabTaskData.TopicCode
					// If its wrong again then idk what to do.
				}
			// Choose translation from voice
			case 22:
				log.Println("Doing Topic mode 22.")
				ans := answer.FindAnswer(22, answer.VocabTaskStruct(vocabTaskData), "")
				// Submit the ans or guess one
				if !ans.Found {
					ans.Index = append(ans.Index, rand.Intn(4))
				}
				// Summon JSON
				verifyJSON := fmt.Sprintf(`{"answer":%d,"topic_code":"%s","timestamp":%d,"version":"%s","sign":"%s","app_type":%s}`,
					ans.Index[0],
					topicCode,
					time.Now().UnixMilli(),
					taskDetail.Versions,
					summonSign(),
					taskDetail.AppType,
				)

				// Verify the answer
				resp = doPOST(fmt.Sprintf("https://app.vocabgo.com/student/api/Student/%s/VerifyAnswer", taskPath), verifyJSON)
				if !strings.Contains(resp, `"code":1`) {
					info.SetText("Failed to verify the answer.")
					log.Println("Failed to verify the answer, topic mode 22.")
					break MainLoop
				}
				// Parse the JSON
				var verifyAnswerResponse StartAnswerResponseStruct
				json.Unmarshal([]byte(resp), &verifyAnswerResponse)
				// Get rid of salt and parse
				_, verifyAnswerResponseDesalted := splitSalt(verifyAnswerResponse.Data)
				// Decode the base64
				verifyAnswerResponseDesaltedDecoded, _ := base64.StdEncoding.DecodeString(verifyAnswerResponseDesalted)
				// Unmarshal the data using json
				var verifyAnswerResponseDesaltedDecodedData VerifyAnswerResponseStruct
				json.Unmarshal(verifyAnswerResponseDesaltedDecoded, &verifyAnswerResponseDesaltedDecodedData)
				// Update topicCode
				topicCode = vocabTaskData.TopicCode
				// Check whether the answer is correct
				if verifyAnswerResponseDesaltedDecodedData.AnswerResult != 1 {
					// Its incorrect. Lets guess another one.
					ans.Index[0] = rand.Intn(4)
					// Summon JSON
					verifyJSON := fmt.Sprintf(`{"answer":%d,"topic_code":"%s","timestamp":%d,"version":"%s","sign":"%s","app_type":%s}`,
						ans.Index[0],
						topicCode,
						time.Now().UnixMilli(),
						taskDetail.Versions,
						summonSign(),
						taskDetail.AppType,
					)
					resp = doPOST(fmt.Sprintf("https://app.vocabgo.com/student/api/Student/%s/VerifyAnswer", taskPath), verifyJSON)
					if !strings.Contains(resp, `"code":1`) {
						info.SetText("Failed to verify the answer.")
						log.Println("Failed to verify the answer after retrying, topic mode 22.")
						break MainLoop
					}
					// Parse the JSON
					var verifyAnswerResponse StartAnswerResponseStruct
					json.Unmarshal([]byte(resp), &verifyAnswerResponse)
					// Get rid of salt and parse
					_, verifyAnswerResponseDesalted := splitSalt(verifyAnswerResponse.Data)
					// Decode the base64
					verifyAnswerResponseDesaltedDecoded, _ := base64.StdEncoding.DecodeString(verifyAnswerResponseDesalted)
					// Unmarshal the data using json
					var verifyAnswerResponseDesaltedDecodedData VerifyAnswerResponseStruct
					json.Unmarshal(verifyAnswerResponseDesaltedDecoded, &verifyAnswerResponseDesaltedDecodedData)
					// Update topicCode
					topicCode = vocabTaskData.TopicCode
					// If its wrong again then idk what to do.
				}
			// Choose word pair
			case 31:
				log.Println("Doing Topic mode 31.")
				ans := answer.FindAnswer(31, answer.VocabTaskStruct(vocabTaskData), "")
				// If does not found then idk wtf
				if !ans.Found {
					times := rand.Intn(len(vocabTaskData.Options) + 1)
					var rangeIndex []int
					for i := 0; i < len(vocabTaskData.Options); i++ {
						rangeIndex = append(rangeIndex, i)
					}
					// Shuffle the range index
					rand.Shuffle(len(rangeIndex), func(i, j int) { rangeIndex[i], rangeIndex[j] = rangeIndex[j], rangeIndex[i] })
					for i := 0; i < times; i++ {
						ans.Index = append(ans.Index, rangeIndex[i])
					}
				}

				for i := 0; i < len(ans.Index); i++ {
					// Summon JSON
					verifyJSON := fmt.Sprintf(`{"answer":%d,"topic_code":"%s","timestamp":%d,"version":"%s","sign":"%s","app_type":%s}`,
						ans.Index[i],
						topicCode,
						time.Now().UnixMilli(),
						taskDetail.Versions,
						summonSign(),
						taskDetail.AppType,
					)

					// Verify the answer
					resp = doPOST(fmt.Sprintf("https://app.vocabgo.com/student/api/Student/%s/VerifyAnswer", taskPath), verifyJSON)
					if !strings.Contains(resp, `"code":1`) {
						info.SetText("Failed to verify the answer.")
						log.Println("Failed to verify the answer, topic mode 31.")
						break MainLoop
					}
					// Parse the JSON
					var verifyAnswerResponse StartAnswerResponseStruct
					json.Unmarshal([]byte(resp), &verifyAnswerResponse)
					// Get rid of salt and parse
					_, verifyAnswerResponseDesalted := splitSalt(verifyAnswerResponse.Data)
					// Decode the base64
					verifyAnswerResponseDesaltedDecoded, _ := base64.StdEncoding.DecodeString(verifyAnswerResponseDesalted)
					// Unmarshal the data using json
					var verifyAnswerResponseDesaltedDecodedData VerifyAnswerResponseStruct
					json.Unmarshal(verifyAnswerResponseDesaltedDecoded, &verifyAnswerResponseDesaltedDecodedData)
					// Update topicCode
					topicCode = vocabTaskData.TopicCode
					// Check whether its done
					if verifyAnswerResponseDesaltedDecodedData.OverStatus == 1 {
						break
					}
				}

				// overStatus = 1 means its over. 2 means its not over
				// Submit the ans recursively

				// Verify the answer
			// Organize word pieces
			case 32:
				log.Println("Doing Topic mode 32.")
				ans := answer.FindAnswer(32, answer.VocabTaskStruct(vocabTaskData), string(rawDecodedString))
				// If not found then guess one
				var wordsToSubmit []string
				if !ans.Found {
					// Get all options
					for _, option := range vocabTaskData.Options {
						wordsToSubmit = append(wordsToSubmit, option.Content)
					}
					// Shuffle the options
					rand.Shuffle(len(wordsToSubmit), func(i, j int) { wordsToSubmit[i], wordsToSubmit[j] = wordsToSubmit[j], wordsToSubmit[i] })
				} else {
					// Append the found answer in the correct order
					for _, index := range ans.Index {
						wordsToSubmit = append(wordsToSubmit, vocabTaskData.Options[index].Content)
					}
				}
				// Verify the answer
				// Summon answer string
				var answerString string
				for _, word := range wordsToSubmit {
					answerString += word + ","
				}

				// Summon JSON
				verifyJSON := fmt.Sprintf(`{"answer":"%s","topic_code":"%s","timestamp":%d,"version":"%s","sign":"%s","app_type":%s}`,
					answerString[:len(answerString)-1],
					topicCode,
					time.Now().UnixMilli(),
					taskDetail.Versions,
					summonSign(),
					taskDetail.AppType,
				)
				// Verify the answer
				resp = doPOST(fmt.Sprintf("https://app.vocabgo.com/student/api/Student/%s/VerifyAnswer", taskPath), verifyJSON)
				if !strings.Contains(resp, `"code":1`) {
					info.SetText("Failed to verify the answer.")
					log.Println("Failed to verify the answer, topic mode 32.")
					break MainLoop
				}
				// Parse the JSON
				var verifyAnswerResponse StartAnswerResponseStruct
				json.Unmarshal([]byte(resp), &verifyAnswerResponse)
				// Get rid of salt and parse
				_, verifyAnswerResponseDesalted := splitSalt(verifyAnswerResponse.Data)
				// Decode the base64
				verifyAnswerResponseDesaltedDecoded, _ := base64.StdEncoding.DecodeString(verifyAnswerResponseDesalted)
				// Unmarshal the data using json
				var verifyAnswerResponseDesaltedDecodedData VerifyAnswerResponseStruct
				json.Unmarshal(verifyAnswerResponseDesaltedDecoded, &verifyAnswerResponseDesaltedDecodedData)
				// Update topicCode
				topicCode = vocabTaskData.TopicCode
				// Check whether the answer is correct
				if verifyAnswerResponseDesaltedDecodedData.AnswerResult != 1 {
					// Its incorrect. Lets guess another one.
					wordsToSubmit = []string{}
					// Get all options
					for _, option := range vocabTaskData.Options {
						wordsToSubmit = append(wordsToSubmit, option.Content)
					}
					// Shuffle the options
					rand.Shuffle(len(wordsToSubmit), func(i, j int) { wordsToSubmit[i], wordsToSubmit[j] = wordsToSubmit[j], wordsToSubmit[i] })

					// Summon answer string
					var answerString string
					for _, word := range wordsToSubmit {
						answerString += word + ","
					}

					// Summon JSON
					verifyJSON := fmt.Sprintf(`{"answer":"%s","topic_code":"%s","timestamp":%d,"version":"%s","sign":"%s","app_type":%s}`,
						answerString[:len(answerString)-1],
						topicCode,
						time.Now().UnixMilli(),
						taskDetail.Versions,
						summonSign(),
						taskDetail.AppType,
					)
					resp = doPOST(fmt.Sprintf("https://app.vocabgo.com/student/api/Student/%s/VerifyAnswer", taskPath), verifyJSON)
					if !strings.Contains(resp, `"code":1`) {
						info.SetText("Failed to verify the answer.")
						log.Println("Failed to verify the answer, after retrying, topic mode 32.")
						break MainLoop
					}
					// Parse the JSON
					var verifyAnswerResponse StartAnswerResponseStruct
					json.Unmarshal([]byte(resp), &verifyAnswerResponse)
					// Get rid of salt and parse
					_, verifyAnswerResponseDesalted := splitSalt(verifyAnswerResponse.Data)
					// Decode the base64
					verifyAnswerResponseDesaltedDecoded, _ := base64.StdEncoding.DecodeString(verifyAnswerResponseDesalted)
					// Unmarshal the data using json
					var verifyAnswerResponseDesaltedDecodedData VerifyAnswerResponseStruct
					json.Unmarshal(verifyAnswerResponseDesaltedDecoded, &verifyAnswerResponseDesaltedDecodedData)
					// Update topicCode
					topicCode = vocabTaskData.TopicCode
					// If its wrong again then idk what to do.
				}

			// Fill in blanks
			case 51:
				log.Println("Doing Topic mode 51.")
				ans := answer.FindAnswer(51, answer.VocabTaskStruct(vocabTaskData), string(rawDecodedString))
				// If really cant find one, fuck it.
				if !ans.Found && !ans.Detail.Uncertain {
					ans.Detail.Word = answer.WordList[rand.Intn(len(answer.WordList))].Word
				}
				// Summon JSON
				verifyJSON := fmt.Sprintf(`{"answer":"%s","topic_code":"%s","timestamp":%d,"version":"%s","sign":"%s","app_type":%s}`,
					ans.Detail.Word,
					topicCode,
					time.Now().UnixMilli(),
					taskDetail.Versions,
					summonSign(),
					taskDetail.AppType,
				)

				// Verify the answer
				resp = doPOST(fmt.Sprintf("https://app.vocabgo.com/student/api/Student/%s/VerifyAnswer", taskPath), verifyJSON)
				if !strings.Contains(resp, `"code":1`) {
					info.SetText("Failed to verify the answer.")
					log.Println("Failed to verify the answer, topic mode 51.")
					break MainLoop
				}
				// Parse the JSON
				var verifyAnswerResponse StartAnswerResponseStruct
				json.Unmarshal([]byte(resp), &verifyAnswerResponse)
				// Get rid of salt and parse
				_, verifyAnswerResponseDesalted := splitSalt(verifyAnswerResponse.Data)
				// Decode the base64
				verifyAnswerResponseDesaltedDecoded, _ := base64.StdEncoding.DecodeString(verifyAnswerResponseDesalted)
				// Unmarshal the data using json
				var verifyAnswerResponseDesaltedDecodedData VerifyAnswerResponseStruct
				json.Unmarshal(verifyAnswerResponseDesaltedDecoded, &verifyAnswerResponseDesaltedDecodedData)
				// Update topicCode
				topicCode = vocabTaskData.TopicCode
				// Check whether the answer is correct
				if verifyAnswerResponseDesaltedDecodedData.AnswerResult != 1 {
					// Its incorrect. Lets guess another one.
					ans.Detail.Word = answer.WordList[rand.Intn(len(answer.WordList))].Word
					// Summon JSON
					verifyJSON := fmt.Sprintf(`{"answer":"%s","topic_code":"%s","timestamp":%d,"version":"%s","sign":"%s","app_type":%s}`,
						ans.Detail.Word,
						topicCode,
						time.Now().UnixMilli(),
						taskDetail.Versions,
						summonSign(),
						taskDetail.AppType,
					)
					resp = doPOST(fmt.Sprintf("https://app.vocabgo.com/student/api/Student/%s/VerifyAnswer", taskPath), verifyJSON)
					if !strings.Contains(resp, `"code":1`) {
						info.SetText("Failed to verify the answer.")
						log.Println("Failed to verify the answer, after retrying, topic mode 51.")
						break MainLoop
					}
					// Parse the JSON
					var verifyAnswerResponse StartAnswerResponseStruct
					json.Unmarshal([]byte(resp), &verifyAnswerResponse)
					// Get rid of salt and parse
					_, verifyAnswerResponseDesalted := splitSalt(verifyAnswerResponse.Data)
					// Decode the base64
					verifyAnswerResponseDesaltedDecoded, _ := base64.StdEncoding.DecodeString(verifyAnswerResponseDesalted)
					// Unmarshal the data using json
					var verifyAnswerResponseDesaltedDecodedData VerifyAnswerResponseStruct
					json.Unmarshal(verifyAnswerResponseDesaltedDecoded, &verifyAnswerResponseDesaltedDecodedData)
					// Update topicCode
					topicCode = vocabTaskData.TopicCode
					// If its wrong again then idk what to do.
				}
			default:
				info.SetText("Unknown topic mode: " + strconv.Itoa(vocabTaskData.TopicMode))
				log.Println("Unknown topic mode: ", vocabTaskData.TopicMode)
				break MainLoop
			}

			time.Sleep(1 * time.Second)

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
			// Update the rawDecodedString
			rawDecodedString = string(submitAnswerAndSaveResponseDesaltedDecoded)
		}
		info.SetText("Automation finished.")
		progressBar.SetValue(1)
	}
}

// The mechanism of signing is not sure, so now its some md5 from timestamp.
func summonSign() string {
	// return md5 of timestamp
	return fmt.Sprintf("%x", md5.Sum([]byte(strconv.Itoa(int(time.Now().Unix())))))
}
