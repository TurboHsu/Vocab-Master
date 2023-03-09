package answer

type WordInfo struct {
	Word    string
	Content []WordInfoContent
}

type WordInfoContent struct {
	Meaning        string
	Usage          []string
	ExampleEnglish []string
}
