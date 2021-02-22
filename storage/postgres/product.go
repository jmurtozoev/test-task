package postgres

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/jmurtozoev/test-task/models"
	"github.com/jmurtozoev/test-task/storage/repo"
)

type productRepo struct {
	db *sqlx.DB
}

func NewProductRepo(db *sqlx.DB) repo.Product {
	return &productRepo{db}
}

func (r *productRepo) Create(product models.Product) error {
	var query = `INSERT INTO products(name, price) VALUES($1, $2) ON CONFLICT DO UPDATE set price = EXCLUDED.price`

	_, err := r.db.Exec(query, product.Name, product.Price)
	return err
}

var allowedFilters = map[string]string{
	"cost_min": "price",
	"cost_max": "price",
	"name":     "name",
}

func (r *productRepo) List(page, limit int, filters map[string]interface{}) ([]models.Product, int, error) {
	var products []models.Product
	var count int

	offset := (page - 1) * limit
	var query = `SELECT id,
						name,
						price,
						update_count
					FROM products
					%s 
					ORDER BY name LIMIT ? OFFSET ?`

	query, varList, err := buildSearchQuery(query, filters, allowedFilters)
	if err != nil {
		return nil, count, err
	}

	varList = append(varList, limit, offset)

	rows, err := r.db.Queryx(query, varList...)
	switch err {
	case sql.ErrNoRows:
		return products, count, nil
	case nil:
		break
	default:
		return nil, count, err
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product
		err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.UpdateCount,
		)
		if err != nil {
			return nil, count, err
		}

		products = append(products, product)
	}

	// get rows count
	query = "SELECT COUNT(1) FROM products  %s "
	query, varList, err = buildSearchQuery(query, filters, allowedFilters)
	if err != nil {
		return nil, 0, err
	}

	row := r.db.QueryRow(query, varList...)
	err = row.Scan(&count)

	return products, count, err
}
