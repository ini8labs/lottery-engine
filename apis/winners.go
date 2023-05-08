package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ini8labs/lsdb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s Server) validateEventId(str string) (bool, error) {
	eventIdExist := true

	resp, err := s.GetAllEvents()
	if err != nil {
		s.Logger.Error(err.Error())
		return false, err
	}

	for i := 0; i < len(resp); i++ {
		if resp[i].EventUID == stringToPrimitive(str) {
			eventIdExist = true
			break
		}
		if resp[i].EventUID != stringToPrimitive(str) {
			eventIdExist = false
		}
	}
	return eventIdExist, nil
}

func winnerDecider(betNumbers []int, winNumbers []int, amount int) (string, int) {
	count := 0

	for i := 0; i < len(winNumbers); i++ {
		for j := 0; j < len(betNumbers); j++ {
			if winNumbers[i] == betNumbers[j] {
				count++
			}
		}
	}

	amountWon := 0
	winType := ""
	switch {
	case count == 1 && len(betNumbers) == 1:
		amountWon = amount*40 - amount
		winType = "Direct-1"

	case count == 2 && len(betNumbers) == 2:
		amountWon = amount*240 - amount
		winType = "Direct-2"

	case count == 3 && len(betNumbers) == 3:
		amountWon = amount*2100 - amount
		winType = "Direct-3"

	case count == 4 && len(betNumbers) == 4:
		amountWon = amount*6000 - amount
		winType = "Direct-4"

	case count == 5 && len(betNumbers) == 5:
		amountWon = amount*44000 - amount
		winType = "Direct-5"

	case count == 2 && len(betNumbers) != 2:
		amountWon = amount*240 - amount
		winType = "Perm-2"

	case count == 3 && len(betNumbers) != 3:
		amountWon = amount*2100 - amount
		winType = "Perm-3"

	case count == 4 && len(betNumbers) != 4:
		amountWon = amount*6000 - amount
		winType = "Perm-4"

	case count == 5 && len(betNumbers) != 5:
		amountWon = amount*44000 - amount
		winType = "Perm-5"
	}
	return winType, amountWon
}

func (s Server) winnerSelector(eventId primitive.ObjectID) ([]lsdb.WinnerInfo, error) {

	resp, err := s.Client.GetParticipantsInfoByEventID(eventId)
	if err != nil {
		s.Logger.Error(err.Error())
		return []lsdb.WinnerInfo{}, err
	}

	resp1, err := s.Client.GetEventInfoByEventId(eventId)

	if err != nil {
		return []lsdb.WinnerInfo{}, err
	}
	arr := initializeEventWinnerInfo(resp, resp1)

	return arr, nil
}

func initializeEventWinnerInfo(resp []lsdb.EventParticipantInfo, resp1 []lsdb.LotteryEventInfo) []lsdb.WinnerInfo {
	winNumbers := resp1.WinningNumber
	var arr []lsdb.WinnerInfo

	for i := 0; i < len(resp); i++ {
		betNumbers := resp[i].BetNumbers

		winType, amountWon := winnerDecider(betNumbers, winNumbers, resp[i].Amount)

		winnerInfo := lsdb.WinnerInfo{
			EventID:   resp[i].EventUID,
			UserID:    resp[i].UserID,
			WinType:   winType,
			AmountWon: amountWon,
		}
		arr = append(arr, winnerInfo)
	}
	return arr
}

func (s Server) addNewWinner(c *gin.Context) {
	eventId := c.Param("eventId")

	validation, _ := s.validateEventId(eventId)
	if !validation {
		c.JSON(http.StatusBadRequest, "EventId does not exist")
		s.Logger.Error("invalid event id")
		return
	}

	resp, err := s.winnerSelector(stringToPrimitive(eventId))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		s.Logger.Error(err.Error())
		return
	}

	for i := 0; i < len(resp); i++ {
		if err := s.AddNewWinner(resp[i]); err != nil {
			c.JSON(http.StatusInternalServerError, "something is wrong with the server")
			s.Logger.Error(err.Error())
			return
		}
	}
	c.JSON(http.StatusCreated, "Winner added successfully")
}

func (s Server) GetEventWinners(c *gin.Context) {
	eventid := c.Query("eventid")

	validation, _ := s.validateEventId(eventid)
	if !validation {
		c.JSON(http.StatusBadRequest, "EventId does not exist")
		s.Logger.Error("invalid event id")
		return
	}

	resp, err := s.Client.GetEventWinners(stringToPrimitive(eventid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "something is wrong with the server")
		s.Logger.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
