package db

import (
	"context"
	"log"

	"google.golang.org/api/option"

	"cloud.google.com/go/spanner"
)

func NewClient(ctx context.Context, dsn string) (*spanner.Client, error) {
	client, err := spanner.NewClient(ctx, dsn, option.WithoutAuthentication())
	if err != nil {
		log.Fatalf("failed to create spanner client: %v", err)
	}

	return client, nil
}
