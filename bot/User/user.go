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

	"Vocbl_2.0/Chat"
	"Vocbl_2.0/Expretion"
	"Vocbl_2.0/ExpretionData"

	Track "Vocbl_2.0/TrackConfig"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type User struct {
	UserId     int64     `json:"user_id"`
	Chat       Chat.Chat `json:"chat"`
	Tracks     map[string]Track.Track
	TracksKeys []string
	TestsLeft  int
}

const (
	CreationDate = iota
	UkrTestDate
	EngTestDate
)

var StartEroor = errors.New("Start error")

func (user User) StartMenue() {

	menue := []Chat.MessageComand{Chat.MessageComand{"Vocbl services", "services"}, Chat.MessageComand{"How do I use Vocbl?", "how"}}

	user.Chat.SendMessegeComand(menue, "Hello, I'm Vocbl!", 2)
}

func (user User) Services() {
	user.Chat.SendCommands([]Chat.MessageComand{Chat.MessageComand{"Add Card", "/add"}, Chat.MessageComand{"Remove Card", "/remove"}, Chat.MessageComand{"Get Quized", "/quiz"}, Chat.MessageComand{"Find Card", "/card"}, Chat.MessageComand{"Test", "/test"}, Chat.MessageComand{"Study", "/study"}}, "What can I do for you😁?", 3)
}

func (user User) FetchPronaunces() {
	for i, item := range user.Storage {

		data, err := ExpretionData.GetEpretionData(item.Data, 3)
		if err != nil {
			continue
		}
		user.Storage[i].PronunciationPath = data.Pronunciation.Path

	}
	user.SaveUsersData()
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

		case EngTestDate:
			if expretion.EngTest.ReapeatDate == date {
				expretionsToRepeat = append(expretionsToRepeat, expretion)
			} else {
				expretionsToStay = append(expretionsToStay, expretion)
			}
		case UkrTestDate:
			if expretion.UkrTest.ReapeatDate == date {
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
	defer storage.Close()
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

	data = strings.ToLower(data)

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

func (user User) SendCards(expretions []Expretion.Expretion) {
	for _, expretion := range expretions {
		expretion.SendCard(user.Chat.Bot, user.Chat.ChatId)
	}
}

func (user User) Quiz(expretions []Expretion.Expretion, test bool, testType string) (int, []Expretion.Expretion, error) {

	var totalAnswers = len(expretions)
	rand.Shuffle(totalAnswers, func(i, j int) {
		expretions[i], expretions[j] = expretions[j], expretions[i]
	})

	var fakeAnswers = user.getWrongAnswers(totalAnswers, testType)

	var wrongAnswers int
	var wrongCards []Expretion.Expretion

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
			return 0, nil, StartEroor
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

			answers := []Chat.MessageComand{Chat.MessageComand{experetion.Data, experetion.Data}, Chat.MessageComand{fakeAnswer2, fakeAnswer2}, Chat.MessageComand{fakeUnswer1, fakeAnswer2}, Chat.MessageComand{"Not sure", "false"}}

			rand.Shuffle(3, func(i, j int) {
				answers[i], answers[j] = answers[j], answers[i]
			})
			user.Chat.SendMessegeComand(answers, "Choose correct answer:", 4)
			status = user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
				switch update.CallbackQuery.Data {
				case experetion.Data:
					return correct
				case "false":
					return notSure
				case fakeAnswer2:
					return notSure
				default:
					return -1
				}
			})
		}
		if status == Chat.Start {
			return 0, nil, StartEroor
		}

		if status == correct {
			user.Chat.SendMessege("Correct✅")
			continue
		} else {
			wrongAnswers++
			if !test {
				wrongCards = append(wrongCards, experetion)

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

				return 0, nil, StartEroor

			}
		}

	}

	return 100 - int(float64(wrongAnswers)/float64(totalAnswers)*100.0), wrongCards, nil

}

func (user User) getWrongAnswers(amount int, testType string) []string {

	var totalExpretionAmount = len(user.Storage)

	amount *= 2

	if amount > totalExpretionAmount {
		amount = totalExpretionAmount
	}
	var wrongAnswers = make([]string, amount)

	if testType == "eng" {

		for i := 0; i < amount; i++ {

			wrongAnswers[i] = user.Storage[rand.Int63n(int64(amount))].Data

		}
	} else {
		for i := 0; i < amount; i++ {

			wrongAnswers[i] = strings.Join(user.Storage[rand.Int63n(int64(amount))].TranslatedData, ", ")

		}
	}
	return wrongAnswers
}

func (user User) GetEngTestExpretions() []Expretion.Expretion {
	var testExpretions []Expretion.Expretion

	for _, expretion := range user.Storage {
		if expretion.EngTest.ReapeatDate == time.Now().Format("2006.01.02") {
			testExpretions = append(testExpretions, expretion)
		}
	}
	return testExpretions
}

func (user User) GetUkrTestExpretions() []Expretion.Expretion {
	var testExpretions []Expretion.Expretion

	for _, expretion := range user.Storage {
		if expretion.UkrTest.ReapeatDate == time.Now().Format("2006.01.02") {

			expretion.Data, expretion.TranslatedData = strings.Join(expretion.TranslatedData, ", "), []string{expretion.Data}
			testExpretions = append(testExpretions, expretion)
		}
	}

	return testExpretions
}

