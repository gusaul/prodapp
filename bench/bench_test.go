package bench

import (
	"testing"

	"prodapp.com/app/info"
)

var data = map[string]string{
	"product_id":   "123",
	"shop_id":      "321",
	"child_cat_id": "456",
	"name":         "so klin",
	"short_desc":   "detergen",
}

func BenchmarkSetDirectly(b *testing.B) {
	for i := 0; i < b.N; i++ {
		basic := new(info.Basic)
		SetDirectly(basic, data)
	}
}

func BenchmarkSetWithReflection(b *testing.B) {
	for i := 0; i < b.N; i++ {
		basic := new(info.Basic)
		SetWithReflection(basic, data)
	}
}
