package sdmht_conn

const (
	ClientEventTypeAdd    = 1
	ClientEventTypeRemove = 2
)

type ClientEvent struct {
	Type      int
	AccountID uint64
	Client    *Client
}

func NewClientEvent(eventType int, accountID uint64, c *Client) ClientEvent {
	return ClientEvent{
		Type:      eventType,
		AccountID: accountID,
		Client:    c,
	}
}
