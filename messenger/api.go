package fb

import (
	"bytes"
	"encoding/json"
	"errors"
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

	responseBody, err := readBodyJson(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return responseBody, nil
	} else {
		return nil, errors.New(fmt.Sprintf("%+v", *responseBody))
	}
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

	responseBody, err := readBodyJson(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return responseBody, nil
	} else {
		return nil, errors.New(fmt.Sprintf("%+v", *responseBody))
	}
}

func (m *MessengerAPI) UserInfo(psid string) (*interface{}, error) {
	resp, err := m.http.Get(fmt.Sprintf("https://graph.facebook.com/%s?access_token=%s", psid, m.accessToken))
	if err != nil {
		return nil, err
	}

	body, err := readBodyJson(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return body, nil
	} else {

		return nil, mapError(*body)
	}
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

func mapError(err interface{}) APIError {
	apiErr := APIError{}

	errMap := err.(map[string]interface{})
	errObj := errMap["error"].(map[string]interface{})

	apiErr.Code = errObj["code"].(float64)
	apiErr.SubCode = errObj["error_subcode"].(float64)
	apiErr.Message = errObj["message"].(string)
	apiErr.FBTraceID = errObj["fbtrace_id"].(string)
	apiErr.Type = errObj["type"].(string)

	return apiErr
}

func (e APIError) Error() string {
	return e.Message
}
