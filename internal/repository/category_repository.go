package repository

import "github.com/thenopholo/my-money/internal/repository/postgres"

type categoryRepository struct {
	queries *postgres.Queries
}

func NewCategoryRepository(q *postgres.Queries) *categoryRepository {
	return &categoryRepository{
		queries: q,
	}
}

