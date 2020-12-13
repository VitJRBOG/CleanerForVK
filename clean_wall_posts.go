package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	govkapi "github.com/VitJRBOG/GoVkApi"
)

// RunWallPostsCleaning запускает алгоритм удаления постов со стены
func RunWallPostsCleaning(accessToken string, ownerID, authorID int, msgChannel chan string) {
	var wpCleaner WallPostsCleaner
	wpCleaner.init(msgChannel, accessToken, ownerID, authorID)
	for {
		wpCleaner.requestWallPosts()
		wpCleaner.showProgress()
		wpCleaner.selectAuthorsWallPosts()
		itWasLastWallPosts := wpCleaner.checkEndOfWall()
		if itWasLastWallPosts {
			break
		} else {
			wpCleaner.enlargeOffset()
		}
	}
	wpCleaner.deleteAuthorsWallPosts()
}

// WallPostsCleaner хранит информацию для алгоритмов удаления постов со стены
type WallPostsCleaner struct {
	AccessToken        string
	OwnerID            int
	AuthorID           int
	NumberReqWallPosts int
	Offset             int
	WallPosts          []WallPost
	AuthorsWallPosts   []WallPost
	MsgChannel         chan string
}

func (w *WallPostsCleaner) init(msgChannel chan string, accessToken string, ownerID, authorID int) {
	w.MsgChannel = msgChannel
	w.AccessToken = accessToken
	w.OwnerID = ownerID
	w.AuthorID = authorID
	w.NumberReqWallPosts = 100
	w.Offset = 0
}

func (w *WallPostsCleaner) requestWallPosts() {
	paramsMap := map[string]string{
		"owner_id": strconv.Itoa(w.OwnerID),
		"filter":   "all",
		"count":    strconv.Itoa(w.NumberReqWallPosts),
		"offset":   strconv.Itoa(w.Offset),
		"v":        "5.126",
	}
	time.Sleep(335 * time.Millisecond)
	response, err := govkapi.SendRequestVkApi(w.AccessToken, "wall.get", paramsMap)
	if err != nil {
		panic(err.Error())
	}
	w.parseResponse(response)
}

func (w *WallPostsCleaner) parseResponse(response []byte) {
	var wallPosts []WallPost

	var f interface{}
	err := json.Unmarshal(response, &f)
	if err != nil {
		panic(err.Error())
	}
	valuesMap := f.(map[string]interface{})

	for _, itemMap := range valuesMap["items"].([]interface{}) {
		var wallPost WallPost

		item := itemMap.(map[string]interface{})

		wallPost.ID = int(item["id"].(float64))
		wallPost.OwnerID = int(item["owner_id"].(float64))
		if _, exist := item["signer_id"]; exist == true {
			wallPost.FromID = int(item["signer_id"].(float64))
		} else {
			wallPost.FromID = int(item["from_id"].(float64))
		}

		wallPosts = append(wallPosts, wallPost)
	}

	w.WallPosts = wallPosts
}

func (w *WallPostsCleaner) checkEndOfWall() bool {
	if len(w.WallPosts) >= w.NumberReqWallPosts {
		return false
	}
	return true
}

func (w *WallPostsCleaner) enlargeOffset() {
	if w.Offset == 0 {
		w.Offset++
	}
	w.Offset += w.NumberReqWallPosts
}

func (w *WallPostsCleaner) selectAuthorsWallPosts() {
	var authorsWallPosts []WallPost
	for i := 0; i < len(w.WallPosts); i++ {
		if w.AuthorID == w.WallPosts[i].FromID {
			authorsWallPosts = append(authorsWallPosts, w.WallPosts[i])
		}
	}
	w.AuthorsWallPosts = append(w.AuthorsWallPosts, authorsWallPosts...)
}

func (w *WallPostsCleaner) deleteAuthorsWallPosts() {
	if len(w.AuthorsWallPosts) > 0 {
		for i := 0; i < len(w.AuthorsWallPosts); i++ {
			paramsMap := map[string]string{
				"owner_id": strconv.Itoa(w.AuthorsWallPosts[i].OwnerID),
				"post_id":  strconv.Itoa(w.AuthorsWallPosts[i].ID),
				"v":        "5.126",
			}
			time.Sleep(335 * time.Millisecond)
			_, err := govkapi.SendRequestVkApi(w.AccessToken, "wall.delete", paramsMap)
			if err != nil {
				panic(err.Error())
			} else {
				msg := fmt.Sprintf("Wallpost https://vk.com/wall%d_%d"+
					" has been successfully deleted.",
					w.AuthorsWallPosts[i].OwnerID, w.AuthorsWallPosts[i].ID)
				w.MsgChannel <- msg
			}
		}
		w.MsgChannel <- "Done!"
	} else {
		w.MsgChannel <- "No wallposts from this author..."
	}
}

func (w *WallPostsCleaner) showProgress() {
	if len(w.WallPosts) > 0 {
		if w.Offset == 0 {
			w.MsgChannel <- fmt.Sprintf("Progress: %d wallposts has been viewed...",
				len(w.WallPosts)+w.Offset)
		} else {
			w.MsgChannel <- fmt.Sprintf("Progress: %d wallposts has been viewed...",
				len(w.WallPosts)+w.Offset-1)
		}
	}
}

// WallPost хранит информацию о посте со стены
type WallPost struct {
	ID      int
	OwnerID int
	FromID  int
}
