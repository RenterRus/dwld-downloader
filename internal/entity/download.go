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
