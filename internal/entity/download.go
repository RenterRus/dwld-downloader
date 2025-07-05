package entity

type Status int

const (
	NEW Status = iota + 1
	WORK
	SENDING
	DONE
)

var StatusMapping map[Status]string = map[Status]string{
	NEW:     "NEW",
	WORK:    "WORK",
	SENDING: "SENDING",
	DONE:    "DONE",
}
