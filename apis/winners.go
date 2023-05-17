package apis

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ini8labs/lsdb"
)

func (s Server) validateEventId(eventId string) bool {

	resp, err := s.GetAllEvents()
	if err != nil {
		s.Logger.Error(err.Error())
		return false
	}

	return eventIDExist(eventId, resp)
}

func eventIDExist(eventID string, eventIDArray []lsdb.LotteryEventInfo) bool {
	eventIdPrimitive := stringToPrimitive(eventID)

	for i := 0; i < len(eventIDArray); i++ {
		if eventIDArray[i].EventUID == eventIdPrimitive {
			return true
		}
	}

	return false
}

func countMatchNumber(betNumbers []int, winNumbers []int) int {
	count := 0
	for i := 0; i < len(winNumbers); i++ {
		for j := 0; j < len(betNumbers); j++ {
			if winNumbers[i] == betNumbers[j] {
				count++
			}
		}
	}
	return count
}

func winnerDecider(betNumbers []int, winNumbers []int, amount int) (string, int) {

	matchNumber := countMatchNumber(betNumbers, winNumbers)

	amountWon := 0
	winType := ""
	switch {
	case matchNumber == 1 && len(betNumbers) == 1:
		amountWon = amount*40 - amount
		winType = "Direct-1"

	case matchNumber == 2 && len(betNumbers) == 2:
		amountWon = amount*240 - amount
		winType = "Direct-2"

	case matchNumber == 3 && len(betNumbers) == 3:
		amountWon = amount*2100 - amount
		winType = "Direct-3"

	case matchNumber == 4 && len(betNumbers) == 4:
		amountWon = amount*6000 - amount
		winType = "Direct-4"

	case matchNumber == 5 && len(betNumbers) == 5:
		amountWon = amount*44000 - amount
		winType = "Direct-5"

	case matchNumber == 2 && len(betNumbers) != 2:
		amountWon = amount*240 - amount
		winType = "Perm-2"

	case matchNumber == 3 && len(betNumbers) != 3:
		amountWon = amount*2100 - amount
		winType = "Perm-3"

	case matchNumber == 4 && len(betNumbers) != 4:
		amountWon = amount*6000 - amount
		winType = "Perm-4"

	case matchNumber == 5 && len(betNumbers) != 5:
		amountWon = amount*44000 - amount
		winType = "Perm-5"
	}

	return winType, amountWon
}

func (s Server) winnerSelector(eventId primitive.ObjectID) ([]lsdb.WinnerInfo, error) {

	resp, err := s.Client.GetParticipantsInfoByEventID(eventId)
	if err != nil {
		s.Logger.Error(err)
		return []lsdb.WinnerInfo{}, err
	}

	if resp == nil {
		return []lsdb.WinnerInfo{}, err
	}

	// resp1, err := s.Client.GetEventInfoByEventId(eventId)

	// if err != nil {
	// 	return []lsdb.WinnerInfo{}, err
	// }
	eventWinnersInfoArr := initializeEventWinnerInfo(resp)

	return eventWinnersInfoArr, nil
}

func initializeEventWinnerInfo(eventParticipantInfoArray []lsdb.EventParticipantInfo) []lsdb.WinnerInfo {
	var winNumbers = []int{81, 30, 3}
	var winnerInfoArr []lsdb.WinnerInfo

	for i := 0; i < len(eventParticipantInfoArray); i++ {
		betNumbers := eventParticipantInfoArray[i].BetNumbers

		winType, amountWon := winnerDecider(betNumbers, winNumbers, eventParticipantInfoArray[i].Amount)
		winnerInfo := lsdb.WinnerInfo{
			EventID:   eventParticipantInfoArray[i].EventUID,
			UserID:    eventParticipantInfoArray[i].UserID,
			WinType:   winType,
			AmountWon: amountWon,
		}
		winnerInfoArr = append(winnerInfoArr, winnerInfo)

	}
	fmt.Println(eventParticipantInfoArray)
	return winnerInfoArr
}

func (s Server) generateWinners(c *gin.Context) {
	eventId := c.Query("eventId")

	valid := s.validateEventId(eventId)
	if !valid {
		c.JSON(http.StatusNotFound, "EventId does not exist")
		s.Logger.Error("event id does not exist")
		return
	}

	resp, err := s.winnerSelector(stringToPrimitive(eventId))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		s.Logger.Error(err)
		return
	}

	for i := 0; i < len(resp); i++ {
		if err := s.AddNewWinner(resp[i]); err != nil {
			c.JSON(http.StatusInternalServerError, "something is wrong with the server")
			s.Logger.Error(err)
			return
		}
	}
	c.JSON(http.StatusCreated, "Winners added successfully")
}

func (s Server) getEventWinners(c *gin.Context) {
	eventId := c.Query("eventId")

	valid := s.validateEventId(eventId)
	if !valid {
		c.JSON(http.StatusBadRequest, "EventId does not exist")
		s.Logger.Error("invalid event id")
		return
	}

	resp, err := s.GetEventWinners(stringToPrimitive(eventId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "something is wrong with the server")
		s.Logger.Error(err)
		return
	}

	resp1, err := s.GetParticipantsInfoByEventID(stringToPrimitive(eventId))
	if err != nil {
		s.Logger.Error(err)
		return
	}

	winnerInfoArr := initializeWinnersInfo(resp, resp1)
	c.JSON(http.StatusOK, winnerInfoArr)
}

func initializeWinnersInfo(eventWinnerInfo []lsdb.WinnerInfo, eventParticipantInfo []lsdb.EventParticipantInfo) []Winners {
	var winnerInfoArr []Winners

	for i := 0; i < len(eventWinnerInfo); i++ {
		for j := 0; j < len(eventParticipantInfo); j++ {
			winnerInfo := Winners{
				EventUID:  primitiveToString(eventWinnerInfo[i].EventID),
				UserID:    primitiveToString(eventWinnerInfo[i].UserID),
				AmountWon: eventWinnerInfo[i].AmountWon,
				WinType:   eventWinnerInfo[i].WinType,
				BetID:     primitiveToString(eventParticipantInfo[j].BetUID),
			}
			winnerInfoArr = append(winnerInfoArr, winnerInfo)
		}
	}

	return winnerInfoArr
}
