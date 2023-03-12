package automatic

import (
	"encoding/json"
	"fmt"

	"github.com/lqqyt2423/go-mitmproxy/proxy"
)

var taskDetail TaskDetail

// TaskDetailProcessor is a function that processes the task detail JSON response.
func TaskDetailProcessor(f *proxy.Flow) {
	// clear all data
	clearTrigger()

	//Update Cookie and Header
	dataset.RequestInfo.Cookies = f.Request.Raw().Cookies()
	dataset.RequestInfo.Header = f.Request.Raw().Header

	rawJSON, _ := f.Response.DecodedBody()
	// Start processing
	json.Unmarshal(rawJSON, &taskDetail)

	// Parse query
	taskDetail.TaskID = f.Request.URL.Query().Get("task_id")
	taskDetail.TaskType = f.Request.URL.Query().Get("task_type")
	taskDetail.Versions = f.Request.URL.Query().Get("version")
	taskDetail.ReleaseID = f.Request.URL.Query().Get("release_id")
	taskDetail.AppType = f.Request.URL.Query().Get("app_type")

	// Append some info to dataset
	dataset.CurrentTask.TaskSet = taskDetail.Data.WordList[0].ListID
	dataset.CurrentTask.TaskID = taskDetail.Data.WordList[0].CourseID
	dataset.RequestInfo.Versions = taskDetail.Versions
	taskDetail.TaskType = fmt.Sprint(taskDetail.Data.TaskType)

	// Process pending words
	for _, word := range taskDetail.Data.WordList {
		if doGetFullScore && word.Score < 10.0 {
			pendingWord = append(pendingWord, word.Word)
			continue
		}
		if word.Progress < 100 {
			pendingWord = append(pendingWord, word.Word)
		}
	}

	// Trigger valid
	valid = true
	confirmed = false
	list.Refresh()
}

// This clears everything.
func clearTrigger() {
	taskDetail = TaskDetail{}
	pendingWord = []string{}
}
