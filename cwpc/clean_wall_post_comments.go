package cwpc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/VitJRBOG/CleaningVk/cwp"
	govkapi "github.com/VitJRBOG/GoVkApi/v2"
)

// RunWallPostCommentsCleaning запускает алгоритмы для удаления комментариев из-под постов со стены
func RunWallPostCommentsCleaning(accessToken string, ownerID, authorID int,
	msgChannel chan string) {
	var wpcCleaner WallPostCommentsCleaner
	wpcCleaner.init(msgChannel, accessToken, ownerID, authorID)
	for {
		wpcCleaner.requestWallPosts()
		wpcCleaner.requestWallPostComments()
		wpcCleaner.showProgress()
		itWasLastWallPosts := wpcCleaner.checkEndOfWall()
		if itWasLastWallPosts {
			break
		} else {
			wpcCleaner.enlargeWallPostsOffset()
		}
	}
	if wpcCleaner.AuthorID == 0 {
		wpcCleaner.deleteWallPostComments(wpcCleaner.WallPostComments)
	} else {
		wpcCleaner.selectAuthorsWallPostComments()
		wpcCleaner.deleteWallPostComments(wpcCleaner.AuthorsWallPostComments)
	}
}

// WallPostCommentsCleaner хранит информацию для алгоритмов удаления комментариев из-под постов со стены
type WallPostCommentsCleaner struct {
	AccessToken               string
	OwnerID                   int
	AuthorID                  int
	NumberReqWallPosts        int
	NumberReqWallPostComments int
	WallPostsOffset           int
	WallPostCommentsOffset    int
	WallPosts                 []cwp.WallPost
	WallPostComments          []WallPostComment
	AuthorsWallPostComments   []WallPostComment
	MsgChannel                chan string
}

func (w *WallPostCommentsCleaner) init(msgChannel chan string, accessToken string, ownerID,
	authorID int) {
	w.MsgChannel = msgChannel
	w.AccessToken = accessToken
	w.OwnerID = ownerID
	w.AuthorID = authorID
	w.NumberReqWallPosts = 100
	w.NumberReqWallPostComments = 100
	w.WallPostsOffset = 0
	w.WallPostCommentsOffset = 0
}

func (w *WallPostCommentsCleaner) requestWallPosts() {
	paramsMap := map[string]string{
		"owner_id": strconv.Itoa(w.OwnerID),
		"filter":   "all",
		"count":    strconv.Itoa(w.NumberReqWallPosts),
		"offset":   strconv.Itoa(w.WallPostsOffset),
		"v":        "5.126",
	}
	time.Sleep(335 * time.Millisecond)
	response, err := govkapi.Method("wall.get", w.AccessToken, paramsMap)
	if err != nil {
		panic(err.Error())
	}
	w.parseWallPostsResponse(response)
}

