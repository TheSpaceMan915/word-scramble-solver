package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Word list
var words = []string{"golang", "programming", "concurrent", "goroutine", "channel", "developer"}

// Scrambles a word
func scrambleWord(word string) string {
	runes := []rune(word)
	rand.Shuffle(len(runes), func(i, j int) { runes[i], runes[j] = runes[j], runes[i] })
	return string(runes)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Create the Fyne app
	myApp := app.New()
	myWindow := myApp.NewWindow("Word Scramble Solver")
	myWindow.Resize(fyne.NewSize(400, 300))

	// Select a random word
	originalWord := words[rand.Intn(len(words))]
	scrambledWord := scrambleWord(originalWord)

	// UI Elements
	wordLabel := widget.NewLabel("Unscramble this word: " + scrambledWord)
	timerLabel := widget.NewLabel("⏳ Time Left: 15s")
	hintLabel := widget.NewLabel("")
	inputEntry := widget.NewEntry()
	inputEntry.SetPlaceHolder("Type your guess here...")

	// Result label
	resultLabel := widget.NewLabel("")

	// Layout
	layout := container.NewVBox(
		wordLabel,
		timerLabel,
		hintLabel,
		inputEntry,
		resultLabel,
	)

	// Timer, Hints, and Game Logic
	timeLeft := 15
	ticker := time.NewTicker(1 * time.Second)
	hintTimer := time.NewTimer(5 * time.Second)

	// Goroutine for handling game logic
	go func() {
		for {
			select {
			case <-ticker.C:
				timeLeft--
				timerLabel.SetText(fmt.Sprintf("⏳ Time Left: %ds", timeLeft))
				if timeLeft == 0 {
					ticker.Stop()
					dialog.ShowInformation("Time's Up!", "The correct word was: "+originalWord, myWindow)
					return
				}

			case <-hintTimer.C:
				hintLabel.SetText("Hint: Starts with '" + string(originalWord[0]) + "'")
				time.AfterFunc(5*time.Second, func() {
					hintLabel.SetText("Hint: Ends with '" + string(originalWord[len(originalWord)-1]) + "'")
				})
			}
		}
	}()

	// User input handling
	inputEntry.OnChanged = func(text string) {
		if strings.EqualFold(text, originalWord) {
			ticker.Stop()
			dialog.ShowInformation("Congratulations!", "🎉 You solved it!", myWindow)
		}
	}

	// Display the window
	myWindow.SetContent(layout)
	myWindow.ShowAndRun()
}
