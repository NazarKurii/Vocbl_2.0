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

		if update.Message != nil {
			user, err1 := Storage.DefineUser(update.Message.Chat.ID)

			if err1 != nil {

				user = Storage.CreateUser(bot, update.Message.Chat.ID, updates)
			} else {
				user.Chat = Chat.Chat{bot, updates, update.Message.Chat.ID}
			}

			var err = User.StartEroor
			message := update.Message.Commands()
			switch message {
			case "/start":
				user.StartMenue()

			case "/card":
				card(user)
			case "/test":
				daylyTest(user)
			case "/quiz":
				quizCommand(user)
			case "/study":
				study(user)
			case "/remove":
				remove(user)
			case "/edit":
				edit(user)
			case "/vocbl":
				user.Chat.SendMessege(fmt.Sprintf("Your vocbl cosists \"%v\" cards😁", len(user.Storage)))
			case "/vocbl_all":
				vocblAll(user)
			case "/pronaunce":
				pronaunce(user)
			default:
				addExpretion(user, message)
			}
			if err != nil {
				user.Services()
			}
		} else if update.CallbackQuery.Data != "" {

			user, err1 := Storage.DefineUser(update.CallbackQuery.Message.Chat.ID)

			if err1 != nil {

				user = Storage.CreateUser(bot, update.CallbackQuery.Message.Chat.ID, updates)
			} else {
				user.Chat = Chat.Chat{bot, updates, update.CallbackQuery.Message.Chat.ID}
			}
			callBacks(user, update.CallbackQuery.Data)
		}

	}
}

func pronaunce(user User.User) {
	user.FetchPronaunces()
}

