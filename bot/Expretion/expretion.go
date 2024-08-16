package Expretion

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Expretion struct {
	Data           string   `json:"data"`
	TranslatedData []string `json:"translated_data"`
	Examples       []string `json:"example"`
	Notes          string   `json:"notes"`
	EngTest        struct {
		ReapeatDate string `json:"repeat_date"`
		Repeated    int    `json:"repeated"`
	}

	UkrTest struct {
		ReapeatDate string `json:"repeat_date"`
		Repeated    int    `json:"repeated"`
	}
	CreationDate      string `json:"creation_date"`
	Pronunciation     string `json:"pronunciation"`
	PronunciationPath string `json:"pronunciation_path"`
}

func (e *Expretion) DefineEngRepeatDate() {

	date, _ := time.Parse("2006.01.02", e.EngTest.ReapeatDate)

	switch e.EngTest.Repeated {
	case 0:
		date = date.AddDate(0, 0, 1)
	case 1:
		date = date.AddDate(0, 0, 1)
	case 2:
		date = date.AddDate(0, 0, 1)
	case 3:
		date = date.AddDate(0, 0, 1)
	case 4:
		date = date.AddDate(0, 0, 4)
	case 5:
		date = date.AddDate(0, 0, 7)
	case 6:
		date = date.AddDate(0, 0, 16)
	case 7:
		date = date.AddDate(0, 0, 30)
	case 8:
		date = date.AddDate(0, 0, 60)
	case 9:
		date = date.AddDate(0, 0, 240)
	}

	e.EngTest.ReapeatDate = date.Format("2006.01.02")
}

func (e *Expretion) DefineUkrRepeatDate() {

	date, _ := time.Parse("2006.01.02", e.UkrTest.ReapeatDate)

	switch e.UkrTest.Repeated {
	case 0:
		date = date.AddDate(0, 0, 1)
	case 1:
		date = date.AddDate(0, 0, 1)
	case 2:
		date = date.AddDate(0, 0, 1)
	case 3:
		date = date.AddDate(0, 0, 1)
	case 4:
		date = date.AddDate(0, 0, 4)
	case 5:
		date = date.AddDate(0, 0, 7)
	case 6:
		date = date.AddDate(0, 0, 16)
	case 7:
		date = date.AddDate(0, 0, 30)
	case 8:
		date = date.AddDate(0, 0, 60)
	case 9:
		date = date.AddDate(0, 0, 240)
	}

	e.UkrTest.ReapeatDate = date.Format("2006.01.02")
}

type Card struct {
	voice    bool
	msg      tgbotapi.MessageConfig
	voiceMsg tgbotapi.VoiceConfig
}

func (e Expretion) SendCard(bot *tgbotapi.BotAPI, chatId int64) {

	card := fmt.Sprintf("• %v", strings.ToUpper(e.Data))
	if e.Pronunciation != "" {
		card += "\n: " + e.Pronunciation
	}

	card += "\n\nTranslations: "

	for _, translation := range e.TranslatedData {

		card += strings.ToLower(translation) + ", "

	}
	card = card[:len(card)-2]

	if len(e.Examples) != 0 {

		card += "\n\nExamples:"
		for _, exapmle := range e.Examples {
			card += fmt.Sprintf("\n-%v", exapmle)
		}
		card = card[:len(card)-1]
	}

	if e.Notes != "" {
		card += fmt.Sprintf("\n\nNotes: %v", e.Notes)
	}

	card += fmt.Sprintf("\n\nUkr test date: %v", e.UkrTest.ReapeatDate)
	card += fmt.Sprintf("\n\nEng test date: %v", e.EngTest.ReapeatDate)

	voice := tgbotapi.NewVoiceUpload(chatId, e.PronunciationPath)
	voice.Caption = card

	_, err := bot.Send(voice)
	if err != nil {
		msg := tgbotapi.NewMessage(chatId, card)
		bot.Send(msg)
	}

}

func (e Expretion) SendCardAddindg(bot *tgbotapi.BotAPI, chatId int64) {

	card := fmt.Sprintf("• %v", strings.ToUpper(e.Data))
	if e.Pronunciation != "" {
		card += "\n: " + e.Pronunciation
	}

	card += "\n\nTranslations: "

	for _, translation := range e.TranslatedData {

		card += strings.ToLower(translation) + ", "

	}
	card = card[:len(card)-2]

	if len(e.Examples) != 0 {

		card += "\n\nExamples:"
		for _, exapmle := range e.Examples {
			card += fmt.Sprintf("\n-%v", exapmle)
		}
		card = card[:len(card)-1]
	}

	if e.Notes != "" {
		card += fmt.Sprintf("\n\nNotes: %v", e.Notes)
	}

	card += fmt.Sprintf("\n\nUkr test date: %v", e.UkrTest.ReapeatDate)
	card += fmt.Sprintf("\n\nEng test date: %v", e.EngTest.ReapeatDate)

	voice := tgbotapi.NewVoiceUpload(chatId, e.PronunciationPath)
	voice.Caption = card

	voice.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("Save", "yes"), tgbotapi.NewInlineKeyboardButtonData("Leave", "no")}, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("Fill again", "again")})
	_, err := bot.Send(voice)
	if err != nil {
		msg := tgbotapi.NewMessage(chatId, card)
		bot.Send(msg)
	}

}

func (e Expretion) Translations() string {
	var translations string

	for _, translation := range e.TranslatedData {

		translations += translation + ", "
	}

	if translations == "" {

		return translations

	}
	return translations[:len(translations)-2]
}
