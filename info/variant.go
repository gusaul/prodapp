package info

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type VariantReg struct{}

func (r *VariantReg) New(productID int64) DataProp {
	v := new(Variant)
	v.Key = fmt.Sprintf("cache:pid:%d", productID)
	v.Fields = []string{"parent_id", "is_parent", "is_variant", "children_ids"}

	v.Query = `
		SELECT product_id, parent_id, is_parent, is_variant, children_ids
		FROM ws_product_tree
		WHERE product_id IN (?)
	`
	v.Identifier = productID
	return v
}

type Variant struct {
	CacheProp
	DBProp

	ParentID       int64   `db:"parent_id" cache:"parent_id" setter:"SetProductID"`
	IsParent       bool    `db:"is_parent" cache:"is_parent" setter:"SetIsParent"`
	IsVariant      bool    `db:"is_variant" cache:"is_variant" setter:"SetIsVariant"`
	ChildrenIDsRaw string  `db:"children_ids"`
	ChildrenIDs    []int64 `cache:"children_ids" setter:"SetChildrenIDs"`
}

func (v *Variant) SetProductID(value string) (err error) {
	v.ParentID, err = strconv.ParseInt(value, 10, 64)
	return err
}

func (v *Variant) SetIsParent(value string) (err error) {
	v.IsParent, err = strconv.ParseBool(value)
	return err
}

func (v *Variant) SetIsVariant(value string) (err error) {
	v.IsVariant, err = strconv.ParseBool(value)
	return err
}

func (v *Variant) SetChildrenIDs(value string) (err error) {
	return json.Unmarshal([]byte(value), &v.ChildrenIDs)
}

func (v *Variant) PostQueryProcess() (err error) {
	raws := strings.Split(v.ChildrenIDsRaw, ",")
	v.ChildrenIDs = make([]int64, len(raws))
	for i, r := range raws {
		v.ChildrenIDs[i], err = strconv.ParseInt(r, 10, 64)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *Variant) GetCacheMap() []string {
	childIds := "[]"
	if len(v.ChildrenIDs) > 0 {
		cids, _ := json.Marshal(v.ChildrenIDs)
		childIds = string(cids)
	}
	return []string{
		"parent_id", strconv.FormatInt(v.ParentID, 10),
		"is_parent", strconv.FormatBool(v.IsParent),
		"is_variant", strconv.FormatBool(v.IsVariant),
		"children_ids", childIds,
	}
}
