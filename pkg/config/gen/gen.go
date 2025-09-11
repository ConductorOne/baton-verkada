package main

import (
	"github.com/conductorone/baton-sdk/pkg/config"
	cfg "github.com/conductorone/baton-verkada/pkg/config"
)

func main() {
	config.Generate("verkada", cfg.Config)
}
