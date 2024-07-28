package User

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/NazarKurii/Vocbl_2.0.git/Chat"
	"github.com/NazarKurii/Vocbl_2.0.git/Expretion"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type User struct {
	UserId  int64                 `json:"user_id"`
	Storage []Expretion.Expretion `json:"storage"`
	Chat    Chat.Chat             `json:"chat"`

	TestInfo struct {
		DaylyTestTries int    `json:"dayly_test_tries"`
		LastFailDate   string `json:"last_fail_date"`
		CanPassTest    bool   `json:"can_pass_test"`
	}
}

const (
	CreationDate = iota
	TestDate
)

var StartEroor = errors.New("Start error")

func (user User) StartMenue() {

	menue := []Chat.MessageComand{Chat.MessageComand{"Vocbl services", "services"}, Chat.MessageComand{"How do I use Vocbl?", "how"}}

	user.Chat.SendMessegeComand(menue, "Hello, I'm Vocbl!", 2)

}

func (user User) Services() {
	user.Chat.SendCommands([]Chat.MessageComand{Chat.MessageComand{"Add Card", "/add"}, Chat.MessageComand{"Remove Card", "/remove"}, Chat.MessageComand{"Get Quized", "/quiz"}, Chat.MessageComand{"Find Card", "/card"}, Chat.MessageComand{"Test", "/test"}, Chat.MessageComand{"Study", "/study"}}, "What can I do for you😁?", 3)
}

func (user User) FindExpretionsByDate(date string, dateType int) ([]Expretion.Expretion, []Expretion.Expretion) {
	var expretionsToRepeat []Expretion.Expretion
	var expretionsToStay []Expretion.Expretion

	for _, expretion := range user.Storage {
		switch dateType {
		case CreationDate:
			if expretion.CreationDate == date {
				expretionsToRepeat = append(expretionsToRepeat, expretion)
			}

		case TestDate:
			if expretion.ReapeatDate == date {
				expretionsToRepeat = append(expretionsToRepeat, expretion)
			} else {
				expretionsToStay = append(expretionsToStay, expretion)
			}

		}

	}
	return expretionsToRepeat, expretionsToStay
}

func (user User) RemoveFromUserStorage(expretion Expretion.Expretion) {

	user.Storage = slices.DeleteFunc(user.Storage, func(exp Expretion.Expretion) bool {
		return strings.ToLower(exp.Data) == strings.ToLower(expretion.Data)
	})
	user.SaveUsersData()

}

func (user User) AddToUserStorage(expretion Expretion.Expretion) {

	user.Storage = append(user.Storage, expretion)
	user.SaveUsersData()
}

func (user User) showListOfExpretions(expretions []Expretion.Expretion) {
	var message string
	for _, expretion := range expretions {
		message += fmt.Sprintf("•%v - %v\n\n", expretion.Data, expretion.TranslatedData)
	}
	user.Chat.SendMessege(message)

}

func (user User) SaveUsersData() {
	storage, _ := os.OpenFile("../vocblStorage/storage.json", os.O_RDWR, 0644)

	oldStorageData, _ := io.ReadAll(storage)
	var storageData []User
	_ = json.Unmarshal(oldStorageData, &storageData)

	var i = slices.IndexFunc(storageData, func(u User) bool {
		return u.UserId == user.UserId
	})

	storageData[i] = user

	var newStorageDataJson, _ = json.Marshal(storageData)

	_ = storage.Truncate(0)

	_, _ = storage.WriteAt(newStorageDataJson, 0)
	storage.Sync()
}

func (user User) FindExpretion(data string) (Expretion.Expretion, bool) {
	i := slices.IndexFunc(user.Storage, func(e Expretion.Expretion) bool {
		return e.Data == data
	})

	if i == -1 {
		return Expretion.Expretion{}, false
	} else {
		return user.Storage[i], true
	}
}

const (
	showAnswers = iota
	notSure
	correct
)

