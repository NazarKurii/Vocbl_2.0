package AddindProces

import (
	"errors"
	"fmt"
	"time"

	"github.com/NazarKurii/Vocbl_2.0.git/Chat"
	"github.com/NazarKurii/Vocbl_2.0.git/Expretion"
	"github.com/NazarKurii/Vocbl_2.0.git/ExpretionData"
	"github.com/NazarKurii/Vocbl_2.0.git/User"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	AgainError  = errors.New("User wants to fill out data again")
	RefuseError = errors.New("User does not want to add card")
)

func CustomAddindProcces(user User.User, data ExpretionData.ExpretionData, newExpretion Expretion.Expretion) (Expretion.Expretion, error) {
	return Expretion.Expretion{}, nil
}

func AutoAddindProcces(user User.User, data ExpretionData.ExpretionData, newExpretion Expretion.Expretion) (Expretion.Expretion, error) {

	newExpretion.TranslatedData = chooseTranslations(user, data.Translations, Choise{"Choose translations:", "Custom translation:", "Translations", "Translation"})
	newExpretion.Examples = chooseExamples(user, data.Translations, Choise{"Choose exmples:", "Custom example:", "Examples", "Example"})

	user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Yes", "yes"}, Chat.MessageComand{"No", "no"}}, "Any notes", 1)
	status := user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
		if update.Message != nil {
			return Contine
		}
		switch update.CallbackQuery.Data {
		case "yes":
			return Yes
		case "no":
			return No
		default:
			return Contine
		}
	})

	switch status {
	case Yes:
		user.Chat.SendMessege("Write me your notes:")
		newExpretion.Notes = user.Chat.GetUpdate()
		user.Chat.SendMessege("Notes added!")
	}

	todaysDate := time.Now().Format("2006s.01.02")
	newExpretion.CreationDate = todaysDate
	newExpretion.ReapeatDate = todaysDate
	newExpretion.PronunciationPath = data.Pronunciation.Path
	newExpretion.Pronunciation = data.Pronunciation.Phonetic

	newExpretion.SendCard(*user.Chat.Bot, user.Chat.ChatId)
	user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Yes", "yes"}, Chat.MessageComand{"No", "no"}, Chat.MessageComand{"Fill again", "again"}}, "Add cart?", 2)
	status = user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
		if update.Message != nil {
			return Contine
		}
		switch update.CallbackQuery.Data {
		case "yes":
			return Yes
		case "no":
			return No
		case "again":
			return Again
		default:
			return Contine
		}
	})
	switch status {
	case Yes:
		return newExpretion, nil
	case No:
		return Expretion.Expretion{}, RefuseError
	case Again:
		return Expretion.Expretion{}, AgainError

	}
	return Expretion.Expretion{}, nil

}

func chooseExamples(user User.User, translationsWithExamples []ExpretionData.Translation, choise Choise) []string {
	var translationsCommands = make([]Chat.MessageComand, len(translationsWithExamples))
	for i, translationWithExamples := range translationsWithExamples {
		for _, example := range translationWithExamples.Examples {
			translationsCommands[i].Command = example
			translationsCommands[i].Callback = example
		}
	}

	menues := Chat.NewMenues(translationsCommands, choise.customMessage)
	messageId := user.Chat.SendMenue(menues[0], choise.choseMessage)

	return getTransltions(user, menues, 0, messageId, false, false, []string{}, choise)
}

type Choise struct {
	choseMessage  string
	customMessage string
	choises       string
	added         string
}

func chooseTranslations(user User.User, translationsWithExamples []ExpretionData.Translation, choise Choise) []string {

	var translationsCommands = make([]Chat.MessageComand, len(translationsWithExamples))
	for i, translationWithExample := range translationsWithExamples {
		translationsCommands[i].Command = translationWithExample.Translation
		translationsCommands[i].Callback = translationWithExample.Translation
	}

	menues := Chat.NewMenues(translationsCommands, choise.customMessage)

	messageId := user.Chat.SendMenue(menues[0], choise.choseMessage)

	return getTransltions(user, menues, 0, messageId, false, false, []string{}, choise)

}

func getTransltions(user User.User, menues []tgbotapi.InlineKeyboardMarkup, menueIndex int, messageId int, edit, addMore bool, translations []string, choise Choise) []string {

	if edit {
		user.Chat.Bot.Send(tgbotapi.NewEditMessageReplyMarkup(user.Chat.ChatId, messageId, menues[menueIndex]))
	} else if addMore {
		messageId = user.Chat.SendMenue(menues[menueIndex], choise.choseMessage)

	}

	switch translation, status := getTransltionsUpdates(user); status {
	case Save:
		return append(translations, translation)
	case Back:
		return getTransltions(user, menues, menueIndex-1, messageId, true, false, translations, choise)
	case Forward:
		return getTransltions(user, menues, menueIndex+1, messageId, true, false, translations, choise)
	case Added:
		return getTransltions(user, menues, menueIndex, messageId, false, false, append(translations, translation), choise)
	default:
		return getTransltions(user, menues, menueIndex, messageId, false, false, translations, choise)
	case Custom:
		user.Chat.SendMessege(choise.customMessage)
		translation = user.Chat.GetUpdate()
		translations = append(translations, translation)
		user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Save", "save"}, Chat.MessageComand{choise.choises, "add"}}, fmt.Sprintf("%v \"%v\" added", choise.added, translation), 1)
		status = user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
			if update.Message != nil {
				return Contine
			}
			switch update.CallbackQuery.Data {
			case "save":
				return Save
			case "add":
				return AddMore
			default:
				return Contine
			}
		})
		switch status {
		case AddMore:
			return getTransltions(user, menues, menueIndex, messageId, false, true, translations, choise)
		case Save:
			return translations
		}
	}
	return []string{}

}

const (
	Contine = -1
	Custom  = iota
	Save
	Back
	Forward
	Added
	AddMore
	Yes
	No
	Again
)

func getTransltionsUpdates(user User.User) (string, int) {
	var translation string
	status := user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
		if update.Message != nil {
			return Contine
		}
		switch data := update.CallbackQuery.Data; data {
		case "save":
			return Save
		case "custom":
			return Custom
		case "back":
			return Back
		case "forward":
			return Forward
		default:
			translation = data
			return Added
		}
	})

	return translation, status
}
