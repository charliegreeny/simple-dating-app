package match

import (
	"context"
	"github.com/charliegreeny/simple-dating-app/app"
	"github.com/charliegreeny/simple-dating-app/appctx"
	"slices"
)

type Matcher interface {
	Match(ctx context.Context, input *SwipeInput) (*SwipeOutput, error)
}

type swipe struct {
	store map[string][]string
}

func NewMatcher() Matcher {
	return &swipe{store: map[string][]string{}}
}

func (s swipe) Match(ctx context.Context, input *SwipeInput) (*SwipeOutput, error) {
	if !input.Swipe {
		return nil, nil
	}
	currentUser := appctx.GetUserFromCtx(ctx)
	if currentUser.ID == input.MatchID {
		return nil, app.ErrBadRequest{Message: "It is good to love oneself but you can not swipe self"}
	}
	go s.updateSwipes(currentUser.ID, input.MatchID)
	matchSwipes, ok := s.store[input.MatchID]
	if !ok {
		return &SwipeOutput{
			Matched: false,
			MatchID: input.MatchID,
		}, nil
	}
	if slices.Contains(matchSwipes, currentUser.ID) {
		return &SwipeOutput{
			Matched: true,
			MatchID: input.MatchID,
		}, nil
	}
	return &SwipeOutput{
		Matched: false,
		MatchID: input.MatchID,
	}, nil
}

func (s swipe) updateSwipes(currentId string, swipedID string) {
	swiped, ok := s.store[currentId]
	if !ok {
		swiped = []string{}
	}
	s.store[currentId] = append(swiped, swipedID)
}
