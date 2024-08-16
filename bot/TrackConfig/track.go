package Track

import (
	"errors"
	"fmt"
	"time"

	"Vocbl_2.0/Expretion"
)

const (
	Passed = iota
	Failed
	Prepared
)

type Track struct {
	Storage      []Expretion.Expretion `json:"storage"`
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
	WarningMessage(int, int) string
	FailedTestUpdateData()
}

func (t *FromLanguage) InteruptionError() string {
	if t.DaylyTestTries == 1 {
		return "Test was interupted, but you can still try once..."
	} else {

		return "Test was interupted and that was your last try! Study and try tomorrow..."

	}
}

func (t *ToLanguage) InteruptionError() string {
	if t.DaylyTestTries == 1 {
		return "Test was interupted, but you can still try once..."
	} else {

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

func (t FromLanguage) WarningMessage(maxMistakes, amountOfTestExpretions int) string {
	if t.DaylyTestTries == 2 {
		return fmt.Sprintf("Total quetions: %v\nMax mistakes to have second try: %v\nReady?", amountOfTestExpretions, maxMistakes)
	} else {
		return "If you make at least one mistake, you are out🥶"
	}
}

func (t ToLanguage) WarningMessage(maxMistakes, amountOfTestExpretions int) string {
	if t.DaylyTestTries == 2 {
		return fmt.Sprintf("Total quetions: %v\nMax mistakes to have second try: %v\nReady?", amountOfTestExpretions, maxMistakes)
	} else {
		return "If you make at least one mistake, you are out🥶"
	}
}

const (
	CreationDate = iota
	UkrTestDate
	EngTestDate
)

var StartEroor = errors.New("Start error")

func (t *FromLanguage) FailedTestUpdateData() {

	t.Status = Failed
	t.DaylyTestTries = 2
	t.LastFailDate = time.Now().Format("2006.01.02")

}

func (t *ToLanguage) FailedTestUpdateData() {

	t.Status = Failed
	t.DaylyTestTries = 2
	t.LastFailDate = time.Now().Format("2006.01.02")

}
