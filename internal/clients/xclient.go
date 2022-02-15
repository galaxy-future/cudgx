package clients

import "net/http"

var bridgxClient, schedulxClient *Client

type Client struct {
	ServerAddress string
	HttpClient    *http.Client
}

func InitializeBridgxClient(bridgxServerAddress string) {
	bridgxClient = NewBridgxClient(bridgxServerAddress)
}

func InitializeSchedulxClient(schedulxServerAddress string) {
	schedulxClient = NewSchedulxClient(schedulxServerAddress)
}