func callBacks(user User.User, data string) {
	switch data {
	case "services":
		user.Services()
	case "how":
		user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Show services", "s"}}, `✅ Vocble Overview:
		Vocble was developed to help people extend their vocabulary. That means you have to be fair with Vocble and, more importantly, with yourself. If you want to make progress, Vocble strictly recommends following the technique provided below, which explains how to pass daily tests and quizzes.
		
		✅ Adding New Cards:
		Cards - expretions with their translations, examples, pronunciations, and additional notes.
		You can add cards to your Vocble using "/add" command and remove them using "/remove" command. You can choose translations/examples given to you by Vocble or create custom ones during adding-card procces.
		
		✅ Testing Schedule:
		Vocble automatically defines the date when you need to test already added expressions(cards), so you won't forget them. Test days count from the day of the first test: 1/2/3/7/14/30/60/120/360. However, if you fail, the test date will be changed to the next day and day counter will be set to "2".
		
		✅ Studying Cards:
		Vocble also helps you study new expressions. You can get a list of cards you need to study today, and then you will be quizzed by Vocble so you can learn the expressions better. Vocable allows you to study cards witch are used in todays test if you have already failed it, othervise you cannot get acces to todays test cards.
		
		✅ Types of Quizzes:
		There are two types of quizzes Vocble provides:
		1. Test Quiz - used for the daily test.
		2. Studying Quiz - used for the studying process.
		
		The difference between them is Vocble's reaction to mistakes. If you make a mistake in the studying quiz, you get the wrong answered Expression's Card. On the other hand, if you make a mistake in the test quiz, you are just sent to the next quiz card.
		
		✅ Quiz Interaction Process:
		After you enter the quiz, you are sent a quiz expression's translations and two options: "Answer" and "Not Sure." Your next move is thinking about whether you know the English translation of the expression or not. If you don't, you press "Not Sure," and that counts as a mistake. If you know the answer, you press the "Answer" button, and next, you are provided with three possible answers, in which only one is correct.
		
		This system was designed to avoid remembering the answer by seeing it among possible answers. In real conversation, you won't be given three possible words, in which only one fits for your next move; you always have to remember that one correct word.
		
		❗️ Vocble strictly recommends following quiz rules if you want to get a much better result.
		
		
		❗️If you need to leave any of the mentioned processes, you can write "/start" at any moment.`, 1)

	case "s":
		user.Services()
	default:
		user.Chat.SendMessege("Button is no longer awailable🥲")

	}

}

func vocblAll(user User.User) {
	for _, card := range user.Storage {
		card.SendCard(user.Chat.Bot, user.Chat.ChatId)
	}
}

func study(user User.User) error {
	user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"New cards", "new"}, Chat.MessageComand{"Test cards", "test"}}, "Welcome to studying mode!\nWhar are we studying today?", 1)

	status := user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
		switch update.CallbackQuery.Data {
		case "new":
			return new
		case "test":
			return test
		default:
			return -1
		}
	})

	if status == Chat.Start {
		user.Chat.SendMessege("Hope you're ready for tests, othervise I'm kicking your lazy ass😄")
		return User.StartEroor
	}

	for true {
		var successRate int
		switch status {
		case new:
			expretiondsToStudy := user.GetNewCardsToStudy()
			if len(expretiondsToStudy) == 0 {
				user.Chat.SendMessege("There are no new cards today...")
				return nil
			}

			rate, err := user.StudyingProcces(expretiondsToStudy)
			successRate = rate
			if err != nil {
				user.Chat.SendMessege("Hope you're ready for tests, othervise I'm kicking your lazy ass😄")
				return User.StartEroor
			}

		case test:
			if !user.TestInfo.CanPassTest {

				expretiondsToStudy := user.GetTestExpretions()
				if len(expretiondsToStudy) == 0 {
					user.Chat.SendMessege("There are no test cards today...")
					return nil
				}

				rate, err := user.StudyingProcces(expretiondsToStudy)
				successRate = rate
				if err != nil {
					user.Chat.SendMessege("Hope you're ready for tests, othervise I'm kicking your lazy ass😄")
					return User.StartEroor
				}
			} else {
				user.Chat.SendMessege("You need to fuck your todays test up, if you want to repeat the test cards...")
				return nil
			}
		}
		var rateMessage string

		switch {
		case successRate == 100:
			rateMessage = "Your success rate is 💯%"
		case successRate < 100 && successRate > 70:
			rateMessage = fmt.Sprintf("Your success rate  is \"%v%%\"", successRate)
		case successRate < 70:
			rateMessage = fmt.Sprintf("Your success rate is \"%v%%\"\nYou need to work harder...", successRate)
		}

		user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Try again", "true"}, Chat.MessageComand{"Leave", "false"}}, rateMessage, 2)
		status := user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
			switch update.CallbackQuery.Data {
			case "true":
				return True
			case "false":
				return False
			default:
				return -1
			}
		})

		switch status {
		case Chat.Start:
			user.Chat.SendMessege("Hope you're ready for tests, othervise I'm kicking your lazy ass😄")
			return User.StartEroor
		case False:
			user.Chat.SendMessege("Hope you're ready for tests, othervise I'm kicking your lazy ass😄")
			return User.StartEroor
		case True:
		}

	}
	return nil
}

const (
	Yes = iota
	No
	Stop
	new
	test
	Continue = -1
)

func edit(user User.User) error {
	user.Chat.SendMessege("Card to edit:")
	userReprly := user.Chat.GetUpdate()

	if userReprly == "/start" {
		user.Chat.SendMessege("Removing procces was interupted...")

		return User.StartEroor
	}

	userReprly = strings.ToLower(userReprly)

	if oldExpretion, exists := user.FindExpretion(userReprly); exists {

		for true {
			oldExpretion.SendCard(user.Chat.Bot, user.Chat.ChatId)

			user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Translations/Examples", "translations"}, Chat.MessageComand{"Notes", "notes"}}, "What would you like to edit?", 1)

			status := user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
				switch update.CallbackQuery.Data {
				case "translations":
					return translations
				case "notes":
					return notes
				default:
					return -1
				}
			})

			if status == Chat.Start {
				user.Chat.SendMessege("Edditing was canceled...")
				return User.StartEroor
			}

			switch status {
			case translations:
				user.Chat.SendMessege("Wait, looking for data")
				expretionsData, _ := ExpretionData.GetEpretionData(oldExpretion.Data, ExpretionData.RequestAttemts)
				newTranslations, err := AddindProces.ChooseTranslations(user, expretionsData.Translations, AddindProces.Choise{"Choose translations:", "Custom translation:", "Translations", "Translation"})
				if err != nil {
					user.Chat.SendMessege("Edditing was canceled...")
					return User.StartEroor
				}
				oldExpretion.TranslatedData = newTranslations
				user.Chat.SendMessege("Translations changed. Change examples according to new translations:")
				newExamples, err := AddindProces.ChooseExamples(user, expretionsData.Translations, AddindProces.Choise{"Choose examples:", "Custom example:", "Examples", "Example"})
				if err != nil {
					user.Chat.SendMessege("Edditing was canceled...")
					return User.StartEroor
				}
				oldExpretion.Examples = newExamples
				oldExpretion.SendCard(user.Chat.Bot, user.Chat.ChatId)

				user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Save", "save"}, Chat.MessageComand{"Cancel changes", "cancel"}, Chat.MessageComand{"Edit(edit new card)", "edit"}}, "Card was succesfully edited", 2)

				status := user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
					switch update.CallbackQuery.Data {
					case "save":
						return save
					case "cancel":
						return cancel
					case "edit":
						return Edit
					default:
						return -1
					}
				})

				if status == Chat.Start {
					user.Chat.SendMessege("Edditing was canceled...")
					return User.StartEroor
				}

				switch status {
				case save:

					user.RemoveFromUserStorage(Expretion.Expretion{Data: oldExpretion.Data})
					user.AddToUserStorage(oldExpretion)
					user.Chat.SendMessege("Card was eddited and saved😁")
					return nil
				case cancel:
					user.Chat.SendMessege("Card was NOT eddited...")
					return nil
				case Edit:
					continue
				}

			case notes:
				for true {

					user.Chat.SendMessege("Write me new notes:")
					notes := user.Chat.GetUpdate()
					user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Save", "save"}, Chat.MessageComand{"Edit", "edit"}, Chat.MessageComand{"Cancel changes", "cancel"}}, fmt.Sprintf("New notes: \n\n%v", notes), 2)

					status = user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
						switch update.CallbackQuery.Data {
						case "save":
							return save
						case "cancel":
							return cancel
						case "edit":
							return Edit
						default:
							return -1
						}
					})

					if status == Chat.Start {
						user.Chat.SendMessege("Card was NOT eddited...")
						return User.StartEroor
					}
					oldExpretion.Notes = notes

					switch status {
					case save:

						user.RemoveFromUserStorage(Expretion.Expretion{Data: oldExpretion.Data})
						user.AddToUserStorage(oldExpretion)
						user.Chat.SendMessege("Card was eddited and saved😁")
						return nil
					case cancel:
						user.Chat.SendMessege("Card was NOT eddited...")
						return nil
					case Edit:
						continue
					}
				}

			}

		}

	} else {
		user.Chat.SendMessege(fmt.Sprintf("There is no \"%v\" in your vocbl😢", userReprly))

	}
	return nil
}

const (
	translations = iota
	examples
	notes
	save
	change
	Edit
	cancel
)

func remove(user User.User) error {
	user.Chat.SendMessege("Card to remove:")
	userReprly := user.Chat.GetUpdate()

	if userReprly == "/start" {
		user.Chat.SendMessege("Removing procces was interupted...")

		return User.StartEroor
	}

	if oldExpretion, exists := user.FindExpretion(userReprly); exists {
		oldExpretion.SendCard(user.Chat.Bot, user.Chat.ChatId)
		user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Yes", "yes"}, Chat.MessageComand{"No", "no"}}, "Remove card?", 1)
		status := user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
			if update.Message != nil {
				return Continue
			}
			switch update.CallbackQuery.Data {
			case "no":
				return No
			case "yes":
				return Yes
			default:
				return Continue
			}
		})
		switch status {
		case Chat.Start:
			return User.StartEroor
		case No:
			user.Chat.SendMessege("Card wasnt removed...")

		case Yes:
			go user.RemoveFromUserStorage(oldExpretion)
			user.Chat.SendMessege("Card was removed...")

		}
	} else {
		user.Chat.SendMessege(fmt.Sprintf("There is no \"%v\" in your vocbl😢", userReprly))

	}
	return nil
}

func card(user User.User) error {
	user.Chat.SendMessege("Provide card:")
	userReprly := user.Chat.GetUpdate()
	if userReprly == "/start" {
		return User.StartEroor
	}
	if exp, exists := user.FindExpretion(strings.ToLower(userReprly)); exists {
		exp.SendCard(user.Chat.Bot, user.Chat.ChatId)
	} else {
		user.Chat.SendMessege(fmt.Sprintf("There is no \"%v\" in your vocbl😢", userReprly))
	}
	return nil
}

func addExpretion(user User.User, new string) error {

	userReprly := new

	go user.Chat.SendStartCommand()

	if userReprly == "/start" {
		user.Chat.SendMessege("Adding procces was interupted...")

		return User.StartEroor
	}

	var newExpretion = Expretion.Expretion{Data: userReprly}
	if oldExpretion, exists := user.FindExpretion(userReprly); exists {
		oldExpretion.SendCard(user.Chat.Bot, user.Chat.ChatId)
		user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Yes", "yes"}, Chat.MessageComand{"Leave both", "no"}, Chat.MessageComand{"Leave(stop adding)", "leave"}}, "Card already exists in your vocbl. \n\nWant to replase it?", 1)
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
			user.Chat.SendMessege("The old card removed.\nContinue adding new card")
		case No:
			user.Chat.SendMessege("Ok, we're leaving both😁")
		case Stop:
			user.Chat.SendMessege("Adding procces was stopped...")
			return User.StartEroor
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
				user.Chat.SendMessege("Card was added")
				user.Chat.SendMessege("What can I do for you😁?")

			}
		case AddindProces.RefuseError:
			user.Chat.SendMessege("Card wasn't added...")
			user.Chat.SendMessege("What can I do for you😁?")
		case User.StartEroor:
			user.Chat.SendMessege("Adding procces was stopped...")
			return User.StartEroor
		}

	} else {

		user.AddToUserStorage(newExpretion)
		user.Chat.SendMessege("Translation added")

	}
	return nil
}

func getQuizExpretionsAmount(user User.User, totalStorageAmount int) (int, error) {
	var amountOfExpretions int

	for true {
		user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"10", "10"}, Chat.MessageComand{"20", "20"}, Chat.MessageComand{"50", "50"}, Chat.MessageComand{"Custom", "custom"}}, "Choose amount of expretions:", 2)
		status := user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
			switch update.CallbackQuery.Data {
			case "10":
				amountOfExpretions = 10
				return 1
			case "20":
				amountOfExpretions = 20
				return 1
			case "50":
				amountOfExpretions = 50
				return 1
			case "custom":
				return 0
			default:
				return -1
			}
		})

		if status == Chat.Start {
			return 0, User.StartEroor
		}

		if amountOfExpretions == 0 {
			amountOfExpretions = getExpretionAmount(user)
		}

		if totalStorageAmount < amountOfExpretions {
			user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Save", "true"}, Chat.MessageComand{"Change", "false"}}, fmt.Sprintf("Your vocbl consists less that \"%v\" cards. The amount is set to \"%v\"", amountOfExpretions, totalStorageAmount), 1)
			status := user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
				switch update.CallbackQuery.Data {
				case "true":
					amountOfExpretions = totalStorageAmount
					return True
				case "false":
					return False
				default:
					return -1
				}
			})
			if status == Chat.Start {
				return 0, User.StartEroor
			}

			if status == True {
				break
			}

		} else {
			break
		}
	}
	return amountOfExpretions, nil
}

const (
	True = iota
	False
)

func quizCommand(user User.User) error {
	user.Chat.SendMessege("Welcome to quize!")
	totalStorageAmount := len(user.Storage)
	if totalStorageAmount < 5 {
		user.Chat.SendMessege("Sorry, but to get quized you have to add al leat 5 cards to your vocbl😁...")
		return User.StartEroor
	}

	amountOfExpretions, err := getQuizExpretionsAmount(user, totalStorageAmount)
	if err != nil {
		user.Chat.SendMessege("Quiz procces is over...")
		return User.StartEroor
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
		var successRate, err = user.Quiz(expretionsToQuize, false)
		if err != nil {
			user.Chat.SendMessege("Quiz procces is over...")
			return User.StartEroor
		}

		var rateMessage string

		switch {
		case successRate == 100:
			rateMessage = "Your success rate is 💯%"
		case successRate < 100 && successRate > 70:
			rateMessage = fmt.Sprintf("Your success rate  is \"%v%%\"", successRate)
		case successRate < 70:
			rateMessage = fmt.Sprintf("Your success rate is \"%v%%\"\nYou need to work harder...", successRate)
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
	return nil
}

func getExpretionAmount(user User.User) int {
	var result int
	var isntNumber = true

	for isntNumber {
		isntNumber = false

		user.Chat.SendMessege("Provide cards amount:")
		usersReply := user.Chat.GetUpdate()
		amountOfExpretions, err := strconv.Atoi(usersReply)
		if err != nil {
			user.Chat.SendMessege("Must be a number😁")
			isntNumber = true
			continue
		}
		if amountOfExpretions <= 0 {
			user.Chat.SendMessege("Must be greater than \"0\"😁")
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

func daylyTest(user User.User) error {
	var testExpretions = user.GetTestExpretions()

	var amountOfTestExpretions = len(testExpretions)
	var maxMistakes = int(float64(amountOfTestExpretions) * 0.2)
	if amountOfTestExpretions == 0 && user.TestInfo.CanPassTest {
		user.Chat.SendMessege("There are no cards to test today😢...")
		return User.StartEroor
	} else if !user.TestInfo.CanPassTest {
		user.Chat.SendMessege("You've failed, study and try tomorrow🖕...")
		return User.StartEroor
	}

	for user.TestInfo.DaylyTestTries > 0 {
		var testMessage string
		if user.TestInfo.DaylyTestTries == 2 {
			testMessage = fmt.Sprintf("Total quetions: %v\nMax mistakes to have second try: %v\nReady?", amountOfTestExpretions, maxMistakes)
		} else {
			testMessage = "Last chance, ready?"
		}
		user.Chat.SendMessegeComand([]Chat.MessageComand{Chat.MessageComand{"Start", "start"}, Chat.MessageComand{"Leave", "leave"}}, testMessage, 1)

		status := user.Chat.GetUpdateFunc(func(update tgbotapi.Update) int {
			switch update.CallbackQuery.Data {
			case "start":
				return True
			case "leave":
				return False
			default:
				return -1
			}
		})

		if status == False || status == Chat.Start {

			return User.StartEroor
		}

		successRate, err := user.Quiz(testExpretions, true)
		user.TestInfo.DaylyTestTries--
		go user.SaveUsersData()
		if err != nil {
			if user.TestInfo.DaylyTestTries == 1 {
				user.Chat.SendMessege("Test was interupted, but you can still try once...")
				return User.StartEroor
			} else {
				user.Chat.SendMessege("Test was interupted and that was your last try! Study and try tomorrow...")
				go user.FailedTestUpdateDates()

			}

			return User.StartEroor

		}
		if successRate == 100 {
			if user.TestInfo.DaylyTestTries == 1 {

				user.Chat.SendMessege("Good job! Next test is yours also😁")
			} else {
				user.Chat.SendMessege("If I have to give you another chance next time, I'll kick your lazy ass... You've passed. Go study!")
			}
			go user.PassedTestUpdateDates()
			return nil
		} else if successRate >= 80 && user.TestInfo.DaylyTestTries > 0 {
			user.Chat.SendMessege("You failed, but you haven't made anough mistakes so I could kick your ass out of here😒...")
		} else {
			break
		}

	}
	go user.FailedTestUpdateDates()
	user.Chat.SendMessege("You've failed, study and try tomorrow🖕...")
	return nil

}
