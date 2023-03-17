package grab

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
}

type VocabRawJSONStruct struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
	Jv   string `json:"jv"`
	Cv   string `json:"cv"`
}

type WordInfoJSON struct {
	CourseID      string `json:"course_id"`
	ListID        string `json:"list_id"`
	Word          string `json:"word"`
	UpdateVersion string `json:"update_version"`
	Means         []struct {
		Mean   []string `json:"mean"`
		PhInfo struct {
			PhEn    string `json:"ph_en"`
			PhEnURL string `json:"ph_en_url"`
			PhUs    string `json:"ph_us"`
			PhUsURL string `json:"ph_us_url"`
		} `json:"ph_info"`
		Usages []struct {
			Usage        interface{}   `json:"usage"`
			Phrases      []string      `json:"phrases"`
			PhrasesInfos []interface{} `json:"phrases_infos"`
			Examples     []struct {
				SenID         string `json:"sen_id"`
				SenContent    string `json:"sen_content"`
				SenMeanCn     string `json:"sen_mean_cn"`
				SenSource     string `json:"sen_source"`
				SenSourceCode string `json:"sen_source_code"`
				AudioFile     string `json:"audio_file"`
			} `json:"examples"`
		} `json:"usages"`
	} `json:"means,omitempty"`
	Version  string `json:"version"`
	HasAu    int    `json:"has_au"`
	AuAddr   string `json:"au_addr"`
	AuWord   string `json:"au_word"`
	WordInfo struct {
		StoreStatus int `json:"store_status"`
	} `json:"word_info"`
	PhEn    string `json:"ph_en,omitempty"`
	PhUs    string `json:"ph_us,omitempty"`
	Options []struct {
		Content struct {
			Mean   string `json:"mean"`
			PhInfo struct {
				PhEn    string `json:"ph_en"`
				PhEnURL string `json:"ph_en_url"`
				PhUs    string `json:"ph_us"`
				PhUsURL string `json:"ph_us_url"`
			} `json:"ph_info"`
			UsageInfos []struct {
				SenID      string `json:"sen_id"`
				SenContent string `json:"sen_content"`
				SenMeanCn  string `json:"sen_mean_cn"`
				AudioFile  string `json:"audio_file"`
			} `json:"usage_infos"`
			Usage   []string `json:"usage"`
			Example []struct {
				SenID         string `json:"sen_id"`
				SenContent    string `json:"sen_content"`
				SenMeanCn     string `json:"sen_mean_cn"`
				SenSource     string `json:"sen_source"`
				SenSourceCode string `json:"sen_source_code"`
				AudioFile     string `json:"audio_file"`
			} `json:"example"`
		} `json:"content"`
	} `json:"options,omitempty"`
}

type ListInfoJson struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		CourseID string `json:"course_id"`
		TaskList []struct {
			TaskID    int     `json:"task_id"`
			TaskType  int     `json:"task_type"`
			CourseID  string  `json:"course_id"`
			ListID    string  `json:"list_id"`
			TaskName  string  `json:"task_name"`
			Score     float64 `json:"score"`
			Progress  int     `json:"progress"`
			TimeSpent int     `json:"time_spent"`
			Free      int     `json:"free"`
		} `json:"task_list"`
		CourseName     string `json:"course_name"`
		CourseImgURL   string `json:"course_img_url"`
		Progress       int    `json:"progress"`
		TimeSpent      int    `json:"time_spent"`
		Free           int    `json:"free"`
		CourseStatus   int    `json:"course_status"`
		OffInstruction string `json:"off_instruction"`
	} `json:"data"`
	Jv string `json:"jv"`
	Cv string `json:"cv"`
}

type TaskInfoJson struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		TaskID   int    `json:"task_id"`
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
		Grade         int `json:"grade"`
		GradeInfoList []struct {
			Value        int    `json:"value"`
			Text         string `json:"text"`
			TopicModeNum int    `json:"topic_mode_num"`
			BaseTime     int    `json:"base_time"`
			Time         int    `json:"time"`
			Remark       string `json:"remark"`
		} `json:"grade_info_list"`
		Score     float64 `json:"score"`
		Progress  int     `json:"progress"`
		TimeSpent int     `json:"time_spent"`
		AudioAddr string  `json:"audio_addr"`
	} `json:"data"`
	Jv string `json:"jv"`
	Cv string `json:"cv"`
}