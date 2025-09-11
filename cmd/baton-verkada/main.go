package main

import (
	"context"
	"fmt"
	"os"

	cfg "github.com/conductorone/baton-verkada/pkg/config"

	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/types"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"

	configschema "github.com/conductorone/baton-sdk/pkg/config"

	"github.com/conductorone/baton-verkada/pkg/connector"
)

var version = "dev"
var connectorName = "baton-verkada"

func main() {
	ctx := context.Background()
	_, cmd, err := configschema.DefineConfiguration(
		ctx,
		connectorName,
		getConnector,
		cfg.Config,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cmd.Version = version
	err = cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getConnector(ctx context.Context, cc *cfg.Verkada) (types.ConnectorServer, error) {
	l := ctxzap.Extract(ctx)

	if err := cfg.ValidateConfig(cc); err != nil {
		return nil, err
	}

	cb, err := connector.New(
		ctx,
		cc.ApiKey,
		cc.Region,
	)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	c, err := connectorbuilder.NewConnector(ctx, cb)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	return c, nil
}
