package repository

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	//	"fmt"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Repo struct {
	Database *sql.DB
}

func NewRepo(dsn string) (*Repo, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errors.WithMessage(err, "Repository: opening Error DB")
	}
	logrus.Debug("Repository: success opening DB")

	if err := db.Ping(); err != nil {
		return nil, errors.WithMessage(err, "Repository: pinging DB")
	}
	logrus.Debug("Repository: ping was successful, New Repository created")
	return &Repo{Database: db}, nil
}

func (r *Repo) DoRequest(ctx context.Context, textRequest string) error {
	_, err := r.Database.ExecContext(ctx, textRequest)
	if err != nil {
		return errors.WithMessage(err, "Repository: Request error")
	}

	return nil
}
