package main

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

type ProxyState struct {
	Enabled bool
	Server  string
	Type    string
	Device  string
}

type Platform struct {
	DataDir string
	CertDir string
	Font    string
}
