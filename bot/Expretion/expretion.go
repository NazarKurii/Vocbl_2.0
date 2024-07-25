package Expretion

import (
	"fmt"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Expretion struct {
	Data              string   `json:"data"`
	TranslatedData    []string `json:"translated_data"`
	Examples          []string `json:"example"`
	Notes             string   `json:"notes"`
	ReapeatDate       string   `json:"repeat_date"`
	Repeated          int      `json:"repeated"`
	CreationDate      string   `json:"creation_date"`
	Pronunciation     string   `json:"pronunciation"`
	PronunciationPath string   `json:"pronunciation_path"`
}

func (e Expretion) DefineRepeatDate() string {

	date, _ := time.Parse("2006.01.02", e.ReapeatDate)

	switch e.Repeated {
	case 1:
		date = date.AddDate(0, 0, 1)
	case 2:
		date = date.AddDate(0, 0, 2)
	case 3:
		date = date.AddDate(0, 0, 3)
	case 4:
		date = date.AddDate(0, 0, 7)
	case 5:
		date = date.AddDate(0, 0, 14)
	case 6:
		date = date.AddDate(0, 0, 30)
	case 7:
		date = date.AddDate(0, 0, 60)
	case 8:
		date = date.AddDate(0, 0, 120)
	case 9:
		date = date.AddDate(0, 0, 360)
	}

	return date.Format("2006.01.02")
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

	card += "\n\nExamples:"
	for _, exapmle := range e.Examples {
		card += fmt.Sprintf("\n-%v", exapmle)
	}
	card = card[:len(card)-1]

	if e.Notes != "" {
		card += fmt.Sprintf("\n\nNotes: %v", e.Notes)
	}

	card += fmt.Sprintf("\n\nTest date: %v", e.ReapeatDate)

	voiceFile, err := os.Open(e.PronunciationPath)
	fmt.Println(err, "<<<<,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,")
	if err == nil {

		defer voiceFile.Close()

		voice := tgbotapi.NewVoiceUpload(chatId, tgbotapi.FileReader{
			Name:   "voice.ogg",
			Reader: voiceFile,
			Size:   -1,
		})

		voice.Caption = card

		bot.Send(voice)
	} else {

		bot.Send(tgbotapi.NewMessage(chatId, card))
		fmt.Println("Error opening voice file:", err)

	}
}

func (e Expretion) Translations() string {
	var translations string

	for _, translation := range e.TranslatedData {

		translations += translation + ", "
	}

	return translations[:len(translations)-2]
}
