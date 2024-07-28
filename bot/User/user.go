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

func (user User) StartMenue() string {
	var action string

	menue := []Chat.MessageComand{Chat.MessageComand{"Vocbl services", "services"}, Chat.MessageComand{"How do I use Vocbl?", "how"}}

	user.Chat.SendMessegeComand(menue, "Hello, I'm Vocbl!", 2)
	status := user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
		switch update.CallbackQuery.Data {
		case "services":
			return 1
		case "how":
			return 2
		default:
			return -1
		}
	})

	switch status {
	case Chat.Start:
		return user.StartMenue()
	case 1:
		user.Services()
		return ""
	case 2:
		user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Show services", "s"}}, `✅ Vocble Overview:
		Vocble was developed to help people extend their vocabulary. That means you have to be fair with Vocble and, more importantly, with yourself. If you want to make progress, Vocble strictly recommends following the technique provided below, which explains how to pass daily tests and quizzes.
		
		✅ Adding New Cards:
		Cards - expretions with their translations, examples, pronunciations, and additional notes.
		You can add cards to your Vocble using "/add" command and remove them using "/remove" command. You can choose translations/examples given to you by Vocble or create custom ones during adding-card procces.
		
		✅ Testing Schedule:
		Vocble automatically defines the date when you need to test already added expressions(cards), so you won't forget them. Test days count from the day of the first test: 1/2/3/7/14/30/60/120/360. However, if you fail, the test date will be changed to the next day and day counter will be set to "2".
		
		✅ Studying Cards:
		Vocble also helps you study new expressions. You can get a list of cards you need to study today, and then you will be quizzed by Vocble so you can learn the expressions better. Vocable allows you to study cards witch are used in todays test if you have already failed it, othervise you cannot get acces to todays test cards.
		
		✅ Types of Quizzes:
		There are two types of quizzes Vocble provides:
		1. Test Quiz - used for the daily test.
		2. Studying Quiz - used for the studying process.
		
		The difference between them is Vocble's reaction to mistakes. If you make a mistake in the studying quiz, you get the wrong answered Expression's Card. On the other hand, if you make a mistake in the test quiz, you are just sent to the next quiz card.
		
		✅ Quiz Interaction Process:
		After you enter the quiz, you are sent a quiz expression's translations and two options: "Answer" and "Not Sure." Your next move is thinking about whether you know the English translation of the expression or not. If you don't, you press "Not Sure," and that counts as a mistake. If you know the answer, you press the "Answer" button, and next, you are provided with three possible answers, in which only one is correct.
		
		This system was designed to avoid remembering the answer by seeing it among possible answers. In real conversation, you won't be given three possible words, in which only one fits for your next move; you always have to remember that one correct word.
		
		❗️ Vocble strictly recommends following quiz rules if you want to get a much better result.
		
		
		❗️If you need to leave any of the mentioned processes, you can write "/start" at any moment.`, 1)

		status = user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
			if update.CallbackQuery.Data == "s" {
				return 1
			} else {
				return -1
			}
		})

		if status == 1 {
			user.Services()
			return ""
		}

	}
	return action
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
