package Chat

import (
	"fmt"
	"math"
	"slices"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Chat struct {
	Bot     *tgbotapi.BotAPI        `json:"-"`
	Updates tgbotapi.UpdatesChannel `json:"-"`
	ChatId  int64                   `json:"chat_id"`
}

type MessageComand struct {
	Command string `json:"command"`

	Callback string `json:"call_back"`
}

func (chat Chat) SendCommands(commandMessages []string, message string, rowAmount int) {
	var commandsRows = make([][]tgbotapi.KeyboardButton, rowAmount)

	var step = int(math.Ceil(float64(len(commandMessages)) / float64(rowAmount)))

	for i, commandMessage := range commandMessages {
		commandsRows[i/step] = append(commandsRows[i/step], tgbotapi.NewKeyboardButton(commandMessage))
	}

	msg := tgbotapi.NewMessage(chat.ChatId, message)
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		commandsRows...)
	_, _ = chat.Bot.Send(msg)

}

func (chat Chat) SendMessege(message string) {
	msg := tgbotapi.NewMessage(chat.ChatId, message)
	_, err := chat.Bot.Send(msg)
	if err != nil {
		panic(err)
	}
}

func (chat Chat) SendMessegeComand(commandMessages []MessageComand, message string, rowAmount int) {
	var commandsRows = make([][]tgbotapi.InlineKeyboardButton, rowAmount)

	var step = int(math.Ceil(float64(len(commandMessages)) / float64(rowAmount)))

	for i, commandMessage := range commandMessages {
		commandsRows[i/step] = append(commandsRows[i/step], tgbotapi.NewInlineKeyboardButtonData(commandMessage.Command, commandMessage.Callback))
	}

	msg := tgbotapi.NewMessage(chat.ChatId, message)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		commandsRows...)
	_, _ = chat.Bot.Send(msg)

}

func (chat Chat) GetUpdate() string {
	for update := range chat.Updates {
		if update.Message == nil {
			continue
		}
		if text := update.Message.Text; text != "" {
			return text
		}
	}
	return ""
}

func (chat Chat) GetUpdateFunc(f func(update tgbotapi.Update) int) int {
	var result int
	for update := range chat.Updates {
		if result = f(update); result != -1 {

			break
		}
	}
	return result

}

var bot, _ = tgbotapi.NewBotAPI("7421574054:AAH1pp0hDxoNQPPxFF1x5x6viuC6PX7UlJ4")

func NewMenues(options []MessageComand, customInfo string) []tgbotapi.InlineKeyboardMarkup {

	options = slices.DeleteFunc(options, func(option MessageComand) bool {
		return len(option.Callback) > 64
	})

	length := int(math.Ceil(float64(len(options)) / 5))
	if length == 0 {
		length++
	}
	var optionBunches = make([][][]tgbotapi.InlineKeyboardButton, length)

	for i, option := range options {
		optionBunches[i/5] = append(optionBunches[i/5], tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(option.Command, option.Callback)))
	}

	var max = len(optionBunches) - 1

	fmt.Println(max, "MAXXXXXXXXX>>>>>>>>>>>>>>>>>?????????????")
	var menues = make([]tgbotapi.InlineKeyboardMarkup, 0, len(optionBunches))
	for i, options := range optionBunches {

		var commandRaws = make([][]tgbotapi.InlineKeyboardButton, 2)

		if i == 0 {
			commandRaws[1] = append(commandRaws[1], tgbotapi.NewInlineKeyboardButtonData("✖️", "_"))
		} else {
			commandRaws[1] = append(commandRaws[1], tgbotapi.NewInlineKeyboardButtonData("⬅️", "back"))
		}

		commandRaws[1] = append(commandRaws[1], tgbotapi.NewInlineKeyboardButtonData("Save", "save"))

		if i == max {
			commandRaws[1] = append(commandRaws[1], tgbotapi.NewInlineKeyboardButtonData("✖️", "_"))
		} else {
			commandRaws[1] = append(commandRaws[1], tgbotapi.NewInlineKeyboardButtonData("➡️", "forward"))
		}

		commandRaws[0] = append(commandRaws[0], tgbotapi.NewInlineKeyboardButtonData(customInfo, "custom"))
		var menu [][]tgbotapi.InlineKeyboardButton = commandRaws
		if len(options) != 0 {
			menu = append(options, menu...)
		}
		menues = append(menues, tgbotapi.NewInlineKeyboardMarkup(menu...))
	}

	return menues
}
func (chat Chat) SendMenue(replyMurkUp tgbotapi.InlineKeyboardMarkup, message string) int {
	fmt.Println("Test3........................", chat.Bot == nil)
	msg := tgbotapi.NewMessage(chat.ChatId, message)
	msg.ReplyMarkup = replyMurkUp
	sentMessage, err := chat.Bot.Send(msg)
	if err != nil {
		panic(err)
	}
	return sentMessage.MessageID
}

func (chat Chat) SendMessegeMenue(menue tgbotapi.InlineKeyboardMarkup, message string) {

	msg := tgbotapi.NewMessage(chat.ChatId, message)
	msg.ReplyMarkup = menue

}
