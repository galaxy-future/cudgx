package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/galaxy-future/cudgx/internal/predict/consts"
)

func NewBridgxClient(serverAddress string) *Client {
	return &Client{
		ServerAddress: serverAddress,
		HttpClient: &http.Client{
			Timeout: 5000 * time.Millisecond,
		},
	}
}

func authXClient() (token string, err error) {
	request := struct {
		Username string
		Password string
	}{Username: consts.XClientUsername, Password: consts.XClientPassword}
	data, _ := json.Marshal(&request)
	resp, err := bridgxClient.HttpClient.Post(fmt.Sprintf("%s/user/login", bridgxClient.ServerAddress), "application/json", bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	response := struct {
		Code int
		Msg  string
		Data string
	}{}
	err = json.Unmarshal(respData, &response)
	if err != nil {
		return "", err
	}
	if response.Code != 200 {
		return "", fmt.Errorf("Wrong code: %v, response : %s", err, respData)
	}
	return response.Data, nil
}
