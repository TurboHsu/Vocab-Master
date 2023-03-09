package automatic

import "net/http"

type VocabDataset struct {
	CurrentTask struct {
		WordList []string
		TaskSet  string
		TaskID   string
	}
	RequestInfo struct {
		Versions string
		Cookies  []*http.Cookie
		Header   http.Header
	}
	IsEnabled bool
	IsAuto    bool
}
