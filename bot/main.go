package main

import (
	"fmt"
	"os"

	"github.com/NazarKurii/Vocbl_2.0.git/Chat"
	AddindProces "github.com/NazarKurii/Vocbl_2.0.git/Components/AddingProcces"
	"github.com/NazarKurii/Vocbl_2.0.git/Expretion"
	"github.com/NazarKurii/Vocbl_2.0.git/ExpretionData"
	"github.com/NazarKurii/Vocbl_2.0.git/Storage"
	"github.com/NazarKurii/Vocbl_2.0.git/User"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	go Storage.VerifyStorage()

	var bot, err = tgbotapi.NewBotAPI("7421574054:AAH1pp0hDxoNQPPxFF1x5x6viuC6PX7UlJ4")
	if err != nil {
		//do
	}

	bot.Debug = true

	var u = tgbotapi.NewUpdate(0)
	u.Timeout = 60

	var updates, err2 = bot.GetUpdatesChan(u)
	if err2 != nil {
		//do
	}

	for update := range updates {

		if update.Message == nil {
			continue
		}
		user, err1 := Storage.DefineUser(update.Message.Chat.ID)

		if err1 != nil {
			user = Storage.CreateUser(bot, update.Message.Chat.ID, updates)
		} else {
			user.Chat = Chat.Chat{bot, updates, update.Message.Chat.ID}
		}

		switch update.Message.Text {
		case "/add":

			addExpretion(user)

		case "/card":
			card(user)
		case "/test":
			test(user)

		case "/study":
			//study(Chat.Chat{bot, updates, update.Message.Chat.ID})
		default:
			//Chat.Chat{bot, updates, update.Message.Chat.ID}.SendMessege("Unknown comand:(")
		}

	}

}

func test(user User.User) {
	voiceFile, err := os.Open("/home/nazar/nazzar/vocbl/audio/fuck")
	if err == nil {
		defer voiceFile.Close()

		voice := tgbotapi.NewVoiceUpload(user.UserId, tgbotapi.FileReader{
			Name:   "voice.ogg", // Assuming the file is an ogg file. Change if needed.
			Reader: voiceFile,
			Size:   -1,
		})

		voice.Caption = "jfnjrnrwwrjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjj"
		user.Chat.Bot.Send(voice)
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
		exp.SendCard(*user.Chat.Bot, user.Chat.ChatId)
	} else {
		user.Chat.SendMessege(fmt.Sprintf("There is no \"%v\" in your vocbl😢", userReprly))
	}
}

func addExpretion(user User.User) {

	user.Chat.SendMessege("Expretion to tranlate:")
	userReprly := user.Chat.GetUpdate()

	var newExpretion = Expretion.Expretion{Data: userReprly}
	if oldExpretion, exists := user.FindExpretion(userReprly); exists {
		oldExpretion.SendCard(*user.Chat.Bot, user.Chat.ChatId)
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
	fmt.Println(data, ".....................................")
	if err != nil {

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

		if newExpretion, err = AddindProces.AutoAddindProcces(user, data, newExpretion); err != nil {
			switch err {
			case AddindProces.AgainError:
				newExpretion, err = AddindProces.AutoAddindProcces(user, data, newExpretion)
				if err != nil {
					switch err {
					case AddindProces.AgainError:
						user.Chat.SendMessege("Let us start over. Send me \"/add\"")
					case AddindProces.RefuseError:
						user.Chat.SendMessege("Card wasn't added...")
						user.Chat.SendMessege("What can I do for you😁?")
					}
				} else {
					user.AddToUserStorage(newExpretion)
					user.Chat.SendMessege("Translation added")
					user.Chat.SendMessege("What can I do for you😁?")

				}
			case AddindProces.RefuseError:
				user.Chat.SendMessege("Card wasn't added...")
				user.Chat.SendMessege("What can I do for you😁?")
			}

		} else {

			user.AddToUserStorage(newExpretion)
			user.Chat.SendMessege("Translation added")
			user.Chat.SendMessege("What can I do for you😁?")

		}
	}
}
