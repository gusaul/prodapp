package main

import (
	"fmt"

	"prodapp.com/app/info"
)

func main() {

	data, err := info.GetInfos([]int64{123, 456}, info.BasicOpt|info.StatsOpt|info.VariantOpt)
	fmt.Println("result:", err)

	for _, v := range data {
		fmt.Println("product_id:", v.Basic.ProductID)
		fmt.Println("shop_id:", v.Basic.ShopID)
		fmt.Println("child_cat_id:", v.Basic.ChildCategoryID)
		fmt.Println("name:", v.Basic.Name)
		fmt.Println("short_desc:", v.Basic.ShortDesc)
		fmt.Println("transaction_success:", v.Stats.TransactionSuccess)
		fmt.Println("transaction_reject:", v.Stats.TransactionReject)
		fmt.Println("count_sold:", v.Stats.CountSold)
		fmt.Println("parent_id:", v.Variant.ParentID)
		fmt.Println("is_parent:", v.Variant.IsParent)
		fmt.Println("is_variant:", v.Variant.IsVariant)
		fmt.Println("children_ids:", v.Variant.ChildrenIDs)
		fmt.Println("==============================")
	}
}
