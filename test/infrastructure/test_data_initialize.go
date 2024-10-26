package infrastructure

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
)

var INSERT_PRODUCTS = `INSERT INTO products (name, price, discount,store) 
VALUES('360Hz 24" Monitor',2000.0, 12.0, 'BENQ'),
('EC-2B Mouse',1200.0, 10.0, 'Zowie'),
('RTX 5090',10000.0, 20.0, 'Nvidia'),
('Iphone 17',3000.0, 0.0, 'Apple');
`

func TestDataInitialize(ctx context.Context, dbPool *pgxpool.Pool) {
	insertProductsResult, insertProductsErr := dbPool.Exec(ctx, INSERT_PRODUCTS)
	if insertProductsErr != nil {
		log.Error(insertProductsErr)
	} else {
		log.Info(fmt.Sprintf("Products data created with %d rows", insertProductsResult.RowsAffected()))
	}
}
