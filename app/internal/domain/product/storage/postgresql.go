package storage

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"prod/internal/domain/product/model"
	"prod/pkg/client/postgresql"
	db "prod/pkg/client/postgresql/model"
)

type ProductStorage struct {
	queryBuilder sq.StatementBuilderType
	client       PostgreSQLClient
}

func NewProductStorage(client PostgreSQLClient) *ProductStorage {
	return &ProductStorage{
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		client:       client,
	}
}

const (
	scheme = "public"
	table  = "product"
)

func (s *ProductStorage) All(ctx context.Context) ([]model.Product, error) {
	query := s.queryBuilder.Select("id").
		Column("name").
		Column("description").
		Column("image_id").
		Column("price").
		Column("currency_id").
		Column("rating").
		Column("category_id").
		Column("specification").
		Column("created_at").
		Column("updated_at").
		From(scheme + "." + table)

	// TODO filtering and sorting

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := s.client.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	list := make([]model.Product, 0)
	for rows.Next() {
		p := model.Product{}
		if err = rows.Scan(
			&p.Id, &p.Name, &p.Description, &p.ImageId, &p.Price, &p.CurrencyId, &p.Rating, &p.CategoryId,
			&p.Specification, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			err = db.ErrScan(postgresql.ParsePgError(err))
			return nil, err
		}
		list = append(list, p)
	}

	return list, nil
}
