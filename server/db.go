package server

import (
	// "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/stdlib"
)

func (s *server) setupDB() {
	c, err := pgx.ParseConfig(s.connectionString())
	if err != nil {
		s.logErr("postgres URI parsing error", err)
		return
	}

	c.Logger = logrusadapter.NewLogger(s.logger)
	db := stdlib.OpenDB(*c)

	// TODO add migration assistant
	// err = validateSchema(db)

	s.db = db
	// TODO add querier and statement builder
}
