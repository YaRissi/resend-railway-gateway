package resend

import (
	"github.com/igorrius/resend-railway-gateway/internal/domain"
	resendgo "github.com/resend/resend-go/v2"
)

// Client implements the Resend API email sender adapter.
// It wraps the Resend Go SDK and adapts it to the domain OutboundEmailSender interface.
type Client struct {
	client *resendgo.Client
	logger domain.MessageLogger
}

// NewClient creates a new Resend client with the given API key.
func NewClient(apiKey string, logger domain.MessageLogger) *Client {
	return &Client{client: resendgo.NewClient(apiKey), logger: logger}
}

// Send converts the domain Email to Resend's format and sends it via the API.
func (c *Client) Send(email domain.Email) error {
	attachmentMeta := make([]map[string]any, 0, len(email.Attachments))
	for _, a := range email.Attachments {
		attachmentMeta = append(attachmentMeta, map[string]any{
			"filename":     a.Filename,
			"content_type": a.ContentType,
			"size":         len(a.Content),
			"content_id":   a.ContentID,
		})
	}
	c.logger.Debug("resend_request_start", map[string]any{
		"subject":     email.Subject,
		"to":          email.To,
		"text_body":   email.Text,
		"html_body":   email.HTML,
		"attachments": attachmentMeta,
	})

	attachments := make([]*resendgo.Attachment, 0, len(email.Attachments))
	for _, a := range email.Attachments {
		attachments = append(attachments, &resendgo.Attachment{
			Filename:    a.Filename,
			Content:     a.Content,
			ContentType: a.ContentType,
			ContentId:   a.ContentID,
		})
	}
	tags := make([]resendgo.Tag, 0, len(email.Tags))
	for _, t := range email.Tags {
		tags = append(tags, resendgo.Tag{Name: t.Name, Value: t.Value})
	}
	request := &resendgo.SendEmailRequest{
		From:        email.From,
		To:          email.To,
		Cc:          email.Cc,
		Bcc:         email.Bcc,
		ReplyTo:     email.ReplyTo,
		Subject:     email.Subject,
		Html:        email.HTML,
		Text:        email.Text,
		Attachments: attachments,
		Tags:        tags,
		Headers:     email.Headers,
	}
	resp, err := c.client.Emails.Send(request)
	if err == nil {
		c.logger.Debug("resend_request_success", map[string]any{
			"id": resp.Id,
		})
	}
	return err
}

var _ domain.OutboundEmailSender = (*Client)(nil)
