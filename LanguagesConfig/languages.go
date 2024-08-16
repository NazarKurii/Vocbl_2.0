package Track

import (
	"errors"
	"fmt"
	"time"
)

const (
	Passed = iota
	Failed
	Prepared
)

type Track struct {
	FromLanguage FromLanguage
	ToLanguage   ToLanguage
}

type FromLanguage struct {
	DaylyTestTries int    `json:"dayly_test_tries"`
	LastFailDate   string `json:"last_fail_date"`
	Status         int
}

type ToLanguage struct {
	DaylyTestTries int    `json:"dayly_test_tries"`
	LastFailDate   string `json:"last_fail_date"`
	Status         int
}

type Test interface {
	InteruptionError() string
	AfterTestMessage(int) string
}

func (t *FromLanguage) InteruptionError() string {
	if t.DaylyTestTries == 1 {
		return "Test was interupted, but you can still try once..."
	} else {
		go t.FailedTestUpdateDates()
		return "Test was interupted and that was your last try! Study and try tomorrow..."

	}
}

func (t *ToLanguage) InteruptionError() string {
	if t.DaylyTestTries == 1 {
		return "Test was interupted, but you can still try once..."
	} else {
		go t.FailedTestUpdateDates()
		return "Test was interupted and that was your last try! Study and try tomorrow..."

	}

}

func (t *FromLanguage) AfterTestMessage(successRate int) string {

	var message string
	if successRate == 100 {
		if t.DaylyTestTries == 1 {
			message = "Good job! Next test is yours also😁"
		} else {
			message = ("If I have to give you another chance next time, I'll kick your lazy ass... You've passed. Go study!")
		}
	} else if successRate >= 80 && t.DaylyTestTries > 0 {
		message = "You failed, but you haven't made anough mistakes so I could kick your ass out of here😒..."
	} else {
		message = "You have fucked this test up. Try tomorow 🤬..."
	}

	return message
}

func (t *ToLanguage) AfterTestMessage(successRate int) string {

	var message string
	if successRate == 100 {
		if t.DaylyTestTries == 1 {
			message = "Good job! Next test is yours also😁"
		} else {
			message = ("If I have to give you another chance next time, I'll kick your lazy ass... You've passed. Go study!")
		}
	} else if successRate >= 80 && t.DaylyTestTries > 0 {
		message = "You failed, but you haven't made anough mistakes so I could kick your ass out of here😒..."
	} else {
		message = "You have fucked this test up. Try tomorow 🤬..."
	}

	return message
}

func (t FromLanguage) MistakesMessage(maxMistakes, amountOfTestExpretions int) string {
	if t.DaylyTestTries == 2 {
		return fmt.Sprintf("Total quetions: %v\nMax mistakes to have second try: %v\nReady?", amountOfTestExpretions, maxMistakes)
	} else {
		return "If you make at least one mistake, you are out🥶"
	}
}

func (t ToLanguage) MistakesMessage(maxMistakes, amountOfTestExpretions int) string {
	if t.DaylyTestTries == 2 {
		return fmt.Sprintf("Total quetions: %v\nMax mistakes to have second try: %v\nReady?", amountOfTestExpretions, maxMistakes)
	} else {
		return "If you make at least one mistake, you are out🥶"
	}
}

func (t *FromLanguage) InteraptingTestMessage(user User.User) {

	if t.DaylyTestTries == 1 {
		user.Chat.SendMessege("Test was interupted, but you can still try once...")

	} else {
		user.Chat.SendMessege("Test was interupted and that was your last try! Study and try tomorrow...")
		t.FailedTestUpdateDates()
	}

}

func (t *ToLanguage) InteraptingTestMessage(user User.User) {

	if t.DaylyTestTries == 1 {
		user.Chat.SendMessege("Test was interupted, but you can still try once...")

	} else {
		user.Chat.SendMessege("Test was interupted and that was your last try! Study and try tomorrow...")
		t.FailedTestUpdateDates()
	}

}

const (
	CreationDate = iota
	UkrTestDate
	EngTestDate
)

var StartEroor = errors.New("Start error")

func (t *TestData) FailedTestUpdateData() {

	t.Status = Failed
	t.DaylyTestTries = 2
	t.LastFailDate = time.Now().Format("2006.01.02")

}
