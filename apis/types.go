package apis

import (
	"github.com/ini8labs/lsdb"
	"github.com/sirupsen/logrus"
)

type Server struct {
	*logrus.Logger
	*lsdb.Client
	Addr string
}

type Date struct {
	Day   int `json:"day,omitempty"`
	Month int `json:"month,omitempty"`
	Year  int `json:"year,omitempty"`
}

type AddWinner struct {
	UserID    string `json:"user_id"`
	EventUID  string `json:"event_id"`
	AmountWon int    `json:"amountWon"`
	EventDate Date   `json:"event_date"`
	WinType   string `json:"winType"`
	BetID     string `json:"bet_id"`
}
