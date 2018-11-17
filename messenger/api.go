/*
	api.go

	This file contains a small subset of the facebook graph api, specifically the api routes used for building chatbots
*/

package fb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type MessengerAPI struct {
	accessToken string
	http        *http.Client
	baseURL     string
}

// CreateMessengerAPI will create a functional messenger api that is setup to use the provided access token.
// This will create an internal http client with a timeout set to 6 seconds
func CreateMessengerAPI(accessToken string) *MessengerAPI {
	mapi := MessengerAPI{
		accessToken: strings.TrimSpace(accessToken),
		baseURL:     "https://graph.facebook.com/v2.6",
	}

	mapi.http = &http.Client{
		Timeout: 6 * time.Second,
	}

	return &mapi
}

// SendMessage will send a messenger message.
// An error will be returned if Facebook returns a non 2xx status code.
// Official documentation can be found here: https://developers.facebook.com/docs/messenger-platform/reference/send-api/
func (m *MessengerAPI) SendMessage(recipient *User, messageType string, message *OutgoingMessage) (interface{}, error) {

	body := map[string]interface{}{
		"messaging_type": messageType,
		"recipient":      recipient,
		"message":        message,
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/me/messages?access_token=%s", m.baseURL, m.accessToken)

	resp, err := m.http.Post(url, "application/json", bytes.NewReader(jsonBytes))
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
		return nil, errors.New(fmt.Sprintf("%+v", responseBody))
	}
}

// SenderAction will set the sender action (typing, read, not typing) for a given psid
// Official documentation can be found here: https://developers.facebook.com/docs/messenger-platform/reference/send-api/
func (m *MessengerAPI) SenderAction(recipient *User, action string) (interface{}, error) {
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
		return nil, errors.New(fmt.Sprintf("%+v", responseBody))
	}
}

// UserInfo will load information on a given psid
func (m *MessengerAPI) UserInfo(psid string) (interface{}, error) {
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

		return nil, mapError(body)
	}
}

func (m *MessengerAPI) MessengerProfile(profile *MessengerProfile) (interface{}, error) {
	jsonBytes, err := json.Marshal(profile)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/me/messenger_profile?access_token=%s", m.baseURL, m.accessToken)

	resp, err := m.http.Post(url, "application/json", bytes.NewReader(jsonBytes))
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
		return nil, errors.New(fmt.Sprintf("%+v", responseBody))
	}
}

// Reads a http response body and parses the json into an interface
func readBodyJson(body io.ReadCloser) (interface{}, error) {
	defer body.Close()

	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var bodyJson interface{}

	err = json.Unmarshal(bytes, &bodyJson)
	if err != nil {
		return string(bytes), nil
	}

	return bodyJson, nil
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
