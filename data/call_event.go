package data

import (
	"log"
	"strings"
	"time"
)

const (
	STRING_DELIMITER = ";"
	TIME_FORMAT      = "02.01.06 15:04:05"
)

type CallEvent struct {
	Date         time.Time
	Operation    string `sql:not null`
	Irrelevant   string
	Line         string `sql:not null`
	SourceNumber string
	TargetNumber string
	InternalPort string
}

func NewEvent(data string) *CallEvent {
	event := &CallEvent{}
	if strings.Contains(data, STRING_DELIMITER) {
		arr := strings.Split(data, STRING_DELIMITER)
		if len(arr) > 3 {
			event.SetFields(arr)
		} else {
			event = nil
		}
	} else {
		event = nil
	}
	log.Println(event)
	return event
}

// 02.01.06 15:04:05;CALL;0;4;034297908900;03429749890;SIP3;
// 02.01.06 15:04:05;RING;0;034297908900;03429749890;SIP3;
// 02.01.06 15:04:05;CONNECT;0;4;03429749890>;
// 02.01.06 15:04:05;DISCONNECT;0;4;

func (e *CallEvent) SetFields(array []string) *CallEvent {

	timeObj, err := time.Parse(TIME_FORMAT, array[0])
	if err != nil {
		timeObj = time.Now()
	}
	e.Date = timeObj

	if len(array) > 2 {
		e.Operation = array[1]
	}

	if len(array) > 3 {
		e.Irrelevant = array[2]
	}

	if len(array) > 4 {
		e.Line = array[3]
	}

	if len(array) > 5 {
		e.SourceNumber = array[4]
	}

	if len(array) > 6 {
		e.TargetNumber = array[5]
	}

	if len(array) > 7 {
		e.InternalPort = array[6]
	}
	return e
}

func (e *CallEvent) IsValid() bool {
	return e.Operation != ""
}