func (w *WallPostCommentsCleaner) parseWallPostsResponse(response []byte) {
	var wallPosts []cwp.WallPost

	var f interface{}
	err := json.Unmarshal(response, &f)
	if err != nil {
		panic(err.Error())
	}
	valuesMap := f.(map[string]interface{})

	for _, itemMap := range valuesMap["items"].([]interface{}) {
		var wallPost cwp.WallPost

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

func (w *WallPostCommentsCleaner) checkEndOfWall() bool {
	if len(w.WallPosts) >= w.NumberReqWallPosts {
		return false
	}
	return true
}

func (w *WallPostCommentsCleaner) enlargeWallPostsOffset() {
	if w.WallPostsOffset == 0 {
		w.WallPostsOffset++
	}
	w.WallPostsOffset += w.NumberReqWallPosts
}

func (w *WallPostCommentsCleaner) requestWallPostComments() {
	for i := 0; i < len(w.WallPosts); i++ {
		for {
			paramsMap := map[string]string{
				"owner_id": strconv.Itoa(w.WallPosts[i].OwnerID),
				"post_id":  strconv.Itoa(w.WallPosts[i].ID),
				"count":    strconv.Itoa(w.NumberReqWallPostComments),
				"offset":   strconv.Itoa(w.WallPostCommentsOffset),
				"sort":     "desc",
				"v":        "5.68",
			}
			time.Sleep(335 * time.Millisecond)
			response, err := govkapi.Method("wall.getComments", w.AccessToken, paramsMap)
			if err != nil {
				panic(err.Error())
			}
			itWasLastWallPostComment := w.parseWallPostCommentsResponse(response, w.WallPosts[i])
			if itWasLastWallPostComment {
				break
			}
		}
		w.WallPostCommentsOffset = 0
	}
}

func (w *WallPostCommentsCleaner) parseWallPostCommentsResponse(response []byte,
	wallPost cwp.WallPost) bool {
	var wallPostComments []WallPostComment

	var f interface{}
	err := json.Unmarshal(response, &f)
	if err != nil {
		panic(err.Error())
	}
	valuesMap := f.(map[string]interface{})

	for _, itemMap := range valuesMap["items"].([]interface{}) {
		var wallPostComment WallPostComment
		wallPostComment.WallPostID = wallPost.ID
		wallPostComment.OwnerID = wallPost.OwnerID

		item := itemMap.(map[string]interface{})

		wallPostComment.ID = int(item["id"].(float64))
		wallPostComment.FromID = int(item["from_id"].(float64))

		wallPostComments = append(wallPostComments, wallPostComment)
	}
	itWasLastWallPostComment := w.checkEndOfPostComments(wallPostComments)
	if !(itWasLastWallPostComment) {
		w.enlargeWallPostCommentsOffset()
	}

	w.WallPostComments = append(w.WallPostComments, wallPostComments...)

	return itWasLastWallPostComment
}

func (w *WallPostCommentsCleaner) checkEndOfPostComments(wallPostComments []WallPostComment) bool {
	if len(wallPostComments) >= w.NumberReqWallPosts {
		return false
	}
	return true
}

func (w *WallPostCommentsCleaner) enlargeWallPostCommentsOffset() {
	w.WallPostCommentsOffset += w.NumberReqWallPostComments
}

func (w *WallPostCommentsCleaner) selectAuthorsWallPostComments() {
	var authorsWallPostComments []WallPostComment
	for i := 0; i < len(w.WallPostComments); i++ {
		if w.AuthorID == w.WallPostComments[i].FromID {
			authorsWallPostComments = append(authorsWallPostComments, w.WallPostComments[i])
		}
	}
	if len(authorsWallPostComments) > 0 {
		w.AuthorsWallPostComments = append(w.AuthorsWallPostComments, authorsWallPostComments...)
	} else {
		w.MsgChannel <- "No comments of wallpost from this author..."
	}
}

func (w *WallPostCommentsCleaner) deleteWallPostComments(wallPostComments []WallPostComment) {
	if len(wallPostComments) > 0 {
		for i := 0; i < len(wallPostComments); i++ {
			paramsMap := map[string]string{
				"owner_id":   strconv.Itoa(wallPostComments[i].OwnerID),
				"comment_id": strconv.Itoa(wallPostComments[i].ID),
				"v":          "5.126",
			}
			time.Sleep(335 * time.Millisecond)
			_, err := govkapi.Method("wall.deleteComment", w.AccessToken, paramsMap)
			if err != nil {
				panic(err.Error())
			} else {
				msg := fmt.Sprintf("Comment https://vk.com/wall%d_%d?reply=%d"+
					" has been successfully deleted.",
					wallPostComments[i].OwnerID, wallPostComments[i].WallPostID,
					wallPostComments[i].ID)
				w.MsgChannel <- msg
			}
		}
		w.MsgChannel <- "Done!"
	}
}

func (w *WallPostCommentsCleaner) showProgress() {
	if len(w.WallPosts) > 0 {
		if w.WallPostsOffset == 0 {
			w.MsgChannel <- fmt.Sprintf("Progress: %d wallposts has been viewed...",
				len(w.WallPosts)+w.WallPostsOffset)
		} else {
			w.MsgChannel <- fmt.Sprintf("Progress: %d wallposts has been viewed...",
				len(w.WallPosts)+w.WallPostsOffset-1)
		}
	}
}

// WallPostComment хранит информацию о комментарии из-под поста со стены
type WallPostComment struct {
	ID         int
	OwnerID    int
	WallPostID int
	FromID     int
}
