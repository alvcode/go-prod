package storage

import (
	"prod/pkg/client/postgresql"
	"prod/pkg/logging"
)

type storage struct {
	client postgresql.Client
	logger *logging.Logger
}
