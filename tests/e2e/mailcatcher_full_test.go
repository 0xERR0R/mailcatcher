package e2e

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestMailcatcherFullWorkflow tests the complete mailcatcher workflow
// This test requires the mailcatcher:test image to be built first
func TestMailcatcherFullWorkflow(t *testing.T) {
	ctx := context.Background()

	// Create a network for the containers to communicate
	networkName := "mailcatcher-test-network"
	network, err := testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{
			Name: networkName,
		},
	})
	require.NoError(t, err)
	defer func() {
		if err := network.Remove(ctx); err != nil {
			t.Logf("Error removing network: %v", err)
		}
	}()

	// Start Mailpit container with SMTP authentication enabled
	mailpitReq := testcontainers.ContainerRequest{
		Image:        "axllent/mailpit:latest",
		ExposedPorts: []string{"1025/tcp", "8025/tcp"},
		Networks:     []string{networkName},
		NetworkAliases: map[string][]string{
			networkName: {"mailpit"},
		},
		Env: map[string]string{
			"MP_SMTP_AUTH_ACCEPT_ANY":     "1",
			"MP_SMTP_AUTH_ALLOW_INSECURE": "1",
		},
		WaitingFor: wait.ForLog("accessible via http://localhost:8025/"),
	}

	mailpitContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: mailpitReq,
		Started:          true,
	})
	require.NoError(t, err)
	defer func() {
		if err := mailpitContainer.Terminate(ctx); err != nil {
			t.Logf("Error terminating mailpit container: %v", err)
		}
	}()

	// Get Mailpit web port
	mailpitWebPort, err := mailpitContainer.MappedPort(ctx, "8025")
	require.NoError(t, err)

	mailpitHost, err := mailpitContainer.Host(ctx)
	require.NoError(t, err)

	// Start Mailcatcher container using pre-built image
	mailcatcherReq := testcontainers.ContainerRequest{
		Image:        "mailcatcher:test",
		ExposedPorts: []string{"1025/tcp"},
		Networks:     []string{networkName},
		Env: map[string]string{
			"MC_PORT":          "1025",
			"MC_HOST":          "test.example.com",
			"MC_REDIRECT_TO":   "test@example.com",
			"MC_SENDER_MAIL":   "mailcatcher@test.example.com",
			"MC_SMTP_HOST":     "mailpit",
			"MC_SMTP_PORT":     "1025",
			"MC_SMTP_USER":     "",
			"MC_SMTP_PASSWORD": "",
		},
		WaitingFor: wait.ForLog("Starting server at"),
	}

	mailcatcherContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: mailcatcherReq,
		Started:          true,
	})
	require.NoError(t, err)
	defer func() {
		if err := mailcatcherContainer.Terminate(ctx); err != nil {
			t.Logf("Error terminating mailcatcher container: %v", err)
		}
	}()

	// Get Mailcatcher port
	mailcatcherPort, err := mailcatcherContainer.MappedPort(ctx, "1025")
	require.NoError(t, err)

	mailcatcherHost, err := mailcatcherContainer.Host(ctx)
	require.NoError(t, err)

	// Create mailpit client
	mailpitClient := NewMailpitClient(mailpitHost, mailpitWebPort.Int())

	// Clear any existing messages
	err = mailpitClient.DeleteAllMessages()
	require.NoError(t, err)

	t.Run("Forward email through mailcatcher", func(t *testing.T) {
		// Send email to mailcatcher
		err := sendTestEmail(mailcatcherHost, mailcatcherPort.Int(), "test@test.example.com")
		require.NoError(t, err)

		// Wait a bit and check mailcatcher logs
		time.Sleep(2 * time.Second)
		logs, err := mailcatcherContainer.Logs(ctx)
		if err == nil {
			logContent, _ := io.ReadAll(logs)
			t.Logf("Mailcatcher logs: %s", string(logContent))
		}

		// Wait for email to be processed and forwarded
		messages, err := mailpitClient.WaitForMessages(1, 10*time.Second)
		require.NoError(t, err)
		assert.Len(t, messages, 1, "Expected 1 email in mailpit")

		if len(messages) > 0 {
			message := messages[0]
			assert.Contains(t, message.Subject, "[MAILCATCHER]", "Email subject should contain [MAILCATCHER] prefix")
			assert.Equal(t, "mailcatcher@test.example.com", message.From.Address, "Email should be from mailcatcher")
			assert.Len(t, message.To, 1, "Email should have one recipient")
			if len(message.To) > 0 {
				assert.Equal(t, "test@example.com", message.To[0].Address, "Email should be forwarded to the configured address")
			}
		}
	})

	t.Run("Ignore invalid domain", func(t *testing.T) {
		// Clear previous messages
		err := mailpitClient.DeleteAllMessages()
		require.NoError(t, err)

		// Send email with invalid domain
		err = sendTestEmail(mailcatcherHost, mailcatcherPort.Int(), "test@invalid-domain.com")
		require.NoError(t, err)

		// Wait a bit and check that no email was forwarded
		time.Sleep(3 * time.Second)

		messages, err := mailpitClient.GetMessages()
		require.NoError(t, err)
		assert.Len(t, messages, 0, "Expected 0 emails in mailpit for invalid domain")
	})
}
