package aid

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/TurboHsu/Vocab-Master/answer"
	"github.com/go-vgo/robotgo"
)

func eventLoop() {
	// parse the wait sec
	sleepms, _ := strconv.Atoi(WaitSecEntry.Text)
	for {
		if !IsEnabled {
			break
		}
		switch answer.CurrentAnswer.TopicMode {
		case 0:
			// just click
			clickPosition("next")
		case 15:
			fallthrough
		case 11:
			fallthrough
		case 21:
			fallthrough
		case 22:
			if len(answer.CurrentAnswer.Index) > 0 {
				// vertical thing
				clickPosition("v" + strconv.Itoa(answer.CurrentAnswer.Index[0]))
				time.Sleep(time.Duration(sleepms) * time.Millisecond)
				clickPosition("next")
			} else {
				// Guessing
				clickPosition("v" + strconv.Itoa(rand.Intn(4)))
				time.Sleep(time.Duration(sleepms) * time.Millisecond)
				clickPosition("next")
			}
		case 31:
			// cannot deal with this
		case 32:
			if len(answer.CurrentAnswer.Index) > 0 {
				for _, i := range answer.CurrentAnswer.Index {
					clickPosition("h" + strconv.Itoa(i))
					time.Sleep(time.Duration(sleepms) * time.Millisecond)
				}
				// Guessing
				clickPosition("h" + strconv.Itoa(rand.Intn(6)))
				time.Sleep(time.Duration(sleepms) * time.Millisecond)
				clickPosition("next")
			} else {
				// Guessing
				clickPosition("h" + strconv.Itoa(rand.Intn(6)))
				time.Sleep(time.Duration(sleepms) * time.Millisecond)
				clickPosition("next")
			}
		case 51:
			// fill in the blank
			robotgo.KeyTap("tab")
			robotgo.TypeStr(answer.CurrentAnswer.Detail.Word)
			clickPosition("submit")
			time.Sleep(time.Duration(sleepms) * time.Millisecond)
			clickPosition("next")
		default:
			// do nothing
		}
		time.Sleep(time.Duration(sleepms) * time.Millisecond)
	}
}

func clickPosition(mapIndex string) {
	pos := posMap[mapIndex]
	robotgo.Move(pos[0], pos[1])
	robotgo.Click("left", false)
}