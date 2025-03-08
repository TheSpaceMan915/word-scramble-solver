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
	resultLabel := widget.NewLabel("")

	// Layout
	layout := container.NewVBox(
		wordLabel,
		timerLabel,
		hintLabel,
		inputEntry,
		resultLabel,
	)

	// Channels for communication
	timerChannel := make(chan int)
	hintChannel := make(chan string)
	gameOverChannel := make(chan bool)

	// Timer Goroutine
	go func() {
		timeLeft := 15
		for timeLeft > 0 {
			select {
			case <-gameOverChannel:
				close(timerChannel)
				return
			default:
				time.Sleep(1 * time.Second)
				timeLeft--
				timerChannel <- timeLeft
			}
		}
		// Notify game over if time runs out
		gameOverChannel <- true
	}()

	// Hint Goroutine
	go func() {
		time.Sleep(5 * time.Second)
		hintChannel <- "Hint: Starts with '" + string(originalWord[0]) + "'"
		time.Sleep(5 * time.Second)
		hintChannel <- "Hint: Ends with '" + string(originalWord[len(originalWord)-1]) + "'"
	}()

	// UI Update Goroutine
	go func() {
		for {
			select {
			case timeLeft := <-timerChannel:
				timerLabel.SetText(fmt.Sprintf("⏳ Time Left: %ds", timeLeft))
			case hint := <-hintChannel:
				hintLabel.SetText(hint)
			case <-gameOverChannel:
				dialog.ShowInformation("Time's Up!", "The correct word was: "+originalWord, myWindow)
				return
			}
		}
	}()

	// User Input Handling
	inputEntry.OnChanged = func(text string) {
		if strings.EqualFold(text, originalWord) {
			gameOverChannel <- true
			dialog.ShowInformation("Congratulations!", "🎉 You solved it!", myWindow)
		}
	}

	// Display the window
	myWindow.SetContent(layout)
	myWindow.ShowAndRun()
}
