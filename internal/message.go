package internal

type Message struct {
	FromName string `json:"from_id"`
	ToName   string `json:"to_id"`
	Data     []byte `json:"data"`
}
