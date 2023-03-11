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
	TaskID    string
	AppType   string
	ReleaseID string
	TaskType  string
	Versions  string
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Data      struct {
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

type StartAnswerResponseStruct struct {
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

type SubmitAnswerAndSaveStruct struct {
	TopicCode   string `json:"topic_code"`
	TimeSpent   int    `json:"time_spent"`
	OptImgW     int    `json:"opt_img_w"`
	OptFontSize int    `json:"opt_font_size"`
	OptFontC    string `json:"opt_font_c"`
	ItImgW      int    `json:"it_img_w"`
	ItFontSize  int    `json:"it_font_size"`
	Timestamp   int64  `json:"timestamp"`
	Version     string `json:"version"`
	Sign        string `json:"sign"`
	AppType     int    `json:"app_type"`
}
