package automatic

import (
	"net/http"
)

type TaskInfo struct {
	Timestamp int
	TopicMode string
	AnswerResult int
	OverStatus int
}

type VocabDataset struct {
	CurrentTask
	RequestInfo
}
type CurrentTask struct {
	WordList []string
	TaskSet  string
	TaskID   string
	TaskType string
}
type RequestInfo struct {
	Versions string
	Cookies  []*http.Cookie
	Header   http.Header
}
