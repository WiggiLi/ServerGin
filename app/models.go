package app

type Event struct {
	Page      string `json:"Page"`
	Count     string `json:"Count"`
	Name      string `json:"Name"`
	Post      string `json:"Post"`
	DateStart string `json:"DateStart"`
	DateEnd   string `json:"DateEnd"`
}

// NewEvent constructs a event object
func NewEvent() *Event {
	return &Event{}
}

type AllEvents []Event

func GetEvents() *AllEvents {
	return &AllEvents{}
}
