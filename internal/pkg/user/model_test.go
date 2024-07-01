package user

import (
	"reflect"
	"testing"
	"time"
)

func setDob(dob string) time.Time {
	t, _ := time.Parse(time.RFC3339, dob)
	return t
}

func TestEntity_toOutput(t *testing.T) {
	tests := []struct {
		name   string
		entity *Entity
		want   *Output
	}{
		{
			name: "successfully convert to output with correct age with d.o.b at end of the year",
			entity: &Entity{
				ID:          "123",
				Name:        "Charlie G",
				Gender:      "Male",
				DateOfBirth: setDob("1997-12-31T00:00:00Z"),
				Email:       "test@email.com",
				Password:    "Password",
			},
			want: &Output{
				ID:       "123",
				Email:    "test@email.com",
				Name:     "Charlie G",
				Gender:   "Male",
				Age:      26,
				Password: "Password",
			},
		},
		{
			name: "successfully convert to output with correct age with d.o.b at start of the year",
			entity: &Entity{
				ID:          "123",
				Name:        "Charlie G",
				Gender:      "Male",
				DateOfBirth: setDob("1997-01-01T00:00:00Z"),
				Email:       "test@email.com",
				Password:    "Password",
			},
			want: &Output{
				ID:       "123",
				Email:    "test@email.com",
				Name:     "Charlie G",
				Gender:   "Male",
				Age:      27,
				Password: "Password",
			},
		},
		{
			name: "successfully convert to output with correct age with at d.o.b on a leap day",
			entity: &Entity{
				ID:          "123",
				Name:        "Charlie G",
				Gender:      "Male",
				DateOfBirth: setDob("1996-02-29T00:00:00Z"),
				Email:       "test@email.com",
				Password:    "Password",
			},
			want: &Output{
				ID:       "123",
				Email:    "test@email.com",
				Name:     "Charlie G",
				Gender:   "Male",
				Age:      28,
				Password: "Password",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.entity.toOutput(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}