func (user User) GreetingTestMessage() error {
	eng, ukr := user.TestInfo.EngTest.Status, user.TestInfo.UkrTest.Status
	var newErr = errors.New("Can't pass any of the tests")

	switch {
	case ukr == Failed && eng == Failed:
		user.Chat.SendMessege("You have failed both your tests, loser...")
		return newErr
	case ukr == Failed && eng == Passed:
		user.Chat.SendMessege("You have passed \"Eng test\",  thank God you've passed the \"Ukr test\" at least...")
		return newErr
	case ukr == Failed && eng == Prepared:
		user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Start", "eng"}, Chat.MessageComand{"Leave", "leave"}}, fmt.Sprintf("You have failed \"Ukr test\", hope you won't fuck the \"Eng\" one up... \n\n✅Tries left:%v\n\n Hit start when you're ready😒", user.TestInfo.EngTest.DaylyTestTries), 1)
		return nil
	case ukr == Passed && eng == Failed:
		user.Chat.SendMessege("You have passed \"Ukr test\",  thank God you've passed the \"Ukr test\" at least...")
		return newErr
	case ukr == Passed && eng == Passed:
		user.Chat.SendMessege("You have passed both tests🥰")
		return newErr
	case ukr == Passed && eng == Prepared:
		user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Start", "eng"}, Chat.MessageComand{"Leave", "leave"}}, fmt.Sprintf("You have passed \"Ukr test\", don't fuck the \"Eng\" one up... \n\n✅Tries left:%v\n\n Hit start when you're ready😒", user.TestInfo.EngTest.DaylyTestTries), 1)
		return nil
	case ukr == Prepared && eng == Failed:
		user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Start", "eng"}, Chat.MessageComand{"Leave", "leave"}}, fmt.Sprintf("You have failed \"Eng test\", hope you won't fuck the \"Ukr\" one up... \n\n✅Tries left:%v\n\n Hit start when you're ready😒", user.TestInfo.UkrTest.DaylyTestTries), 1)
		return nil
	case ukr == Prepared && eng == Passed:
		user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Start", "eng"}, Chat.MessageComand{"Leave", "leave"}}, fmt.Sprintf("You have passed \"Eng test\", don't fuck the \"Ukr\" one up... \n\n✅Tries left:%v\n\n Hit start when you're ready😒", user.TestInfo.UkrTest.DaylyTestTries), 1)
		return nil
	case ukr == Prepared && eng == Prepared:
		user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Eng", "eng"}, Chat.MessageComand{"Ukr", "ukr"}, Chat.MessageComand{"Leave", "leave"}}, fmt.Sprintf("Don`t fuck the tests up...\n\"Eng\" - choose English translation of Ukrainian word\n\"Ukr\" - choose Ukrainian translation of English word \n\n✅Eng tries left:%v\n✅Ukr tries left:%v", user.TestInfo.EngTest.DaylyTestTries, user.TestInfo.UkrTest.DaylyTestTries), 2)
		return nil
	default:
		user.Chat.SendMessege("Something went wrong🥶.\nTry again)")
		return newErr
	}
}

func (user User) PassedEngTestUpdateDates() {
	for i, expretion := range user.Storage {
		if expretion.EngTest.ReapeatDate == time.Now().Format("2006.01.02") {
			user.Storage[i].EngTest.Repeated++
			user.Storage[i].DefineEngRepeatDate()
		}
	}
	user.TestInfo.EngTest.Status = Passed
	user.TestInfo.EngTest.DaylyTestTries = 2
	user.SaveUsersData()
}

func (user User) PassedUkrTestUpdateDates() {
	for i, expretion := range user.Storage {
		if expretion.UkrTest.ReapeatDate == time.Now().Format("2006.01.02") {
			user.Storage[i].UkrTest.Repeated++
			user.Storage[i].DefineUkrRepeatDate()
		}
	}
	user.TestInfo.UkrTest.Status = Passed
	user.TestInfo.EngTest.DaylyTestTries = 2
	user.SaveUsersData()
}

func (user *User) UpdateTestCardsData(testType string, successRate int) {

	if successRate == 100 {
		if testType == "ukr" {
			go user.PassedUkrTestUpdateDates()
		} else {
			go user.PassedEngTestUpdateDates()
		}
	} else {
		var test *TestData
		if testType == "ukr" {
			test = &user.TestInfo.UkrTest
		} else {
			test = &user.TestInfo.EngTest
		}

		if test.DaylyTestTries == 0 {
			test.FailedTestUpdateData()

		}
	}
}

func (user User) StudyingProcces(expretions []Expretion.Expretion, testType string) (int, []Expretion.Expretion, error) {
	user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Get quized", "quize"}, Chat.MessageComand{"See cards", "cards"}}, "Want to see cards first or pass a quize?", 1)
	var seeCards bool
	user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
		switch update.CallbackQuery.Data {
		case "quize":
			return 1
		case "cards":
			seeCards = true
			return 1
		default:
			return -1
		}
	})
	if seeCards {

		for _, expretion := range expretions {
			expretion.SendCard(user.Chat.Bot, user.Chat.ChatId)
		}
		user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Get quized", "true"}, Chat.MessageComand{"Leave", "false"}}, "When you are finished remembering cards, you can be quized🙃", 1)
		status := user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
			switch update.CallbackQuery.Data {
			case "true":
				return 0
			case "false":
				return Chat.Start
			default:
				return -1
			}
		})
		if status == Chat.Start {
			return 0, nil, StartEroor
		}
	}
	if testType == "ukr" {

		for _, expretion := range expretions {
			expretion.Data, expretion.TranslatedData = strings.Join(expretion.TranslatedData, ", "), []string{expretion.Data}
		}
	}
	return user.Quiz(expretions, false, testType)
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
