package Storage

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"slices"
	"time"

	"github.com/NazarKurii/Vocbl_2.0.git/Chat"
	"github.com/NazarKurii/Vocbl_2.0.git/User"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func VerifyStorage() {

	storageData, storage := OpenStorage()
	defer storage.Close()

	todaysDate, _ := time.Parse("2006.01.02", time.Now().Format("2006.01.02"))
	formatedTodaysDate := todaysDate.Format("2006.01.02")
	for j, user := range storageData {
		if !user.TestInfo.CanPassTest {
			if date, _ := time.Parse("2006.01.02", user.TestInfo.LastFailDate); date.Before(todaysDate) {
				storageData[j].TestInfo.CanPassTest = true
				storageData[j].TestInfo.DaylyTestTries = 2
			}

		}
		for i, expretion := range user.Storage {

			repeatDate, _ := time.Parse("2006.01.02", expretion.ReapeatDate)

			if repeatDate.Before(todaysDate) {
				storageData[j].Storage[i].ReapeatDate = formatedTodaysDate

			}

			if expretion.Repeated == 0 {
				creationDate, _ := time.Parse("2006.01.02", expretion.CreationDate)
				if creationDate.Before(todaysDate) {
					storageData[j].Storage[i].CreationDate = formatedTodaysDate
					storageData[j].Storage[i].Repeated = 0

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
	return storage[len(storage)-1], nil
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
			DaylyTestTries int    `json:"dayly_test_tries"`
			LastFailDate   string `json:"last_fail_date"`
			CanPassTest    bool   `json:"can_pass_test"`
		}{
			2,
			"",
			true,
		},
	}

	WriteToStorage(append(users, newUser), storage)

	return newUser
}
