package info

import (
	"fmt"
	"strconv"
)

type BasicReg struct{}

func (r *BasicReg) New(productID int64) DataProp {
	b := new(Basic)
	b.Key = fmt.Sprintf("cache:pid:%d", productID)
	b.Fields = []string{"product_id", "shop_id", "child_cat_id", "name", "short_desc"}

	b.Query = `
		SELECT product_id, shop_id, child_cat_id, name, short_desc
		FROM ws_product
		WHERE product_id IN (?)
	`
	b.Identifier = productID
	return b
}

type Basic struct {
	CacheProp
	DBProp

	ProductID       int64  `db:"product_id" cache:"product_id" setter:"SetProductID"`
	ShopID          int64  `db:"shop_id" cache:"shop_id" setter:"SetShopID"`
	ChildCategoryID int64  `db:"child_cat_id" cache:"child_cat_id" setter:"SetChildCategoryID"`
	Name            string `db:"name" cache:"name" setter:"SetName"`
	ShortDesc       string `db:"short_desc" cache:"short_desc" setter:"SetShortDesc"`
}

func (b *Basic) SetProductID(value string) (err error) {
	b.ProductID, err = strconv.ParseInt(value, 10, 64)
	return err
}

func (b *Basic) SetShopID(value string) (err error) {
	b.ShopID, err = strconv.ParseInt(value, 10, 64)
	return err
}

func (b *Basic) SetChildCategoryID(value string) (err error) {
	b.ChildCategoryID, err = strconv.ParseInt(value, 10, 64)
	return err
}

func (b *Basic) SetName(value string) error {
	b.Name = value
	return nil
}

func (b *Basic) SetShortDesc(value string) error {
	b.ShortDesc = value
	return nil
}

func (b *Basic) GetCacheMap() []string {
	return []string{
		"product_id", strconv.FormatInt(b.ProductID, 10),
		"shop_id", strconv.FormatInt(b.ShopID, 10),
		"child_cat_id", strconv.FormatInt(b.ChildCategoryID, 10),
		"name", b.Name,
		"short_desc", b.ShortDesc,
	}
}
