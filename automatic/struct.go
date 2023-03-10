package automatic

type TaskInfo struct {
	Timestamp    int
	TopicMode    string
	AnswerResult int
	OverStatus   int
	ReleaseID    int
	TaskType     string
}

// This is got from JSON response when accessing task page
type TaskDetail struct {
	TaskID string
	TaskType   string
	Versions   string
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	Data       struct {
		TaskID   int    `json:"task_id"`
		TaskType int    `json:"task_type"`
		CourseID string `json:"course_id"`
		ListID   string `json:"list_id"`
		TaskName string `json:"task_name"`
		WordList []struct {
			Progress  int     `json:"progress"`
			Score     float64 `json:"score"`
			TimeSpent int     `json:"time_spent"`
			Status    int     `json:"status"`
			CourseID  string  `json:"course_id"`
			ListID    string  `json:"list_id"`
			Word      string  `json:"word"`
			WordType  int     `json:"word_type"`
			WordZh    string  `json:"word_zh"`
			WordAudio string  `json:"word_audio"`
		} `json:"word_list"`
		Grade     int     `json:"grade"`
		Score     float64 `json:"score"`
		Progress  int     `json:"progress"`
		TimeSpent int     `json:"time_spent"`
		AudioAddr string  `json:"audio_addr"`
	} `json:"data"`
	Jv string `json:"jv"`
	Cv string `json:"cv"`
}
