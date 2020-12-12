package main

import (
	"fmt"
	"strconv"
)

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

type CleanWallPostsUI struct {
	AccessToken string
	OwnerID     int
	AuthorID    int
	msgChannel  chan string
}

func (c *CleanWallPostsUI) init() {
	c.msgChannel = make(chan string)
}

func (c *CleanWallPostsUI) setAccessToken() {
	fmt.Print("--- Enter your access token and press «Enter» ---\n" +
		"> ")
	var accessToken string
	_, err := fmt.Scan(&accessToken)
	if err != nil {
		panic(err.Error())
	}
	c.AccessToken = accessToken
}

func (c *CleanWallPostsUI) setOwnerID() {
	fmt.Print("--- Now enter ID of owner of wall posts and press «Enter» ---\n" +
		"> ")
	var ownerID string
	_, err := fmt.Scan(&ownerID)
	if err != nil {
		panic(err.Error())
	}
	c.OwnerID, err = strconv.Atoi(ownerID)
	if err != nil {
		panic(err.Error())
	}
}

func (c *CleanWallPostsUI) setAuthorID() {
	fmt.Print("--- And enter ID of author of wall posts and press «Enter» ---\n" +
		"> ")
	var authorID string
	_, err := fmt.Scan(&authorID)
	if err != nil {
		panic(err.Error())
	}
	c.AuthorID, err = strconv.Atoi(authorID)
	if err != nil {
		panic(err.Error())
	}
}

func (c *CleanWallPostsUI) outputtingMessages() {
	for {
		msg := <-c.msgChannel
		fmt.Printf("%v\n", msg)
		if msg == "Done!" || msg == "No wall posts from this author..." {
			break
		}
	}
}

func showCleanWallPostsUI() {
	var cwpUI CleanWallPostsUI
	cwpUI.init()
	fmt.Print("[Cleaning of wallposts]\n")
	cwpUI.setAccessToken()
	cwpUI.setOwnerID()
	cwpUI.setAuthorID()
	go RunWallPostsCleaning(cwpUI.AccessToken, cwpUI.OwnerID, cwpUI.AuthorID, cwpUI.msgChannel)
	cwpUI.outputtingMessages()
}

type CleanWallPostCommentsUI struct {
	AccessToken string
	OwnerID     int
	AuthorID    int
	msgChannel  chan string
}

func (c *CleanWallPostCommentsUI) init() {
	c.msgChannel = make(chan string)
}

func (c *CleanWallPostCommentsUI) setAccessToken() {
	fmt.Print("--- Enter your access token and press «Enter» ---\n" +
		"> ")
	var accessToken string
	_, err := fmt.Scan(&accessToken)
	if err != nil {
		panic(err.Error())
	}
	c.AccessToken = accessToken
}

func (c *CleanWallPostCommentsUI) setOwnerID() {
	fmt.Print("--- Now enter ID of owner of wall post comments and press «Enter» ---\n" +
		"> ")
	var ownerID string
	_, err := fmt.Scan(&ownerID)
	if err != nil {
		panic(err.Error())
	}
	c.OwnerID, err = strconv.Atoi(ownerID)
	if err != nil {
		panic(err.Error())
	}
}

func (c *CleanWallPostCommentsUI) setAuthorID() {
	fmt.Print("--- And enter ID of author of wall post comments and press «Enter» ---\n" +
		"> ")
	var authorID string
	_, err := fmt.Scan(&authorID)
	if err != nil {
		panic(err.Error())
	}
	c.AuthorID, err = strconv.Atoi(authorID)
	if err != nil {
		panic(err.Error())
	}
}

func (c *CleanWallPostCommentsUI) outputtingMessages() {
	for {
		msg := <-c.msgChannel
		fmt.Printf("%v\n", msg)
		if msg == "Done!" || msg == "No wall post comments from this author..." {
			break
		}
	}
}

func showCleanWallPostsCommentsUI() {
	var cwpcUI CleanWallPostCommentsUI
	cwpcUI.init()
	fmt.Print("[Cleaning of comments under wallposts]\n")
	cwpcUI.setAccessToken()
	cwpcUI.setOwnerID()
	cwpcUI.setAuthorID()
	go RunWallPostCommentsCleaning(cwpcUI.AccessToken, cwpcUI.OwnerID, cwpcUI.AuthorID, cwpcUI.msgChannel)
	cwpcUI.outputtingMessages()
}
