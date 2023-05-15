package apis

import (
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_stringToPrimitive(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want primitive.ObjectID
	}{
		{"valid input", args{"645cc94575799d46e3352e0b"}, stringToPrimitive("645cc94575799d46e3352e0b")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringToPrimitive(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stringToPrimitive() = %v, want %v", got, tt.want)
			}
		})
	}
}
