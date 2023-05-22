package apis

import (
	"reflect"
	"testing"

	"github.com/ini8labs/lsdb"
)

func Test_countMatchNumber(t *testing.T) {
	type args struct {
		betNumbers []int
		winNumbers []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"1 matched", args{[]int{1, 2, 3, 4, 5}, []int{5, 6, 7, 8, 9}}, 1},
		{"2 matched", args{[]int{1, 2, 3, 4, 5}, []int{4, 5, 6, 7, 8}}, 2},
		{"3 matched", args{[]int{1, 2, 3, 4, 5}, []int{3, 4, 5, 6, 7}}, 3},
		{"4 matched", args{[]int{1, 2, 3, 4, 5}, []int{2, 3, 4, 5, 6}}, 4},
		{"5 matched", args{[]int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countMatchNumber(tt.args.betNumbers, tt.args.winNumbers); got != tt.want {
				t.Errorf("countMatchNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_winnerDecider(t *testing.T) {
	type args struct {
		betNumbers []int
		winNumbers []int
		amount     int
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 int
	}{
		{"direct-1 win", args{[]int{1}, []int{1, 2, 3, 4, 5}, 10}, "Direct-1", 390},
		{"direct-2 win", args{[]int{1, 2}, []int{1, 2, 3, 4, 5}, 10}, "Direct-2", 2390},
		{"direct-3 win", args{[]int{1, 2, 3}, []int{1, 2, 3, 4, 5}, 10}, "Direct-3", 20990},
		{"direct-4 win", args{[]int{1, 2, 3, 4}, []int{1, 2, 3, 4, 5}, 10}, "Direct-4", 59990},
		{"direct-5 win", args{[]int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}, 10}, "Direct-5", 439990},
		{"perm-2 win", args{[]int{1, 2, 6}, []int{1, 2, 3, 4, 5}, 10}, "Perm-2", 2390},
		{"perm-3 win", args{[]int{1, 2, 3, 7}, []int{1, 2, 3, 4, 5}, 10}, "Perm-3", 20990},
		{"perm-4 win", args{[]int{1, 2, 3, 4, 8}, []int{1, 2, 3, 4, 5}, 10}, "Perm-4", 59990},
		{"perm-5 win", args{[]int{1, 2, 3, 4, 5, 9}, []int{1, 2, 3, 4, 5}, 10}, "Perm-5", 439990},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := winnerDecider(tt.args.betNumbers, tt.args.winNumbers, tt.args.amount)
			if got != tt.want {
				t.Errorf("winnerDecider() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("winnerDecider() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_eventIDExist(t *testing.T) {
	type args struct {
		eventID      string
		eventIDArray []lsdb.LotteryEventInfo
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"valid input",
			args{
				"6452183b3aa8ab565e89897b",
				[]lsdb.LotteryEventInfo{
					{
						EventUID:      stringToPrimitive("6452183b3aa8ab565e89897b"),
						EventDate:     1683244800000,
						Name:          "Friday Bonanza",
						EventType:     "FB",
						WinningNumber: []int{67, 23, 65, 22, 11},
						CreatedAt:     1683101755340,
						UpdatedAt:     1683101755340},
				},
			},
			true,
		},
		{
			"invalid input",
			args{
				"6452183b3aa8ab565e89897b",
				[]lsdb.LotteryEventInfo{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := eventIDExist(tt.args.eventID, tt.args.eventIDArray); got != tt.want {
				t.Errorf("eventIDExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initializeWinnersInfo(t *testing.T) {
	type args struct {
		eventWinnerInfo      []lsdb.WinnerInfo
		eventParticipantInfo []lsdb.EventParticipantInfo
	}
	tests := []struct {
		name string
		args args
		want []Winners
	}{
		{
			"valid input",
			args{
				[]lsdb.WinnerInfo{
					{
						EventID:   stringToPrimitive("6452183b3aa8ab565e89897b"),
						UserID:    stringToPrimitive("644790a68e3540cbb44180b0"),
						WinType:   "Perm-2",
						AmountWon: 1314500,
						CreatedAt: 1683802437767,
					},
					{
						EventID:   stringToPrimitive("6452183b3aa8ab565e89897b"),
						UserID:    stringToPrimitive("644790a68e3540cbb44180b0"),
						WinType:   "Direct-3",
						AmountWon: 11544500,
						CreatedAt: 1683809942653,
					},
				},
				[]lsdb.EventParticipantInfo{
					{
						BetUID:   stringToPrimitive("64529784e5b433802324b3f7"),
						EventUID: stringToPrimitive("6452183b3aa8ab565e89897b"),
						ParticipantInfo: lsdb.ParticipantInfo{
							UserID:     stringToPrimitive("644790a68e3540cbb44180b0"),
							BetNumbers: []int{2, 22, 62},
							Amount:     5500},
						CreatedAt: 1683134340492,
						UpdatedAt: 1683134340492,
					},
				},
			},
			[]Winners{
				{
					UserID:    "644790a68e3540cbb44180b0",
					EventUID:  "6452183b3aa8ab565e89897b",
					AmountWon: 1314500,
					WinType:   "Perm-2",
					BetID:     "64529784e5b433802324b3f7",
				},
				{
					UserID:    "644790a68e3540cbb44180b0",
					EventUID:  "6452183b3aa8ab565e89897b",
					AmountWon: 11544500,
					WinType:   "Direct-3",
					BetID:     "64529784e5b433802324b3f7",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := initializeWinnersInfo(tt.args.eventWinnerInfo, tt.args.eventParticipantInfo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initializeWinnersInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initializeEventWinnerInfo(t *testing.T) {
	type args struct {
		eventParticipantInfoArr []lsdb.EventParticipantInfo
	}

	arr := []lsdb.EventParticipantInfo{
		{
			BetUID:   stringToPrimitive("64529784e5b433802324b3f7"),
			EventUID: stringToPrimitive("6452183b3aa8ab565e89897b"),
			ParticipantInfo: lsdb.ParticipantInfo{
				UserID:     stringToPrimitive("644790a68e3540cbb44180b0"),
				BetNumbers: []int{12, 29, 62}, Amount: 50},
			CreatedAt: 1683134340492, UpdatedAt: 1683134340492},
	}

	tests := []struct {
		name string
		args args
		want []lsdb.WinnerInfo
	}{
		{
			"valid input",
			args{arr},
			[]lsdb.WinnerInfo{
				{
					EventID:   stringToPrimitive("6452183b3aa8ab565e89897b"),
					UserID:    stringToPrimitive("644790a68e3540cbb44180b0"),
					WinType:   "Direct-3",
					AmountWon: 104950,
					CreatedAt: 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := initializeEventWinnerInfo(tt.args.eventParticipantInfoArr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initializeEventWinnerInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
