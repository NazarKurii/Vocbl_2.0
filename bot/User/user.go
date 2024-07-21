package User

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/NazarKurii/Vocbl_2.0.git/Chat"
	"github.com/NazarKurii/Vocbl_2.0.git/Expretion"
)

type User struct {
	UserId         int64                 `json:"user_id"`
	Storage        []Expretion.Expretion `json:"storage"`
	Chat           Chat.Chat             `json:"chat"`
	DaylyTestTries int                   `json:"dayly_test_tries"`
}

const (
	CreationDate = iota
	TestDate
)

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

	expretionToRemove := user.Storage[slices.IndexFunc(user.Storage, func(exp Expretion.Expretion) bool {
		return exp.Data == expretion.Data
	})]

	user.Storage = slices.DeleteFunc(user.Storage, func(exp Expretion.Expretion) bool {
		return exp.Data == expretionToRemove.Data
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
	user.Chat.SendCommands([]string{"Get quized", "Leave"}, "Ready to continue?", 1)
}

func (user User) SaveUsersData() {
	storage, _ := os.OpenFile("/home/nazar/nazzar/vocbl/vocblStorage/storage.json", os.O_RDWR, 0644)

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
