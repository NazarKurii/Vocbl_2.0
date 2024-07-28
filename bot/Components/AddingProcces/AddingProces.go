package AddindProces

import (
	"errors"
	"fmt"
	"slices"
	"strings"
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

func AutoAddindProcces(user User.User, data ExpretionData.ExpretionData, newExpretion Expretion.Expretion) (Expretion.Expretion, error) {

	newExpretion.Data = strings.ToLower(newExpretion.Data)

	translatedData, err := ChooseTranslations(user, data.Translations, Choise{"Choose translations:", "Custom translation", "Translations", "Translation"})
	if err != nil {
		return Expretion.Expretion{}, err
	}
	newExpretion.TranslatedData = translatedData

	data.Translations = slices.DeleteFunc(data.Translations, func(e ExpretionData.Translation) bool {
		var result = true
		for _, translation := range translatedData {
			if translation == e.Translation {
				result = false
			}
		}
		return result
	})

	examples, err := ChooseExamples(user, data.Translations, Choise{"Choose exmples:", "Custom example", "Examples", "Example"})
	if err != nil {
		return Expretion.Expretion{}, err
	}
	newExpretion.Examples = examples

	user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Yes", "yes"}, Chat.MessageComand{"No", "no"}}, "Any notes", 1)
	status := user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {

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
	case Chat.Start:
		return Expretion.Expretion{}, User.StartEroor
	case Yes:
		user.Chat.SendMessege("Write me your notes:")
		newExpretion.Notes = user.Chat.GetUpdate()
		user.Chat.SendMessege("Notes added!")
	}

	todaysDate := time.Now().Format("2006.01.02")
	newExpretion.CreationDate = todaysDate
	newExpretion.ReapeatDate = todaysDate
	newExpretion.PronunciationPath = data.Pronunciation.Path
	newExpretion.Pronunciation = data.Pronunciation.Phonetic
	newExpretion.Repeated = 0

	newExpretion.SendCard(user.Chat.Bot, user.Chat.ChatId)
	user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Yes", "yes"}, Chat.MessageComand{"No", "no"}, Chat.MessageComand{"Fill again", "again"}}, "Add cart?", 2)
	status = user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {

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
	case Chat.Start:
		return Expretion.Expretion{}, User.StartEroor
	case Yes:
		return newExpretion, nil
	case No:
		return Expretion.Expretion{}, RefuseError
	case Again:
		return Expretion.Expretion{}, AgainError

	}
	return Expretion.Expretion{}, nil

}

func ChooseExamples(user User.User, translationsWithExamples []ExpretionData.Translation, choise Choise) ([]string, error) {
	var translationsCommands = make([]Chat.MessageComand, len(translationsWithExamples))
	for i, translationWithExamples := range translationsWithExamples {
		for _, example := range translationWithExamples.Examples {
			translationsCommands[i].Command = example
			translationsCommands[i].Callback = example
		}
	}

	var result, err = ChoseOptions(user, translationsCommands, choise, true)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ChooseTranslations(user User.User, translationsWithExamples []ExpretionData.Translation, choise Choise) ([]string, error) {

	var translationsCommands = make([]Chat.MessageComand, len(translationsWithExamples))
	for i, translationWithExample := range translationsWithExamples {
		translationsCommands[i].Command = translationWithExample.Translation
		translationsCommands[i].Callback = translationWithExample.Translation
	}

	var result, err = ChoseOptions(user, translationsCommands, choise, false)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type Choise struct {
	ChoseMessage  string
	CustomMessage string
	Choises       string
	Added         string
}

func ChoseOptions(user User.User, options []Chat.MessageComand, choise Choise, canBeEmpty bool) ([]string, error) {

	menues := Chat.NewMenues(options, choise.CustomMessage)
	messageId := user.Chat.SendMenue(menues[0], choise.ChoseMessage)
	var translations []string

	var zeroTranslations bool
	for i := 0; i >= 0; {
		if zeroTranslations {
			messageId = user.Chat.SendMenue(menues[i], choise.ChoseMessage)
			zeroTranslations = false
		} else {
			user.Chat.Bot.Send(tgbotapi.NewEditMessageReplyMarkup(user.Chat.ChatId, messageId, menues[i]))
		}
		switch translation, status := GetOptionsUpdate(user); status {
		case Save:
			if len(translations) == 0 && !canBeEmpty {
				user.Chat.SendMessege(fmt.Sprintf("%v field cannot be empty...", choise.Added))
				zeroTranslations = true
				continue
			}
			i = -1
		case Back:
			i--
			continue
		case Forward:
			i++
			continue
		case Added:
			menues[i].InlineKeyboard = slices.DeleteFunc(menues[i].InlineKeyboard, func(t []tgbotapi.InlineKeyboardButton) bool {
				return t[0].Text == translation
			})
			translations = append(translations, translation)

		case Custom:
			translation, status := CustomAdding(user, choise)
			zeroTranslations = true
			switch status {
			case AddMore:
				continue
			case Save:
				translations = append(translations, translation)

			}
		case Chat.Start:
			return nil, User.StartEroor

		}

	}

	return translations, nil

}

func CustomAdding(user User.User, choise Choise) (string, int) {
	user.Chat.SendMessege(choise.CustomMessage + ":")
	translation := user.Chat.GetUpdate()

	user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Save", "save"}, Chat.MessageComand{"Don't save", "false"}}, fmt.Sprintf("%v \"%v\" added", choise.Added, translation), 1)
	status := user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {

		switch update.CallbackQuery.Data {
		case "save":
			return Save
		case "false":
			return AddMore
		default:
			return Contine
		}
	})
	return translation, status
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

func GetOptionsUpdate(user User.User) (string, int) {
	var translation string
	status := user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {

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
