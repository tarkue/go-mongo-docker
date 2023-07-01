package service

import (
	"errors"
	"regexp"

	"github.com/labstack/echo/v4"
)

const WithProtocolURL = `^https?:\/\/(?:www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b(?:[-a-zA-Z0-9()@:%_\+.~#?&\/=]*)$`
const WithoutProtocolURL = `^[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b(?:[-a-zA-Z0-9()@:%_\+.~#?&//=]*)$`

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) MatchLink(ctx echo.Context, link string) error {
	matchedProtocol, _ := regexp.MatchString(WithProtocolURL, link)
	matchedWithoutProtocol, _ := regexp.MatchString(WithoutProtocolURL, link)

	if !matchedProtocol && !matchedWithoutProtocol {
		return errors.New("not matched")
	}

	return nil
}
