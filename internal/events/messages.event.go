package events

type Message interface {
	Type() string
}

type CreatedFeedMessage struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

func (message CreatedFeedMessage) Type() string {
	return "created_feed"
}
