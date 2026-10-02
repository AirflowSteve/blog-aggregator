package main

import (
	"github.com/AirflowSteve/blog-aggregator/internal/config"
	"github.com/AirflowSteve/blog-aggregator/internal/database"
)

type State struct {
	db     *database.Queries
	config *config.Config
}
