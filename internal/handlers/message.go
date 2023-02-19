package handlers

import "net/http"

type StatusMessageError interface {
	error
	Status() int
	Message() string
}

type httpResponseMessage struct {
	status  int
	message string
}

func (mesg httpResponseMessage) Status() int {
	return mesg.status
}

func (mesg httpResponseMessage) Message() string {
	return mesg.message
}

func (mesg httpResponseMessage) Error() string {
	return mesg.message
}

var couldNotDecodeError = httpResponseMessage{
	status:  http.StatusInternalServerError,
	message: "could not decode request body",
}

func createInvalidBodyError(message string) httpResponseMessage {
	return httpResponseMessage{
		status:  http.StatusBadRequest,
		message: message,
	}
}

func createInternalError(message string) httpResponseMessage {
	return httpResponseMessage{
		status:  http.StatusInternalServerError,
		message: message,
	}
}
