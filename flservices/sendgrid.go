package flservices

import (
	"encoding/json"
	"fmt"

	"github.com/nerde/fuji-lane-back/flconfig"
	"github.com/sendgrid/rest"
	sendgrid "github.com/sendgrid/sendgrid-go"
)

// https://sendgrid.com/docs/api-reference/

type sendgridResponse struct {
	PersistedRecipients []string `json:"persisted_recipients"`
}

// SendgridContact represents a contact we can send to Sendgrid
type SendgridContact struct {
	Email string `json:"email"`
}

// Sendgrid service interface
type Sendgrid interface {
	SubscribeNewsletter(SendgridContact) error
}

// SendgridAPI is a Sendgrid implementation that uses their API to perform actions
type SendgridAPI struct {
	key              string
	newsletterListID string
}

// SubscribeNewsletter creates a contact and subscribes it to our newsletter
func (s *SendgridAPI) SubscribeNewsletter(c SendgridContact) error {
	resp, err := s.createContacts([]SendgridContact{c})
	if err != nil {
		return err
	}

	return s.subscribeNewsletter(resp.PersistedRecipients)
}

func (s *SendgridAPI) createContacts(cs []SendgridContact) (resp *sendgridResponse, err error) {
	request := s.getRequest("POST", "/v3/contactdb/recipients")
	request.Body, err = json.Marshal(cs)
	if err != nil {
		return
	}

	var response *rest.Response
	response, err = sendgrid.API(request)
	if err != nil {
		return
	}

	resp = &sendgridResponse{PersistedRecipients: []string{}}
	err = json.Unmarshal([]byte(response.Body), resp)

	return
}

func (s *SendgridAPI) subscribeNewsletter(contactIDs []string) (err error) {
	request := s.getRequest("POST", fmt.Sprintf("/v3/contactdb/lists/%s/recipients", s.newsletterListID))
	request.Body, err = json.Marshal(contactIDs)
	if err != nil {
		return
	}

	_, err = sendgrid.API(request)

	return
}

func (s *SendgridAPI) getRequest(method rest.Method, path string) rest.Request {
	request := sendgrid.GetRequest(s.key, path, "")
	request.Method = method

	return request
}

// NewSendgridAPI creates a SendgridAPI with system configuration
func NewSendgridAPI() *SendgridAPI {
	return &SendgridAPI{key: flconfig.Config.SendgridKey, newsletterListID: flconfig.Config.SendgridNewsletterListID}
}
