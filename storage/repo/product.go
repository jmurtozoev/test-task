package repo

import "github.com/jmurtozoev/test-task/models"

type Product interface {
	Create(product models.Product) error
	List(page, limit int, filter map[string]interface{}) ([]models.Product, int, error)
	Get(productId int) (*models.Product, error)
	Update(product *models.Product) error
}
