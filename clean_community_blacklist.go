package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	govkapi "github.com/VitJRBOG/GoVkApi"
)

// RunCommunityBlacklistCleaning запускает алгоритмы удаления субъектов из черного списка сообщества
func RunCommunityBlacklistCleaning(accessToken string, ownerID int, msgChannel chan string) {
	var cbCleaner CommunityBlacklistCleaner
	cbCleaner.init(accessToken, ownerID, msgChannel)
	for {
		thisWasLastBannedSubjects := cbCleaner.requestBanned()

		if thisWasLastBannedSubjects {
			break
		} else {
			cbCleaner.showProgress()
			cbCleaner.enlargeOffset()
		}
	}

	if len(cbCleaner.BannedSubjects) > 0 {
		cbCleaner.unbanBanned()
	} else {
		cbCleaner.MsgChannel <- "No banned subjects in the blacklist..."
	}
}

// CommunityBlacklistCleaner хранит информацию для алгоритмов удаления субъектов из
// черного списка сообщества
type CommunityBlacklistCleaner struct {
	AccessToken             string
	OwnerID                 int
	NumberReqBannedSubjects int
	Offset                  int
	BannedSubjects          []BannedSubject
	MsgChannel              chan string
}

func (c *CommunityBlacklistCleaner) init(accessToken string, ownerID int, msgChannel chan string) {
	c.MsgChannel = msgChannel
	c.AccessToken = accessToken
	c.OwnerID = ownerID
	c.NumberReqBannedSubjects = 2
	c.Offset = 0
}

func (c *CommunityBlacklistCleaner) requestBanned() bool {
	paramsMap := map[string]string{
		"group_id": strconv.Itoa(c.OwnerID),
		"count":    strconv.Itoa(c.NumberReqBannedSubjects),
		"offset":   strconv.Itoa(c.Offset),
		"v":        "5.126",
	}
	time.Sleep(335 * time.Millisecond)
	response, err := govkapi.SendRequestVkApi(c.AccessToken, "groups.getBanned", paramsMap)
	if err != nil {
		panic(err.Error())
	}
	thisWasLastBannedSubjects := c.parseBannedResponse(response)
	return thisWasLastBannedSubjects
}

func (c *CommunityBlacklistCleaner) parseBannedResponse(response []byte) bool {
	var bannedSubjects []BannedSubject

	var f interface{}
	err := json.Unmarshal(response, &f)
	if err != nil {
		panic(err.Error())
	}
	valuesMap := f.(map[string]interface{})

	for _, itemMap := range valuesMap["items"].([]interface{}) {
		var bannedSubject BannedSubject

		item := itemMap.(map[string]interface{})

		bannedSubject.Type = item["type"].(string)

		bannedSubject.ID = int(item[bannedSubject.Type].(map[string]interface{})["id"].(float64))

		if bannedSubject.Type == "group" {
			bannedSubject.ID = -(bannedSubject.ID)
		}

		bannedSubjects = append(bannedSubjects, bannedSubject)
	}

	c.BannedSubjects = append(c.BannedSubjects, bannedSubjects...)
	thisWasLastBannedSubjects := c.checkEndOfBannedList(bannedSubjects)
	return thisWasLastBannedSubjects
}

func (c *CommunityBlacklistCleaner) checkEndOfBannedList(bannedSubjects []BannedSubject) bool {
	if len(bannedSubjects) == 0 {
		return true
	}
	return false
}

func (c *CommunityBlacklistCleaner) enlargeOffset() {
	c.Offset += c.NumberReqBannedSubjects
}

func (c *CommunityBlacklistCleaner) unbanBanned() {
	if len(c.BannedSubjects) > 0 {
		for i := 0; i < len(c.BannedSubjects); i++ {
			paramsMap := map[string]string{
				"group_id": strconv.Itoa(c.OwnerID),
				"owner_id": strconv.Itoa(c.BannedSubjects[i].ID),
				"v":        "5.126",
			}
			time.Sleep(335 * time.Millisecond)
			_, err := govkapi.SendRequestVkApi(c.AccessToken, "groups.unban", paramsMap)
			if err != nil {
				panic(err.Error())
			} else {
				switch c.BannedSubjects[i].Type {
				case "profile":
					msg := fmt.Sprintf("User https://vk.com/id%d"+
						" has been successfully unbaned.", c.BannedSubjects[i].ID)
					c.MsgChannel <- msg
				case "group":
					msg := fmt.Sprintf("Community https://vk.com/club%d"+
						" has been successfully unbaned.", -(c.BannedSubjects[i].ID))
					c.MsgChannel <- msg
				}
			}
		}
		c.MsgChannel <- "Done!"
	}
}

func (c *CommunityBlacklistCleaner) showProgress() {
	if len(c.BannedSubjects) > 0 {
		c.MsgChannel <- fmt.Sprintf("Progress: %d banned subjects has been viewed...",
			len(c.BannedSubjects))
	}
}

// BannedSubject хранит информацию о заблокированном субъекте
type BannedSubject struct {
	Type string
	ID   int
}
