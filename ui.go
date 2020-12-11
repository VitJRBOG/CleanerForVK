package main

import "fmt"

func MakeUI() {
	var ui UI
	ui.showMainMenu()
	for {
		ui.setUserSelection()
		success := ui.showSelected()
		if success {
			break
		} else {
			ui.showMessageOfWrongInput()
		}
	}
	ui.showMessageOfExit()
}

type UI struct {
	UserSelection string
}

func (ui *UI) showMainMenu() {
	fmt.Print("\n[Main menu]\n" +
		"1. Cleaning of wallposts;\n" +
		"2. Cleaning comments under wallposts.\n" +
		"--- Enter number of menu item and press «Enter» ---\n" +
		"> ")
}

func (ui *UI) setUserSelection() {
	_, err := fmt.Scan(&ui.UserSelection)
	if err != nil {
		panic(err.Error())
	}
}

func (ui *UI) showSelected() bool {
	switch ui.UserSelection {
	case "1":
		showCleanWallPostsUI()
		return true
	case "2":
		showCleanWallPostsCommentsUI()
		return true
	default:
		return false
	}
}

func (ui *UI) showMessageOfWrongInput() {
	fmt.Print("ERROR! Your input is wrong. Please try again...\n" +
		"> ")
}

func (ui *UI) showMessageOfExit() {
	fmt.Print("Exit...\n")
}
func showCleanWallPostsUI() {
	fmt.Println("[Cleaning of wallposts]\n" +
		"Here is empty yet....")
}

func showCleanWallPostsCommentsUI() {
	fmt.Println("[Cleaning of comments under wallposts]\n" +
		"Here is empty yet....")
}
