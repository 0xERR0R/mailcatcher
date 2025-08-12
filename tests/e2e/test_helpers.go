package e2e

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gopkg.in/gomail.v2"
)

// WaitForContainerHealth waits for a container to be healthy
func WaitForContainerHealth(container testcontainers.Container, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for container health")
		default:
			state, err := container.State(ctx)
			if err != nil {
				log.Printf("Error getting container state: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			if state.Running {
				return nil
			}

			time.Sleep(1 * time.Second)
		}
	}
}

// CreateMailpitContainer creates and starts a mailpit container
func CreateMailpitContainer(ctx context.Context, image string) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        image,
		ExposedPorts: []string{"1025/tcp", "8025/tcp"},
		WaitingFor:   wait.ForLog("accessible via http://localhost:8025/"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	return container, err
}

// CreateMailcatcherContainer creates and starts a mailcatcher container
func CreateMailcatcherContainer(ctx context.Context, dockerfilePath string, env map[string]string) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    dockerfilePath,
			Dockerfile: "Dockerfile",
		},
		ExposedPorts: []string{"1025/tcp"},
		Env:          env,
		WaitingFor:   wait.ForLog("Starting server at"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	return container, err
}

// CleanupContainer safely terminates a container
func CleanupContainer(container testcontainers.Container) {
	if container != nil {
		ctx := context.Background()
		if err := container.Terminate(ctx); err != nil {
			log.Printf("Error terminating container: %v", err)
		}
	}
}

// RetryWithTimeout retries a function until it succeeds or timeout is reached
func RetryWithTimeout(operation func() error, timeout time.Duration, interval time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if err := operation(); err == nil {
			return nil
		}
		time.Sleep(interval)
	}

	return fmt.Errorf("operation failed after %v", timeout)
}

// sendTestEmail sends a test email to the specified host and port
func sendTestEmail(host string, port int, to string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "sender@example.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Test Email")
	m.SetBody("text/plain", "This is a test email body")

	d := gomail.NewDialer(host, port, "", "")
	return d.DialAndSend(m)
}
