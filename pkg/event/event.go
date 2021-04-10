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
	RequestParameters RequestParameters `json:"requestParameters"`
}

type RequestParameters struct {
	UserName string `json:"userName"`
}

func (detail *EventDetail) UserName() string {
	return detail.RequestParameters.UserName
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
	user.Name = ed.UserName()
	if user.Name == "" {
		return user, errors.New(`request body is invalid. created user name is missing`)
	}
	return user, nil
}
