package Expretion

import (
	"fmt"
	"os"
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

func (c Card) Send(bot tgbotapi.BotAPI) {
	if c.voice {
		bot.Send(c.voiceMsg)
	} else {
		bot.Send(c.msg)
	}
}

func (e Expretion) Card(chatId int64) Card {

	card := fmt.Sprintf("• %v \n\nTranslation: %v", e.Data, e.TranslatedData)

	if len(e.Examples) > 0 {
		card += "\n\nExamples:"
		for _, exapmle := range e.Examples {
			card += fmt.Sprintf("\n-%v", exapmle)
		}
	}
	if e.Notes != "" {
		card += fmt.Sprintf("\n\nNotes: %v", e.Notes)
	}

	card += fmt.Sprintf("\n\nCreation date: %v\n\nTest date: %v", e.CreationDate, e.ReapeatDate)
	if e.Pronunciation != "" {
		card += "\n\nPronunciation: " + e.Pronunciation
	}

	voiceFile, err := os.Open("e.PronunciationPath")

	if err == nil {

		voice := tgbotapi.NewVoiceUpload(chatId, tgbotapi.FileReader{
			Name:   "",
			Reader: voiceFile,
			Size:   -1,
		})

		voice.Caption = card
		defer voiceFile.Close()
		return Card{
			voiceMsg: voice, voice: true,
		}
	} else {
		defer voiceFile.Close()
		return Card{msg: tgbotapi.NewMessage(chatId, card), voice: false}
	}
}
