package fb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type MessengerAPI struct {
	accessToken string
	http        *http.Client
	baseURL     string
}

func CreateMessengerAPI(accessToken string) *MessengerAPI {
	mapi := MessengerAPI{
		accessToken: accessToken,
		baseURL:     "https://graph.facebook.com/v2.6",
	}

	mapi.http = &http.Client{
		Timeout: 6 * time.Second,
	}

	return &mapi
}

func (m *MessengerAPI) SendMessage(recipient *User, messageType string, message *OutgoingMessage) (*interface{}, error) {

	body := map[string]interface{}{
		"messaging_type": messageType,
		"recipient":      recipient,
		"message":        message,
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := m.http.Post(fmt.Sprintf("%s/me/messages?access_token=%s", m.baseURL, m.accessToken), "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}

	return readBodyJson(resp.Body)
}

func (m *MessengerAPI) SenderAction(recipient *User, action string) (*interface{}, error) {
	body := map[string]interface{}{
		"recipient":     recipient,
		"sender_action": action,
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := m.http.Post(fmt.Sprintf("%s/me/messages?access_token=%s", m.baseURL, m.accessToken), "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}

	return readBodyJson(resp.Body)
}

func (m *MessengerAPI) UserInfo(psid string) (*interface{}, error) {
	resp, err := m.http.Get(fmt.Sprintf("https://graph.facebook.com/%s?access_token=%s", psid, m.accessToken))
	if err != nil {
		return nil, err
	}

	return readBodyJson(resp.Body)
}

// Reads a http response body and parses the json into an interface
func readBodyJson(body io.ReadCloser) (*interface{}, error) {
	defer body.Close()

	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var bodyJson interface{}

	err = json.Unmarshal(bytes, &bodyJson)
	if err != nil {
		return nil, err
	}

	return &bodyJson, nil
}
