package main

import (
	"fmt"
	"math/rand"
	"slices"
	"strconv"
	"strings"

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

	var updates, _ = bot.GetUpdatesChan(u)

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

		case "/quiz":
			quizCommand(user)
		case "/study":
			//study(Chat.Chat{bot, updates, update.Message.Chat.ID})
		case "/remove":

		default:
			//Chat.Chat{bot, updates, update.Message.Chat.ID}.SendMessege("Unknown comand:(")
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
	if exp, exists := user.FindExpretion(strings.ToLower(userReprly)); exists {
		exp.SendCard(user.Chat.Bot, user.Chat.ChatId)
	} else {
		user.Chat.SendMessege(fmt.Sprintf("There is no \"%v\" in your vocbl😢", userReprly))
	}
}

func addExpretion(user User.User) {

	user.Chat.SendMessege("Expretion to tranlate:")
	userReprly := user.Chat.GetUpdate()

	var newExpretion = Expretion.Expretion{Data: userReprly}
	if oldExpretion, exists := user.FindExpretion(userReprly); exists {
		oldExpretion.SendCard(user.Chat.Bot, user.Chat.ChatId)
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

	go user.Chat.SendMessege("Wait, looking for data...")
	var data, _ = ExpretionData.GetEpretionData(userReprly, ExpretionData.RequestAttemts)

	if newExpretion, err := AddindProces.AutoAddindProcces(user, data, newExpretion); err != nil {
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

func quizCommand(user User.User) {
	user.Chat.SendMessege("Welcome to quize!")
	user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"10", "10"}, Chat.MessageComand{"20", "20"}, Chat.MessageComand{"50", "50"}, Chat.MessageComand{"Custom", "custom"}}, "Choose amount of expretions:", 2)
	amountOfExpretions := user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
		switch update.CallbackQuery.Data {
		case "10":
			return 10
		case "20":
			return 20
		case "50":
			return 50
		case "custom":
			return 0
		default:
			return -1
		}
	})

	if amountOfExpretions == 0 {
		amountOfExpretions = getExpretionAmount(user)

	}

	var totalStorageAmount = len(user.Storage)

	if totalStorageAmount < amountOfExpretions {
		user.Chat.SendMessege(fmt.Sprintf("Your vocbl consists less that \"%v\" expretions. The amount is set to \"%v\"", amountOfExpretions, totalStorageAmount))
		amountOfExpretions = totalStorageAmount
	}

	var expretionsToQuize = make([]Expretion.Expretion, amountOfExpretions)

	for i := 0; i < amountOfExpretions; i++ {
		var exist = true
		for exist {

			newExpretion := user.Storage[rand.Intn(totalStorageAmount)]

			exist = slices.ContainsFunc(expretionsToQuize, func(e Expretion.Expretion) bool {
				return e.Data == newExpretion.Data
			})
			if !exist {
				expretionsToQuize[i] = newExpretion
				break
			}
		}

	}

	var passTest = true

	for passTest {
		user.Chat.SendMessege("Good luck!")
		var successRate = int(user.Quiz(expretionsToQuize, false))

		var rateMessage string

		switch {
		case successRate == 100:
			rateMessage = "Your success rate is 💯%"
		case successRate < 100 && successRate > 70:
			rateMessage = fmt.Sprintf("Your success rate is \"%v%\"", successRate)
		case successRate < 70:
			rateMessage = fmt.Sprintf("Your success rate is \"%v%\" \nYou need to work harder...", successRate)
		}
		if successRate != 100 {
			user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Try again", "true"}, Chat.MessageComand{"Leave", "false"}}, rateMessage, 2)
			user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
				switch update.CallbackQuery.Data {
				case "true":

					return 1

				case "false":
					passTest = false
					return 1
				default:
					return -1
				}
			})
		} else {
			user.Chat.SendMessege(rateMessage)
			user.Chat.SendMessege("🎉That was cool!!!🎉")
			passTest = false
		}
	}
	user.Chat.SendMessege("I was happy to give a hand😁")

}

func getExpretionAmount(user User.User) int {
	var result int
	var isntNumber = true

	for isntNumber {
		isntNumber = false

		user.Chat.SendMessege("Provide expretions amount:")
		usersReply := user.Chat.GetUpdate()
		amountOfExpretions, err := strconv.Atoi(usersReply)
		if err != nil {
			user.Chat.SendMessege("Must be a number😁")
			isntNumber = true
			continue
		}
		user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Keep", "true"}, Chat.MessageComand{"Change", "false"}}, fmt.Sprintf("Amount set to \"%v\"", amountOfExpretions), 2)
		user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
			switch update.CallbackQuery.Data {
			case "true":
				result = amountOfExpretions
				return 1
			case "false":
				isntNumber = true
				return No
			default:
				return -1
			}
		})
	}

	return result
}
