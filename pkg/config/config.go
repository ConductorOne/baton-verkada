package config

import (
	"fmt"

	"github.com/conductorone/baton-sdk/pkg/field"
)

var ApiKeyField = field.StringField(
	"api-key",
	field.WithDisplayName("Api key"),
	field.WithDescription("Verkada API key"),
	field.WithRequired(true),
	field.WithIsSecret(true),
)
var RegionField = field.SelectField(
	"region",
	[]string{"US", "EU"},
	field.WithDisplayName("Region"),
	field.WithDescription("API region. Default is US. In case of EU based organization, pass region as EU."),
	field.WithDefaultValue("US"),
)

var BaseURLField = field.StringField(
	"base-url",
	field.WithDescription("Override the Verkada API URL (for testing)"),
	field.WithHidden(true),
	field.WithExportTarget(field.ExportTargetCLIOnly),
)

//go:generate go run ./gen
var Config = field.NewConfiguration(
	[]field.SchemaField{
		ApiKeyField,
		RegionField,
		BaseURLField,
	},
)

// ValidateConfig is run after the configuration is loaded, and should return an
// error if it isn't valid. Implementing this function is optional, it only
// needs to perform extra validations that cannot be encoded with configuration
// parameters.
func ValidateConfig(cfg *Verkada) error {
	if cfg.ApiKey == "" {
		return fmt.Errorf("config: missing api key")
	}

	return nil
}