func (user User) Quiz(expretions []Expretion.Expretion, test bool) (int, error) {

	var totalAnswers = len(expretions)

	var fakeAnswers = user.getWrongAnswers(totalAnswers)

	var wrongAnswers int

	for i, experetion := range expretions {

		user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Answer", "true"}, Chat.MessageComand{"Not sure", "false"}}, experetion.Translations(), 1)
		status := user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
			switch update.CallbackQuery.Data {
			case "true":
				return showAnswers
			case "false":
				return notSure
			default:
				return -1
			}
		})
		if status == Chat.Start {
			return 0, StartEroor
		}

		if status == showAnswers {
			var fakeAnswer2, fakeUnswer1 string

			for true {
				var n1, n2 = rand.Intn(totalAnswers), rand.Intn(totalAnswers)
				if n1 == n2 {
					if n1 == 0 {
						n1++
					} else {
						n1--
					}
				}
				fakeUnswer1, fakeAnswer2 = fakeAnswers[n1], fakeAnswers[n2]
				if fakeAnswer2 != experetion.Data && fakeUnswer1 != experetion.Data {
					break
				}
			}

			answers := []Chat.MessageComand{Chat.MessageComand{experetion.Data, "true"}, Chat.MessageComand{fakeAnswer2, "false"}, Chat.MessageComand{fakeUnswer1, "false"}}

			rand.Shuffle(3, func(i, j int) {
				answers[i], answers[j] = answers[j], answers[i]
			})
			user.Chat.SendMessegeComand(answers, "Choose correct answer:", 3)
			status = user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
				switch update.CallbackQuery.Data {
				case "true":
					return correct
				case "false":

					return notSure
				default:
					return -1
				}
			})
		}
		if status == Chat.Start {
			return 0, StartEroor
		}

		if status == correct {
			user.Chat.SendMessege("Correct✅")
			continue
		} else {
			wrongAnswers++
			if !test {

				experetion.SendCard(user.Chat.Bot, user.Chat.ChatId)

			} else {
				user.Chat.SendMessege("Wrong❌")
			}
		}

		if i != totalAnswers-1 && !test {
			user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Continue", "continue"}}, "Ready to move on?", 1)
			status := user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
				switch update.CallbackQuery.Data {
				case "continue":
					return correct
				default:
					return -1
				}
			})
			if status == Chat.Start {

				return 0, StartEroor

			}
		}

	}

	return 100 - int(float64(wrongAnswers)/float64(totalAnswers)*100.0), nil

}

func (user User) getWrongAnswers(amount int) []string {

	var totalExpretionAmount = len(user.Storage)

	amount *= 2

	if amount > totalExpretionAmount {
		amount = totalExpretionAmount
	}
	var wrongAnswers = make([]string, amount)

	for i := 0; i < amount; i++ {
		wrongAnswers[i] = user.Storage[rand.Int63n(int64(amount))].Data

	}
	return wrongAnswers
}

func (user User) GetTestExpretions() []Expretion.Expretion {
	var testExpretions []Expretion.Expretion

	for _, expretion := range user.Storage {
		if expretion.ReapeatDate == time.Now().Format("2006.01.02") {
			testExpretions = append(testExpretions, expretion)
		}
	}
	return testExpretions
}

func (user User) PassedTestUpdateDates() {
	for i, expretion := range user.Storage {
		if expretion.ReapeatDate == time.Now().Format("2006.01.02") {
			user.Storage[i].Repeated++
			user.Storage[i].DefineRepeatDate()
		}
	}
	user.TestInfo.DaylyTestTries = 2
	user.SaveUsersData()
}

func (user User) FailedTestUpdateDates() {
	for i, expretion := range user.Storage {
		if expretion.ReapeatDate == time.Now().Format("2006.01.02") {
			user.Storage[i].Repeated = 1

		}
	}

	user.TestInfo.CanPassTest = false
	user.TestInfo.DaylyTestTries = 2
	user.TestInfo.LastFailDate = time.Now().Format("2006.01.02")
	user.SaveUsersData()
}

func (user User) StudyingProcces(expretions []Expretion.Expretion) (int, error) {
	for _, expretion := range expretions {
		expretion.SendCard(user.Chat.Bot, user.Chat.ChatId)
	}
	user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Get quized", "true"}}, "When you are finished remembering cards press button below🙃", 1)
	status := user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
		switch update.CallbackQuery.Data {
		case "true":
			return 0
		default:
			return -1
		}
	})
	if status == Chat.Start {
		return 0, StartEroor
	}
	return user.Quiz(expretions, false)
}

func (user User) GetNewCardsToStudy() []Expretion.Expretion {
	var expretionsToStudy []Expretion.Expretion

	todaysDate := time.Now().Format("2006.01.02")
	for _, expretion := range user.Storage {
		if expretion.CreationDate == todaysDate {
			expretionsToStudy = append(expretionsToStudy, expretion)
		}
	}
	return expretionsToStudy
}
