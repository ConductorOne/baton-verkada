package connector

import (
	"context"
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/conductorone/baton-verkada/pkg/verkada"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

type Connector struct {
	client *verkada.Client
}

// ResourceSyncers returns a ResourceSyncer for each resource type that should be synced from the upstream service.
func (v *Connector) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		newUserBuilder(v.client),
		newGroupBuilder(v.client),
	}
}

// Metadata returns metadata about the connector.
func (v *Connector) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "Verkada connector",
		Description: "Connector syncing users and groups from Verkada to Baton.",
	}, nil
}

// Validate is called to ensure that the connector is properly configured. It should exercise any API credentials
// to be sure that they are valid.
func (v *Connector) Validate(ctx context.Context) (annotations.Annotations, error) {
	_, err := v.client.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to validate API credentials: %w", err)
	}
	return nil, nil
}

// New returns a new instance of the connector.
func New(ctx context.Context, apiKey, region string) (*Connector, error) {
	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	if err != nil {
		return nil, err
	}

	return &Connector{
		client: verkada.NewClient(httpClient, apiKey, region),
	}, nil
}
