package client

type Client struct {
	Id     string
	Secret string
}

func NewClient(clientId, secret string) Client {
	return Client{
		Id:     clientId,
		Secret: secret,
	}
}
