package event

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type Parser struct{}

type EventDetail struct {
	SourceIdentifier string
}

func (detail *EventDetail) UserName() string {
	return ""
}

type User struct {
	Name string
}

func (handler *Parser) Parse(ctx context.Context, ev events.CloudWatchEvent) (User, error) {
	ed := EventDetail{}
	user := User{}
	if err := json.Unmarshal(ev.Detail, &ed); err != nil {
		return user, fmt.Errorf("parse a request body detail: %w", err)
	}
	if ed.SourceIdentifier == "" {
		return user, errors.New(`request body is invalid. the field detail.SourceIdentifier is missing`)
	}
	user.Name = ed.UserName()
	return user, nil
}
