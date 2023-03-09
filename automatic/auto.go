package automatic

var dataset VocabDataset

func FetchDataset(data VocabDataset) {
	dataset = data
}

func StartAutomation() {
	// Make the first request.
	//response := doGET(fmt.Sprintf("https://app.vocabgo.com/api/Student/%s/StartAnswer?task_id=%d&",
	//	dataset.CurrentTask.TaskType))
}
