/*
	apitypes.go

	This file contains the types for the Messenger API
*/

package fb

var MessageTypeMessageTag string = "MESSAGE_TAG"
var MessageTypeResponse string = "RESPONSE"
var MessageTypeUpdate string = "UPDATE"

var NotificationTypeRegular string = "REGULAR"
var NotificationTypeSilentPush string = "SILENT_PUSH"
var NotificationTypeNoPush string = "NO_PUSH"

var SenderActionMarkSeen string = "mark_seen"
var SenderActionTypingOn string = "typing_on"
var SenderActionTypingOff string = "typing_off"

var QuickReplyLocation QuickReply = QuickReply{
	ContentType: "location",
}

var QuickReplyPhoneNumber QuickReply = QuickReply{
	ContentType: "user_phone_number",
}

var QuickReplyEmail QuickReply = QuickReply{
	ContentType: "user_email",
}

type User struct {
	ID string `json:"id"`
}

type ThreadOwner struct {
	AppID string `json:"app_id"`
}

type ThreadOwnerResponse struct {
	Data []struct {
		ThreadOwner ThreadOwner `json:"thread_owner"`
	} `json:"data"`
}

type MessageResponse struct {
	RecipientID string `json:"recipient_id"`
	MessageID   string `json:"message_id"`
}

type APIError struct {
	FBTraceID string  `json:"fbtrace_id"`
	Message   string  `json:"message"`
	Type      string  `json:"type"`
	Code      float64 `json:"code"`
	SubCode   float64 `json:"error_subcode"`
}

type OutgoingMessage struct {
	Text             string              `json:"text,omitempty"`
	Metadata         string              `json:"string,omitempty"`
	QuickReplies     []QuickReply        `json:"quick_replies,omitempty"`
	NotificationType string              `json:"notification_type,omitempty"`
	Attachment       *OutgoingAttachment `json:"attachment,omitempty"`
	Tag              string              `json:"tag,omitempty"`
}

type OutgoingShareMessage struct {
	Attachment OutgoingAttachment `json:"attachment"`
}

type OutgoingAttachment struct {
	Type    string          `json:"type,omitempty"`
	Payload TemplatePayload `json:"payload,omitempty"`
}

type TemplatePayload struct {
	URL              string                   `json:"url,omitempty"`
	IsReusable       bool                     `json:"is_reusable,omitempty"`
	AttachmentID     string                   `json:"attachment_id,omitempty"`
	TemplateType     string                   `json:"template_type,omitempty"`
	Elements         []GenericTemplateElement `json:"elements,omitempty"`
	Sharable         bool                     `json:"sharable,omitempty"`
	ImageAspectRatio string                   `json:"image_aspect_ratio,omitempty"`
	Text             string                   `json:"text,omitempty"`
	TopElementStyle  string                   `json:"top_element_style,omitempty"`
	Buttons          []Button                 `json:"buttons,omitempty"`
}

type GenericTemplateElement struct {
	Title         string         `json:"title,omitempty"`
	Subtitle      string         `json:"subtitle,omitempty"`
	ImageURL      string         `json:"image_url,omitempty"`
	DefaultAction *DefaultAction `json:"default_action,omitempty"`
	Buttons       []Button       `json:"buttons,omitempty"`
}

type Button struct {
	Type          string                `json:"type,omitempty"`
	URL           string                `json:"url,omitempty"`
	Title         string                `json:"title,omitempty"`
	Payload       string                `json:"payload,omitempty"`
	ShareContents *OutgoingShareMessage `json:"share_contents,omitempty"`
}

type UploadableAttachment struct {
	Type    string                      `json:"type"`
	Payload UploadableAttachmentPayload `json:"payload"`
}

type UploadableAttachmentPayload struct {
	IsReusable bool   `json:"is_reusable"`
	URL        string `json:"url"`
}

type DefaultAction struct {
	URL                 string `json:"url"`
	Type                string `json:"type"`
	WebviewHeightRatio  string `json:"webview_height_ratio,omitempty"`
	MessengerExtensions bool   `json:"messenger_extensions,omitempty"`
	FallbackURL         string `json:"fallback_url,omitempty"`
	WebviewShareButton  string `json:"webview_share_button,omitempty"`
}

type QuickReply struct {
	ContentType string `json:"content_type"`
	Title       string `json:"title,omitempty"`
	Payload     string `json:"payload,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
}

type MessageSendResponse struct {
	RecipientID string `json:"recipient_id"`
	MessageID   string `json:"message_id"`
}

type MessengerProfile struct {
	AccountLinkingURL string           `json:"account_linking_url,omitempty"`
	GetStarted        GetStarted       `json:"get_started,omitempty"`
	PersistentMenu    []PersistentMenu `json:"persistent_menu,omitempty"`
	Greeting          []Greeting       `json:"greeting,omitempty"`
}

type GetStarted struct {
	Payload string `json:"payload"`
}

type PersistentMenu struct {
	Locale                string     `json:"locale,omitempty"`
	ComposerInputDisabled bool       `json:"composer_input_disabled,omitempty"`
	DisabledSurfaces      []string   `json:"disabled_surfaces,omitempty"`
	CallToActions         []MenuItem `json:"call_to_actions,omitempty"`
}

type MenuItem struct {
	Type               string     `json:"type"`
	Title              string     `json:"title"`
	URL                string     `json:"url,omitempty"`
	Payload            string     `json:"payload,omitempty"`
	CallToActions      []MenuItem `json:"call_to_actions,omitempty"`
	WebviewHeightRatio string     `json:"webview_height_ratio,omitempty"`
	FallbackURL        string     `json:"fallback_url,omitempty"`
	WebviewShareButton string     `json:"webview_share_button,omitempty"`
}

type Greeting struct {
	Locale string `json:"locale"`
	Text   string `json:"text"`
}
