package grab

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/TurboHsu/Vocab-Master/answer"
)

var FetchIdentityTrigger *widget.Check
var FetchProgressBar *widget.ProgressBar

func GenerateNewWindow(app *fyne.App) (window fyne.Window) {
	window = (*app).NewWindow("Database console")
	FetchIdentityTrigger = widget.NewCheck("Fetch identify", func(b bool) {
		FetchIdentity = b
	})
	FetchIdentityTrigger.SetChecked(FetchIdentity)

	fetchCourseIDEntry := widget.NewEntry()
	fetchCourseIDLabel := widget.NewLabel("Course ID:")

	FetchProgressBar = widget.NewProgressBar()
	fetchStartBtn := widget.NewButton("Start Fetching", func() {
		// Check the validation
		if !DatasetValid {
			dialog.NewInformation("Hmm", "Cookie and headers not found.\nEnable identity fetcher to do that.", window)
			return
		}
		// Clear wordlist
		answer.WordList = []answer.WordInfo{}

		// Get all lists
		url := fmt.Sprintf("https://app.vocabgo.com/student/api/Student/StudyTask/List?course_id=%s&timestamp=%d&version=1.2.0&app_type=1",
			fetchCourseIDEntry.Text,
			time.Now().UnixMilli(),
		)
		resp := doGET(url)
		var courseListData ListInfoJson
		if err := json.Unmarshal([]byte(resp), &courseListData); err != nil {
			dialog.NewInformation("Error", "Failed to parse course list data.", window)
			return
		}

		// Grab words in lists recursively
		for i := 0; i < len(courseListData.Data.TaskList); i++ {
			// Get words in the list
			url = fmt.Sprintf("https://app.vocabgo.com/student/api/Student/StudyTask/Info?task_id=-1&course_id=%s&list_id=%s&timestamp=%d&version=1.2.0&app_type=1",
				fetchCourseIDEntry.Text,
				courseListData.Data.TaskList[i].ListID,
				time.Now().UnixMilli(),
			)
			resp := doGET(url)
			var taskInfo TaskInfoJson
			if err := json.Unmarshal([]byte(resp), &taskInfo); err != nil {
				dialog.NewInformation("Error", "Failed to parse task info data.", window)
				return
			}
			// Set ID
			Dataset.CurrentTask.TaskID = taskInfo.Data.CourseID
			Dataset.CurrentTask.TaskSet = taskInfo.Data.ListID
			Dataset.RequestInfo.Versions = "1.2.0"
			// Grab all words
			for j := 0; j < len(taskInfo.Data.WordList); j++ {
				log.Println("[I] Grabbing word", taskInfo.Data.WordList[j].Word)
				GrabWord(taskInfo.Data.WordList[j].Word, &Dataset, 50+rand.Intn(150))
			}
			// Set progress
			FetchProgressBar.SetValue(float64(i) / float64(len(courseListData.Data.TaskList)))
		}
		FetchProgressBar.SetValue(1)
		IsDatabaseLoaded = true
	})
	saveDatabaseBtn := widget.NewButton("Save Database", func() {
		// Summon data
		data, _ := json.Marshal(answer.WordList)
		dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if writer == nil {
				return
			}
			defer writer.Close()
			// Write to file
			writer.Write(data)
		}, window).Show()
	})
	loadDatabaseBtn := widget.NewButton("Load Database", func() {
		dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if reader == nil {
				return
			}
			defer reader.Close()
			// Read from file
			data, _ := io.ReadAll(reader)
			json.Unmarshal(data, &answer.WordList)
		}, window).Show()
		IsDatabaseLoaded = true
	})

	infoLabel := widget.NewLabel("Grabs word, save or load from database.\nTip: KEEP IDENTITY FETCHER OPEN!\nClick somewhere in the vocabgo window constantly to prevent rate limit!")

	window.SetContent(
		container.NewVBox(
			infoLabel,
			FetchIdentityTrigger,
			container.NewHBox(
				fetchCourseIDLabel,
				fetchCourseIDEntry,
			),
			FetchProgressBar,
			container.NewHBox(
				fetchStartBtn,
				loadDatabaseBtn,
				saveDatabaseBtn,
			),
		),
	)

	return
}
