/*
	webhooktypes.go

	This file contains the types for incoming Facebook webhook requests
*/

package fb

type WebhookRequest struct {
	Object string         `json:"object"`
	Entry  []WebhookEntry `json:"entry"`
}

type WebhookEntry struct {
	ID        string          `json:"id"`
	Time      int64           `json:"time"`
	Messaging []MessagingItem `json:"messaging"`
	Standby   []MessagingItem `json:"standby"`
}

type MessagingItem struct {
	IsStandby         bool                `json:"standby"`
	Sender            User                `json:"sender"`
	Recipient         User                `json:"recipient"`
	Timestamp         int64               `json:"timestamp"`
	Message           WHMessage           `json:"message"`
	Postback          WHPostback          `json:"postback"`
	Referral          WHReferral          `json:"referral"`
	TakeThreadControl WHTakeThreadControl `json:"take_thread_control"`
}

type WHMessage struct {
	MID         string         `json:"mid"`
	Text        string         `json:"text"`
	Seq         int64          `json:"seq"`
	IsEcho      bool           `json:"is_echo"`
	AppID       int64          `json:"app_id"`
	Metadata    string         `json:"metadata"`
	Attachments []WHAttachment `json:"attachments"`
	QuickReply  struct {
		Payload string `json:"payload"`
	} `json:"quick_reply"`
}

type WHPostback struct {
	Title    string     `json:"title"`
	Payload  string     `json:"payload"`
	Referral WHReferral `json:"referral"`
}

type WHReferral struct {
	Source     string `json:"source"`
	Type       string `json:"type"`
	Ref        string `json:"ref"`
	AdID       string `json:"ad_id"`
	RefererURI string `json:"referer_uri"`
}

type WHAttachment struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Type    string `json:"type"`
	Payload struct {
		URL         string        `json:"url"`
		Coordinates WHCoordinates `json:"coordinates"`
	} `json:"payload"`
}

type WHTakeThreadControl struct {
	PreviousOwnerAppID int64  `json:"previous_owner_app_id"`
	Metadata           string `json:"metadata"`
}

type WHCoordinates struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}
