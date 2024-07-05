package filter

import (
	"context"
	"github.com/charliegreeny/simple-dating-app/app"
	"github.com/stretchr/testify/assert"
	"testing"
)

var maleUser1 = &app.User{UserOutput: &app.UserOutput{
	ID:     "M1",
	Gender: app.Male,
}}

var maleUser2 = &app.User{UserOutput: &app.UserOutput{
	ID:     "M2",
	Gender: app.Male,
}}

var femaleUser1 = &app.User{UserOutput: &app.UserOutput{
	ID:     "F1",
	Gender: app.Female,
}}
var femaleUser2 = &app.User{UserOutput: &app.UserOutput{
	ID:     "F2",
	Gender: app.Female,
}}

func Test_gender_Apply(t *testing.T) {

	tests := []struct {
		name       string
		perf       *app.Preference
		wantGender string
		wantLen    int
	}{
		{
			name:       "prefer females returns only females",
			perf:       &app.Preference{Gender: app.Female},
			wantGender: app.Female,
			wantLen:    2,
		},
		{
			name:       "prefer males returns only males",
			perf:       &app.Preference{Gender: app.Male},
			wantGender: app.Male,
			wantLen:    2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gender{}
			stubUsers := []*app.User{maleUser1, maleUser2, femaleUser1, femaleUser2}
			got := g.Apply(context.Background(), tt.perf, stubUsers)
			assert.Equal(t, tt.wantLen, len(got))
			for _, user := range got {
				assert.Equal(t, tt.wantGender, user.Gender)
			}
		})
	}
}
