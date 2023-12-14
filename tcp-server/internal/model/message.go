package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

const (
	MessageTypeChallenge = "challenge"
	MessageTypeWow       = "wow"
	MessageTypeSolution  = "solution"
	MessageTypeRequest   = "request"
)

var (
	messageTypes = map[string]bool{
		MessageTypeChallenge: true,
		MessageTypeWow:       true,
		MessageTypeSolution:  true,
		MessageTypeRequest:   true,
	}
)

type Message struct {
	RequestID     string `json:"request_id"`
	MessageType   string `json:"message_type"`
	MessageString string `json:"message_string"`
	Difficulty    int    `json:"difficulty"`
}

func PrepareMessage(rid string, mType string, mString string, d int) Message {
	return Message{
		RequestID:     rid,
		MessageType:   mType,
		MessageString: mString,
		Difficulty:    d,
	}
}

func (m Message) AsJsonString() []byte {
	result, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("Error while marshalling Message: %v", err)
		return nil
	}
	return append(result, []byte("\n")...)
}

func (m Message) GetUint64() (uint64, error) {
	if m.MessageType != MessageTypeSolution {
		return 0, errors.New("not a solution message")
	}
	return strconv.ParseUint(m.MessageString, 10, 64)
}

func ParseServerMessage(message string) (Message, error) {
	var result = Message{}
	err := json.Unmarshal([]byte(message), &result)
	if err != nil {
		return Message{}, err
	}
	if !validateMessageType(result.MessageType) {
		return Message{}, errors.New("wrong messageType")
	}
	return result, nil
}

func validateMessageType(messageType string) bool {
	if _, ok := messageTypes[messageType]; !ok {
		return false
	}
	return true
}
