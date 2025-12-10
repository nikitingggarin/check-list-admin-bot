package infrastructure

import (
	"context"
	"fmt"

	"github.com/nedpals/supabase-go"
)

type DatabaseClient struct {
	client *supabase.Client
}

func NewDatabaseClient(url, apiKey string) (*DatabaseClient, error) {
	if url == "" || apiKey == "" {
		return nil, fmt.Errorf("database URL and API key are required")
	}

	client := supabase.CreateClient(url, apiKey)
	return &DatabaseClient{client: client}, nil
}

func (d *DatabaseClient) HealthCheck(ctx context.Context) error {
	if d.client == nil {
		return fmt.Errorf("database client is not initialized")
	}
	return nil
}

func (d *DatabaseClient) Client() *supabase.Client {
	return d.client
}
