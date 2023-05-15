package apis

import (
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
		{"valid input", args{"6452183b3aa8ab565e89897b", []lsdb.LotteryEventInfo{{EventUID: stringToPrimitive("6452183b3aa8ab565e89897b"), EventDate: 1683244800000, Name: "Friday Bonanza", EventType: "FB", WinningNumber: []int{67, 23, 65, 22, 11}, CreatedAt: 1683101755340, UpdatedAt: 1683101755340}}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := eventIDExist(tt.args.eventID, tt.args.eventIDArray); got != tt.want {
				t.Errorf("eventIDExist() = %v, want %v", got, tt.want)
			}
		})
	}
}
