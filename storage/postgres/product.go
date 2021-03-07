package postgres

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/jmurtozoev/test-task/models"
	"github.com/jmurtozoev/test-task/storage/repo"
	"gopkg.in/guregu/null.v4"
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
	var updatedAt null.Time

	offset := (page - 1) * limit
	var query = `SELECT id,
						count(1) OVER(),
						name,
						price,
						update_count,
						updated_at,
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
			&count,
			&product.Name,
			&product.Price,
			&product.UpdateCount,
			&updatedAt,
		)
		if err != nil {
			return nil, count, err
		}

		if updatedAt.Valid {
			product.UpdatedAt = updatedAt.Time.Format("2006-01-02 15:04:05")
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return products, count, err
}

func (r *productRepo) Get(productId int) (*models.Product, error) {
	var product models.Product

	query := `SELECT id, name, price, update_count FROM products WHERE id = $1;`

	err := r.db.QueryRowx(query, productId).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.UpdateCount)

	return &product, err
}

func (r *productRepo) Update(product *models.Product) error {
	query := `UPDATE products SET 
                    name = $2, 
                    price = $3, 
                    update_count = $4,
					updated_at = CURRENT_TIMESTAMP
				WHERE id = $1;`

	_, err := r.db.Exec(query,
		product.ID,
		product.Name,
		product.Price,
		product.UpdateCount,
		)

	return err
}

