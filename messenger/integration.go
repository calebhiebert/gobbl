package fb

import (
	".."
)

type MessengerIntegration struct {
	API *MessengerAPI
}

func (m *MessengerIntegration) GenericRequest(c *cpn.Context) (*cpn.GenericRequest, error) {
	return &cpn.GenericRequest{
		Text: "Howdy",
	}, nil
}

func (m *MessengerIntegration) User(c *cpn.Context) (*cpn.User, error) {
	return &cpn.User{
		ID: "123456789",
	}, nil
}

func (m *MessengerIntegration) Respond(c *cpn.Context) (*interface{}, error) {
	result, err := m.API.SenderAction(&User{
		ID: c.User.ID,
	}, SenderActionTypingOn)
	if err != nil {
		return nil, err
	}

	return result, nil
}
