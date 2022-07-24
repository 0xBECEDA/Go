package internal

type Message struct {
	FromName string `json:"from_id"`
	FromHost string `json:"from_host"`
	ToName   string `json:"to_id"`
	Data     []byte `json:"data"`
}

type AuthorizeMessage struct {
	Name  string
	Email string
	Host  string
}
