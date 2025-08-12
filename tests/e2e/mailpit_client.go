package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type MailpitClient struct {
	baseURL string
	client  *http.Client
}

type MailpitEmail struct {
	ID      string           `json:"ID"`
	From    MailpitAddress   `json:"From"`
	To      []MailpitAddress `json:"To"`
	Subject string           `json:"Subject"`
	Text    string           `json:"Text"`
	HTML    string           `json:"HTML"`
	Created time.Time        `json:"Created"`
}

type MailpitResponse struct {
	Total          int            `json:"total"`
	Unread         int            `json:"unread"`
	Count          int            `json:"count"`
	MessagesCount  int            `json:"messages_count"`
	MessagesUnread int            `json:"messages_unread"`
	Start          int            `json:"start"`
	Tags           []string       `json:"tags"`
	Messages       []MailpitEmail `json:"messages"`
}

type MailpitAddress struct {
	Name    string `json:"Name"`
	Address string `json:"Address"`
}

func NewMailpitClient(host string, port int) *MailpitClient {
	return &MailpitClient{
		baseURL: fmt.Sprintf("http://%s:%d", host, port),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *MailpitClient) GetMessages() ([]MailpitEmail, error) {
	resp, err := c.client.Get(fmt.Sprintf("%s/api/v1/messages", c.baseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response MailpitResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Messages, nil
}

func (c *MailpitClient) GetMessage(id string) (*MailpitEmail, error) {
	resp, err := c.client.Get(fmt.Sprintf("%s/api/v1/message/%s", c.baseURL, id))
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var message MailpitEmail
	if err := json.Unmarshal(body, &message); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return &message, nil
}

func (c *MailpitClient) DeleteAllMessages() error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/messages", c.baseURL), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete messages: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *MailpitClient) WaitForMessages(expectedCount int, timeout time.Duration) ([]MailpitEmail, error) {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		messages, err := c.GetMessages()
		if err != nil {
			return nil, err
		}

		if len(messages) >= expectedCount {
			return messages, nil
		}

		time.Sleep(500 * time.Millisecond)
	}

	return nil, fmt.Errorf("timeout waiting for %d messages", expectedCount)
}
