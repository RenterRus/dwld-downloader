package entity

type Status int

const (
	NEW Status = iota + 1
	WORK
	TO_SEND
	SENDING
	DONE
)

var StatusMapping map[Status]string = map[Status]string{
	NEW:     "NEW",
	WORK:    "WORK",
	TO_SEND: "TO_SEND",
	SENDING: "SENDING",
	DONE:    "DONE",
}

type Stage struct {
	Positions         int
	AttemptBeforeNext int
	Threads           int
	IsCookie          bool
	IsEmbededCharters bool
	IsFormat          bool
	IsMarkWatched     bool
	Extractors        string
}
