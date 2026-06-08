package subscription

import (
	"database/sql"

	def "github.com/LearLocker/streaming/subscription/internal/repository"
)

var _ def.SubscriptionRepository = (*repository)(nil)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db: db}
}
