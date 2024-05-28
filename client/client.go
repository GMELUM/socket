package client

type IClient interface {
	HandlerEvents(dataType string, data interface{})
}

// Clan struct
type Room struct {
	ID string
}

func (r *Room) UpdateBg() {
	
}

type RoomList map[string]*Room

type Client struct {
	FirstName string
	LastName  string
	Room    *Room
}

func New() *Client {
	return &Client{}
}

func (cl *Client) HandlerEvents(dataType string, data interface{}) {
	switch dataType {
	case "get_clan":
		cl.Send(GET_CLAN(data))
	case "get_user":
		cl.Send(GET_USER(data))
	}
}

func GET_CLAN(data interface{}) *struct{} {
	// request mysql
	return &struct{}{}
}

func GET_USER(data interface{}) *struct{} {
	// request mysql
	return &struct{}{}
}

func (cl *Client) Send(data interface{}) {
	// sending data
}
