package answer

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
)

var WordList []WordInfo

func FindAnswer(topicID int, vocabTaskInfo VocabTaskStruct, rawJSON string) (ans Answer) {
	// Do stuff
	switch topicID {
	// Choose translation from word	
	case 15:
		stemTrimed := strings.ReplaceAll(vocabTaskInfo.Stem.Content, " ", "")
	Loop15:
		for i := 0; i < len(WordList); i++ {
			if WordList[i].Word == stemTrimed {
				// Match the translation
				for j := 0; j < len(WordList[i].Content); j++ {
					for k := 0; k < len(vocabTaskInfo.Options); k++ {
						if compareTranslation(WordList[i].Content[j].Meaning, vocabTaskInfo.Options[k].Content) {
							ans.Found = true
							ans.Detail.Translation = WordList[i].Content[j].Meaning
							ans.Index = append(ans.Index, k)
							break Loop15
						}
					}
				}
			}
		}

	// Choose translation based on some word
	case 11:
		stemConverted := strings.ReplaceAll(vocabTaskInfo.Stem.Content, "  ", " ")
	Loop11:
		for i := 0; i < len(WordList); i++ {
			for j := 0; j < len(WordList[i].Content); j++ {
				for k := 0; k < len(WordList[i].Content[j].ExampleEnglish); k++ {
					if WordList[i].Content[j].ExampleEnglish[k] == stemConverted {
						ans.Detail.Translation = WordList[i].Content[j].Meaning
						ans.Found = true
						break Loop11
					}
				}
			}
		}
		for i := 0; i < len(vocabTaskInfo.Options); i++ {
			regex := regexp.MustCompile(`（.*?）`)
			vocabTaskInfo.Options[i].Content = string(regex.ReplaceAll([]byte(vocabTaskInfo.Options[i].Content), []byte("")))
			if compareTranslation(ans.Detail.Translation, vocabTaskInfo.Options[i].Content) {
				ans.Index = append(ans.Index, i)
				break
			}
		}
		// Debug
		log.Println(vocabTaskInfo.Options)


	// Choose translation based on some voice
	case 22:
	Loop22:
		for i := 0; i < len(WordList); i++ {
			if WordList[i].Word == vocabTaskInfo.Stem.Content {
				for j := 0; j < len(WordList[i].Content); j++ {
					for k := 0; k < len(vocabTaskInfo.Options); k++ {
						regex := regexp.MustCompile(`（.*?）`)
						vocabTaskInfo.Options[k].Content = string(regex.ReplaceAll([]byte(vocabTaskInfo.Options[k].Content), []byte("")))
						if compareTranslation(vocabTaskInfo.Options[k].Content, WordList[i].Content[j].Meaning) {
							ans.Index = append(ans.Index, k)
							ans.Found = true
							break Loop22
						}
					}
				}
			}
		}

	// Choose some word pair
	case 31:
		for i := 0; i < len(vocabTaskInfo.Stem.Remark); i++ {
			for j := 0; j < len(vocabTaskInfo.Options); j++ {
				if strings.Contains(vocabTaskInfo.Stem.Remark[i].SenMarked, vocabTaskInfo.Options[j].Content) {
					ans.Index = append(ans.Index, j)
					ans.Found = true
				}
			}
		}

	// Organize some word pair
	case 32:
		regexFind := regexp.MustCompile(`"remark":".*?"`)
		raw := regexFind.FindString(rawJSON)
		ans.Detail.Raw = raw[10 : len(raw)-1]
	Loop32:
		for i := 0; i < len(WordList); i++ {
			for j := 0; j < len(WordList[i].Content); j++ {
				for k := 0; k < len(WordList[i].Content[j].Usage); k++ {
					if strings.Contains(WordList[i].Content[j].Usage[k], ans.Detail.Raw) {
						ans.Detail.Word = WordList[i].Content[j].Usage[k]
						ans.Found = true
						break Loop32
					}
				}
			}
		}

		// Get the words
		var word string
		var words []string
		lowerWord := strings.ToLower(ans.Detail.Word)

		for i := 0; i < len(lowerWord); i++ {
			if (lowerWord[i] <= 'z' && lowerWord[i] >= 'a') || lowerWord[i] == '\'' {
				word += string(lowerWord[i])
			}
			if lowerWord[i] == ' ' {
				words = append(words, string(word))
				word = ""
			}
		}

		// Find correct order
		for i := 0; i < len(words); i++ {
			for j := 0; j < len(vocabTaskInfo.Options); j++ {
				if vocabTaskInfo.Options[j].Content == words[i] {
					ans.Index = append(ans.Index, j)
					break
				}
			}
		}

	// Some fill in blank
	case 51:
		//Find from remark
		regexFind := regexp.MustCompile(`"remark":".*?"`)
		raw := regexFind.FindString(rawJSON)
		ans.Detail.Raw = raw[10 : len(raw)-1]
	Loop51:
		for i := 0; i < len(WordList); i++ {
			for j := 0; j < len(WordList[i].Content); j++ {
				for k := 0; k < len(WordList[i].Content[j].Usage); k++ {
					if strings.Contains(WordList[i].Content[j].Usage[k], ans.Detail.Raw) {
						ans.Detail.Word = WordList[i].Word
						ans.Found = true
						break Loop51
					}
				}
			}
		}
		//If remark isn't found, then check the word length and wtip.
		if !ans.Found {
			for i := 0; i < len(WordList); i++ {
				if vocabTaskInfo.WLen == len(WordList[i].Word) && vocabTaskInfo.WTip == WordList[i].Word[:len(vocabTaskInfo.WTip)] {
					ans.Detail.Uncertain = true
					ans.Detail.Word = WordList[i].Word
				}
			}
		}
	}

	log.Println("[I] Found ans: " + fmt.Sprintln(ans))
	return
}

// The translation may vary.
func compareTranslation(str1 string, str2 string) bool {
	//If they are actually the same, then that's good
	if str1 == str2 {
		return true
	}

	//Compare length first
	if len(str1) != len(str2) {
		return false
	}

	//Split the classification of current word and translation
	str1Slice := strings.Split(str1, " ")
	str2Slice := strings.Split(str2, " ")

	//Judge the classification of current word
	if str1Slice[0] != str2Slice[0] {
		return false
	}

	//Split
	str1 = strings.ReplaceAll(str1Slice[1], "；", "，")
	str2 = strings.ReplaceAll(str2Slice[1], "；", "，")
	str1Slice = strings.Split(str1, "，")
	str2Slice = strings.Split(str2, "，")

	//Compare split length
	if len(str1Slice) != len(str2Slice) {
		return false
	}

	//Sort and compare content
	sort.Strings(str1Slice)
	sort.Strings(str2Slice)

	for i := 0; i < len(str1Slice); i++ {
		if str1Slice[i] != str2Slice[i] {
			return false
		}
	}

	return true
}
