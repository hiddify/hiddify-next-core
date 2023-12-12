package global

import (
	"encoding/json"
)

var status = Stopped

type StatusMessage struct {
	Status  string  `json:"status"`
	Alert   *string `json:"alert"`
	Message *string `json:"message"`
}

func propagateStatus(newStatus string) string {
	status = newStatus

	msg, _ := json.Marshal(StatusMessage{Status: status})
	return string(msg)
}

func stopAndAlert(alert string, err error) string {
	status = Stopped
	message := err.Error()

	msg, _ := json.Marshal(StatusMessage{Status: status, Alert: &alert, Message: &message})
	return string(msg)
}
