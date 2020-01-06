package bench

import (
	"prodapp.com/app/info"
)

func SetDirectly(target *info.Basic, data map[string]string) {
	if val, ok := data["product_id"]; ok {
		target.SetProductID(val)
	}

	if val, ok := data["shop_id"]; ok {
		target.SetShopID(val)
	}

	if val, ok := data["child_cat_id"]; ok {
		target.SetChildCategoryID(val)
	}

	if val, ok := data["name"]; ok {
		target.SetName(val)
	}

	if val, ok := data["short_desc"]; ok {
		target.SetShortDesc(val)
	}

	// fmt.Printf("%+v\n", target)

}

func SetWithReflection(target *info.Basic, data map[string]string) {
	target.Apply(target, data)

	// fmt.Printf("%+v\n", target)
}
