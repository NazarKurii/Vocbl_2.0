package Storage

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"slices"
	"time"

	"Vocbl_2.0/Chat"
	"Vocbl_2.0/User"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func VerifyStorage() {

	storageData, storage := OpenStorage()
	defer storage.Close()

	todaysDate, _ := time.Parse("2006.01.02", time.Now().Format("2006.01.02"))
	formatedTodaysDate := todaysDate.Format("2006.01.02")
	for j, user := range storageData {
		if !user.TestInfo.EngTest.Passed && !user.TestInfo.EngTest.Failed {
			if date, _ := time.Parse("2006.01.02", user.TestInfo.EngTest.LastFailDate); date.Before(todaysDate) {
				storageData[j].TestInfo.EngTest.Passed, storageData[j].TestInfo.EngTest.Failed = false, false
				storageData[j].TestInfo.EngTest.DaylyTestTries = 2
			}

		}
		if !user.TestInfo.UkrTest.Passed && !user.TestInfo.UkrTest.Failed {
			if date, _ := time.Parse("2006.01.02", user.TestInfo.UkrTest.LastFailDate); date.Before(todaysDate) {
				storageData[j].TestInfo.UkrTest.Passed, storageData[j].TestInfo.UkrTest.Failed = false, false
				storageData[j].TestInfo.UkrTest.DaylyTestTries = 2
			}

		}
		for i, expretion := range user.Storage {

			repeatDate, _ := time.Parse("2006.01.02", expretion.EngTest.ReapeatDate)

			if repeatDate.Before(todaysDate) {
				storageData[j].Storage[i].EngTest.ReapeatDate = formatedTodaysDate

			}

			if expretion.EngTest.Repeated == 0 {
				creationDate, _ := time.Parse("2006.01.02", expretion.CreationDate)
				if creationDate.Before(todaysDate) {
					storageData[j].Storage[i].CreationDate = formatedTodaysDate
					storageData[j].Storage[i].EngTest.Repeated = 0

				}
			}

			repeatDate, _ = time.Parse("2006.01.02", expretion.UkrTest.ReapeatDate)

			if repeatDate.Before(todaysDate) {
				storageData[j].Storage[i].UkrTest.ReapeatDate = formatedTodaysDate

			}

			if expretion.UkrTest.Repeated == 0 {
				creationDate, _ := time.Parse("2006.01.02", expretion.CreationDate)
				if creationDate.Before(todaysDate) {
					storageData[j].Storage[i].CreationDate = formatedTodaysDate
					storageData[j].Storage[i].UkrTest.Repeated = 0

				}
			}

		}

	}

	WriteToStorage(storageData, storage)

}

func OpenStorage() ([]User.User, *os.File) {
	storage, _ := os.OpenFile("../vocblStorage/storage.json", os.O_RDWR, 0644)

	oldStorageData, _ := io.ReadAll(storage)
	var storageData []User.User
	_ = json.Unmarshal(oldStorageData, &storageData)

	return storageData, storage
}

func WriteToStorage(u []User.User, storage *os.File) {
	defer storage.Close()
	var users, _ = json.Marshal(u)

	_ = storage.Truncate(0)

	_, _ = storage.WriteAt(users, 0)
	storage.Sync()

}

func DefineUser(userId int64) (User.User, error) {
	var storage, st = OpenStorage()
	defer st.Close()
	userIndex := slices.IndexFunc(storage, func(u User.User) bool { return u.UserId == userId })

	if userIndex == -1 {
		return User.User{}, errors.New("User does not exist")
	}
	return storage[userIndex], nil
}

func CreateUser(bot *tgbotapi.BotAPI, userId int64, updates tgbotapi.UpdatesChannel) User.User {
	var users, storage = OpenStorage()

	var newUser = User.User{
		UserId: userId,
		Chat: Chat.Chat{
			Bot:     bot,
			Updates: updates,
			ChatId:  userId,
		},

		TestInfo: struct {
			EngTest   User.TestData
			UkrTest   User.TestData
			TestsLeft int
		}{
			EngTest:   User.TestData{2, "2006.01.01", false, false},
			UkrTest:   User.TestData{2, "2006.01.01", false, false},
			TestsLeft: 2,
		},
	}

	WriteToStorage(append(users, newUser), storage)

	return newUser
}
