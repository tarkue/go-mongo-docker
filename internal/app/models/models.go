package models

type Link struct {
	Initial string `json:"initial"`
	Result  string `json:"result"`
	Counter int64  `json:"counter"`
}

type HandlerError struct {
	Error string `json:"error"`
}
