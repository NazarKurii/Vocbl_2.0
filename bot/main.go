package main

import (
	"fmt"

	"github.com/NazarKurii/vocbl/Chat"
	AddindProces "github.com/NazarKurii/vocbl/Components/AddingProcces"
	"github.com/NazarKurii/vocbl/Expretion"
	"github.com/NazarKurii/vocbl/ExpretionData"
	"github.com/NazarKurii/vocbl/Storage"
	"github.com/NazarKurii/vocbl/User"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	//go Storage.VerifyStorage()

	var bot, err = tgbotapi.NewBotAPI("7421574054:AAH1pp0hDxoNQPPxFF1x5x6viuC6PX7UlJ4")
	if err != nil {
		fmt.Println(err)
	}

	bot.Debug = false

	var u = tgbotapi.NewUpdate(0)
	u.Timeout = 60

	var updates, err2 = bot.GetUpdatesChan(u)
	if err2 != nil {
		fmt.Println(err2)
	}

	for update := range updates {

		if update.Message == nil {
			continue
		}

		var user, err = Storage.DefineUser(update.Message.Chat.ID)

		if err != nil {
			user = Storage.CreateUser(bot, update.Message.Chat.ID, updates)
		}

		user.Chat.Bot = bot
		user.Chat.Updates = updates

		switch update.Message.Text {

		case "/add":
			fmt.Println(",....")
			addExpretion(bot, update.Message.Chat.ID, updates)

		case "/remove":

		case "/card":

		}

	}
}

const (
	Yes = iota
	No
	Stop
	Continue = -1
)

func card(user User.User) {
	user.Chat.SendMessege("Provide expretion:")
	userReprly := user.Chat.GetUpdate()
	if exp, exists := user.FindExpretion(userReprly); exists {
		exp.Card(user.Chat.ChatId).Send(*user.Chat.Bot)
	} else {
		user.Chat.SendMessege(fmt.Sprintf("There is no \"%v\" in your vocbl😢", userReprly))
	}
}

func addExpretion(bot *tgbotapi.BotAPI, chatId int64, updates tgbotapi.UpdatesChannel) {
	fmt.Println("....")
	user, err1 := Storage.DefineUser(chatId)
	if err1 != nil {
		user = Storage.CreateUser(bot, chatId, updates)
	}

	user.Chat.SendMessege("Expretion to tranlate:")
	userReprly := user.Chat.GetUpdate()
	var newExpretion = Expretion.Expretion{Data: userReprly}
	if oldExpretion, exists := user.FindExpretion(userReprly); exists {
		oldExpretion.Card(user.Chat.ChatId).Send(*user.Chat.Bot)
		user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Yes", "yes"}, Chat.MessageComand{"Leave both", "no"}, Chat.MessageComand{"Leave(stop adding)", "leave"}}, "Expretion already exists in your vocbl. \n\nWant to replase it?", 1)
		status := user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
			if update.Message != nil {
				return Continue
			}

			switch update.CallbackQuery.Data {
			case "no":
				return No
			case "yes":
				return Yes
			case "leave":

				return Stop
			default:
				return Continue
			}
		})

		switch status {
		case Yes:
			user.RemoveFromUserStorage(newExpretion)
			user.Chat.SendMessege("The old expretion removed.\nContinue adding new Expretion😁")
		case No:
			user.Chat.SendMessege("Ok, we're leaving both😁")
		case Stop:
			user.Chat.SendMessege("What can I do for you😁?")
			return
		}
	}

	var data, err = ExpretionData.GetEpretionData(userReprly, ExpretionData.RequestAttemts)
	if err != nil {
		fmt.Println(err, "........................................")
		if newExpretion, err = AddindProces.CustomAddindProcces(user, data, newExpretion); err != nil {
			switch err {
			case AddindProces.AgainError:
				newExpretion, err = AddindProces.CustomAddindProcces(user, data, newExpretion)
				if err != nil {
					switch err {
					case AddindProces.AgainError:
						user.Chat.SendMessege("Let us start over. Send me \"/add\"")
					case AddindProces.RefuseError:
						user.Chat.SendMessege("Card wasn't added...")
						user.Chat.SendMessege("What can I do for you😁?")
					}
				}
			case AddindProces.RefuseError:
				user.Chat.SendMessege("Card wasn't added...")
				user.Chat.SendMessege("What can I do for you😁?")
			}

		}
	} else {
		if newExpretion, err = AddindProces.CustomAddindProcces(user, data, newExpretion); err != nil {
			switch err {
			case AddindProces.AgainError:
				newExpretion, err = AddindProces.CustomAddindProcces(user, data, newExpretion)
				if err != nil {
					switch err {
					case AddindProces.AgainError:
						user.Chat.SendMessege("Let us start over. Send me \"/add\"")
					case AddindProces.RefuseError:
						user.Chat.SendMessege("Card wasn't added...")
						user.Chat.SendMessege("What can I do for you😁?")
					}
				}
			case AddindProces.RefuseError:
				user.Chat.SendMessege("Card wasn't added...")
				user.Chat.SendMessege("What can I do for you😁?")
			}

		}
	}

	user.AddToUserStorage(newExpretion)
}
