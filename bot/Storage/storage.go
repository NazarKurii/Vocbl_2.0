package Storage

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"slices"
	"time"

	"github.com/NazarKurii/vocbl/Chat"
	"github.com/NazarKurii/vocbl/User"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func VerifyStorage() {

	storageData, storage := OpenStorage()

	todaysDate, _ := time.Parse("2006.01.02", time.Now().Format("2006.01.02"))
	formatedTodaysDate := todaysDate.Format("2006.01.02")
	for j, user := range storageData {
		for i, expretion := range user.Storage {

			repeatDate, _ := time.Parse("2006.01.02", expretion.ReapeatDate)

			if repeatDate.Before(todaysDate) {
				storageData[j].Storage[i].ReapeatDate = formatedTodaysDate
				storageData[j].DaylyTestTries = 2
			}
		}

	}

	WriteToStorage(storageData, storage)

}

func OpenStorage() ([]User.User, *os.File) {
	storage, _ := os.OpenFile("/home/nazar/nazzar/vocbl/vocblStorage/storage.json", os.O_RDWR, 0644)

	oldStorageData, _ := io.ReadAll(storage)
	var storageData []User.User
	_ = json.Unmarshal(oldStorageData, &storageData)

	return storageData, storage
}

func WriteToStorage(u []User.User, storage *os.File) {

	var users, _ = json.Marshal(u)

	_ = storage.Truncate(0)

	_, _ = storage.WriteAt(users, 0)
	storage.Sync()

}

func DefineUser(userId int64) (User.User, error) {
	var storage, _ = OpenStorage()
	userIndex := slices.IndexFunc(storage, func(u User.User) bool { return u.UserId == userId })
	if userIndex == -1 {
		return User.User{}, errors.New("User does not exist")
	}
	return storage[len(storage)-1], nil
}

func CreateUser(bot *tgbotapi.BotAPI, userId int64, updates tgbotapi.UpdatesChannel) User.User {
	var _, storage = OpenStorage()
	var newUser = User.User{
		UserId: userId,
		Chat: Chat.Chat{
			Bot:     bot,
			Updates: updates,
			ChatId:  userId,
		},
		DaylyTestTries: 2,
	}

	WriteToStorage([]User.User{newUser}, storage)

	return newUser
}
