package postgres

import (
	"encoding/json"
	"github.com/Heilartin/bot_support/models"
)

func (db *DB) GetProductByTaskID(pid, storeID string) (*models.Product, error)  {
	query := `SELECT
			product.id   "product.id",
			product.store_id "product.store_id",
			product.wish_list_id "product.wish_list_id",
			product.access_key "product.access_key",
			product.task_id "product.task_id",
			product.pid "product.pid",
			product.name "product.name",
			product.image "product.image",
			product.price "product.price",
			product.symbol "product.symbol",
			product.stock_level "product.stock_level",
		
			COALESCE(json_agg(DISTINCT sizes.*), '[]')::json "sizes"
		FROM mrp_products AS product
				 left JOIN mrp_sizes as sizes ON
			( product.task_id = sizes.task_id)
		
		WHERE product.pid = $1 AND product.store_id = $2 GROUP BY "product.id";`
	var rawProduct models.RawProduct

	err := db.DB.Get(&rawProduct, query, pid, storeID)
	if err != nil {
		db.Logger.Error(err)
		return nil, err
	}
	p, err := db.processRawProduct(rawProduct)
	if err != nil {
		db.Logger.Error(err)
		return nil, err
	}
	return p, nil
}


func (db *DB) processRawProduct(p models.RawProduct) (*models.Product, error) {
	var sizes []*models.Size
	err := json.Unmarshal(p.Sizes, &sizes)
	if err != nil {
		return &models.Product{}, err
	}
	p.Product.Sizes = sizes
	return p.Product, nil
}

func (db *DB) GetPidByBrandName(name string) ([]string, error)  {
	var res []string
	query := `SELECT DISTINCT variant_part_number FROM mrp_scraper WHERE designer_name=$1;`
	err := db.DB.Select(&res, query, name)
	if err != nil {
		db.Logger.Error(err)
		return nil, err
	}
	return res, nil
}

