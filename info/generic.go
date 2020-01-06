package info

import (
	"reflect"
)

const (
	cacheTag  string = "cache"
	setterTag string = "setter"
)

// DataProp for cacheprod and dbprod
type DataProp interface {
	GetCacheKey() string
	GetCacheFields() []string
	Apply(interface{}, map[string]string) bool
	GetCacheMap() []string

	GetQuery() string
	GetIdentifier() int64
	PostQueryProcess() error
}

// CacheProp - this struct is for embedded to main struct data
type CacheProp struct {
	Key    string
	Fields []string

	applied int
}

func (c *CacheProp) GetCacheKey() string {
	return c.Key
}

func (c *CacheProp) GetCacheFields() []string {
	return c.Fields
}

// Apply - to fill in map cache to struct
// it will find map key based on struct field tag and use setter func to assign value
func (c *CacheProp) Apply(dest interface{}, data map[string]string) (isCompleted bool) {
	// get reflection type and value from struct parent
	objType := reflect.TypeOf(dest)
	objVals := reflect.ValueOf(dest)
	for i := 0; i < objType.Elem().NumField(); i++ {
		// get struct field attribute
		cacheField := objType.Elem().Field(i).Tag.Get(cacheTag)
		setterName := objType.Elem().Field(i).Tag.Get(setterTag)
		// ensue field has cache and setter tag
		if cacheField == "" || setterName == "" {
			continue
		}

		// ensure setter method exist before invoke
		_, setterExist := objType.MethodByName(setterName)
		if !setterExist {
			continue
		}

		//  method setter must has only 1 string arg and 1 return field
		setterMethod := objVals.MethodByName(setterName)
		if setterMethod.Type().NumIn() != 1 || setterMethod.Type().In(0).Kind() != reflect.String || setterMethod.Type().NumOut() != 1 {
			continue
		}

		// get value from cache map
		cacheVal, cacheExist := data[cacheField]
		if !cacheExist {
			continue
		}

		// invoke setter method
		result := setterMethod.Call([]reflect.Value{reflect.ValueOf(cacheVal)})
		if len(result) > 0 && result[0].IsNil() {
			c.applied++
		}

	}

	return c.applied == len(c.Fields)
}

// GetCacheMap - abstract method, override this func on child
// use to get map of struct field for set cache
// value should be []string{field1, value1, field2, value2...}
func (c *CacheProp) GetCacheMap() []string {
	return []string{}
}

type DBProp struct {
	Query      string
	Identifier int64
	IDCatcher  int64 `db:"product_id"` //generalize mandatory field for catching structScan from sqlx
}

func (c *DBProp) GetQuery() string {
	return c.Query
}

func (c *DBProp) GetIdentifier() int64 {
	return c.Identifier
}

// PostQueryProcess - abstract method, override this func on child
// use to post processing after query db scan
// usually for formatting some fields
func (c *DBProp) PostQueryProcess() error {
	return nil
}
