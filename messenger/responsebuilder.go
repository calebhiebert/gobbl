/*
	responsebuilder.go

	This file contains utilities for crafting Messenger response objects
*/

package fb

import (
	"errors"
	"math/rand"
	"time"

	"github.com/calebhiebert/gobbl"
)

// MBResponse holds a facebook messenger response object
type MBResponse struct {
	Messages      []OutgoingMessage
	QuickReplies  []QuickReply
	MinTypingTime []time.Duration
}

// MBImmediateResponse is the same as MBResponse, except it can be sent any time
type MBImmediateResponse struct {
	MBResponse
	Context     *gbl.Context
	Integration *MessengerIntegration
}

// CreateResponse will return a pre-populated messenger response object and add it to the context
func CreateResponse(c *gbl.Context) *MBResponse {
	existingResponse := c.R.(*MBResponse)

	if existingResponse.Messages == nil {
		existingResponse.Messages = []OutgoingMessage{}
	}

	if existingResponse.QuickReplies == nil {
		existingResponse.QuickReplies = []QuickReply{}
	}

	if existingResponse.MinTypingTime == nil {
		existingResponse.MinTypingTime = []time.Duration{time.Second}
	}

	return existingResponse
}

// CreateImmediateResponse will create a response object that can be sent manually at any time
func CreateImmediateResponse(c *gbl.Context) *MBImmediateResponse {
	integration := c.Integration.(*MessengerIntegration)

	immediateResponse := MBImmediateResponse{
		Integration: integration,
		Context:     c,
		MBResponse: MBResponse{
			Messages:      []OutgoingMessage{},
			QuickReplies:  []QuickReply{},
			MinTypingTime: []time.Duration{time.Second},
		},
	}

	return &immediateResponse
}

// Send will immediately send all messages in the response
func (im *MBImmediateResponse) Send() error {
	if im.Context.User.ID == "" {
		return errors.New("Missing user ID!")
	}

	return im.Integration.doResponse(im.Context.User.ID, &im.MBResponse)
}

// M adds a new message to the response and returns it
func (m *MBResponse) M(om *OutgoingMessage) {
	m.Messages = append(m.Messages, *om)
}

// Template will add a message with the given template payload
func (m *MBResponse) Template(template *TemplatePayload) {
	m.M(&OutgoingMessage{
		Attachment: &OutgoingAttachment{
			Type:    "template",
			Payload: *template,
		},
	})
}

// Text will add a single text message to the output
func (m *MBResponse) Text(text string) {
	m.Messages = append(m.Messages, OutgoingMessage{
		Text: text,
	})
}

// RandomText will choose one string from the provided list and send it
func (m *MBResponse) RandomText(text ...string) {
	m.Messages = append(m.Messages, OutgoingMessage{
		Text: text[rand.Intn(len(text))],
	})
}

func (m *MBResponse) TypingTime(mtt ...time.Duration) {
	m.MinTypingTime = mtt
}

// Image will send an image with the following url to the chat
func (m *MBResponse) Image(url string) {
	m.Messages = append(m.Messages, OutgoingMessage{
		Attachment: &OutgoingAttachment{
			Type: "image",
			Payload: TemplatePayload{
				URL:        url,
				IsReusable: true,
			},
		},
	})
}

// QR will add one or more quick replies to the repsonse
func (m *MBResponse) QR(qr ...QuickReply) {
	for _, quickReply := range qr {
		m.QuickReplies = append(m.QuickReplies, quickReply)
	}
}

// QRText is a helper function to create a text quickreply
func QRText(title string, payload string) QuickReply {
	return QuickReply{
		ContentType: "text",
		Title:       title,
		Payload:     payload,
	}
}

// ImageCardElement returns a facebook GenericTemplateElement
func ImageCardElement(title, subtitle, imageURL string) GenericTemplateElement {
	return GenericTemplateElement{
		Title:    title,
		Subtitle: subtitle,
		ImageURL: imageURL,
	}
}

// ImageCardElementClickable returns an image card element with the default action set
func ImageCardElementClickable(title, subtitle, imageURL, actionURL string) GenericTemplateElement {
	gt := ImageCardElement(title, subtitle, imageURL)

	gt.DefaultAction = &DefaultAction{
		Type:               "web_url",
		URL:                actionURL,
		WebviewHeightRatio: "tall",
	}

	return gt
}

// Button will add one or more buttons to the template element
func (ge *GenericTemplateElement) Button(buttons ...Button) *GenericTemplateElement {
	if ge.Buttons == nil {
		ge.Buttons = []Button{}
	}

	ge.Buttons = append(ge.Buttons, buttons...)

	return ge
}

// ImageCard returns a ready-to-go template elment
func ImageCard(title, subtitle, imageURL string) TemplatePayload {
	return TemplatePayload{
		TemplateType: "generic",
		Elements:     []GenericTemplateElement{ImageCardElement(title, subtitle, imageURL)},
	}
}

// Carousel will combine multiple image card elements into a carousel
func Carousel(elements ...GenericTemplateElement) TemplatePayload {
	return TemplatePayload{
		TemplateType: "generic",
		Elements:     elements,
	}
}

// Element will add one or more elements to a generic template carousel
// Do not use this mehtod on a non generic template carousel
func (c *TemplatePayload) Element(elements ...GenericTemplateElement) *TemplatePayload {
	if c.Elements != nil {
		c.Elements = append(c.Elements, elements...)
	} else {
		c.Elements = elements
	}

	return c
}

// ButtonURL creates a facebook URL button
func ButtonURL(title, url string) Button {
	return Button{
		Type:  "web_url",
		Title: title,
		URL:   url,
	}
}

// ButtonPostback creates a facbook postback button
func ButtonPostback(title, payload string) Button {
	return Button{
		Type:    "postback",
		Title:   title,
		Payload: payload,
	}
}

// ButtonShare creates a facebook share button
func ButtonShare() Button {
	return Button{
		Type: "element_share",
	}
}
