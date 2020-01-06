package info

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"prodapp.com/app/util"
)

// its just naive approch to get db and redis conn to make it simple
// should be replaced with proper object attribute
var conn *sqlx.DB
var rdsConn util.RedisCache

func init() {
	conn = util.GetDatabaseConn()
	rdsConn = util.GetRedisConn()
}

type Registrar interface {
	New(int64) DataProp
}

type Options int64

const (
	BasicOpt Options = 1 << iota
	VariantOpt
	StatsOpt
)

type Getter struct {
	reg      Registrar
	prop     DataProp
	infoType Options
	cacheKey string
}

type QueryGroup struct {
	ids   []int64
	props map[int64]DataProp
}

type Info struct {
	Basic   Basic
	Stats   Stats
	Variant Variant
}

func GetInfos(productIDs []int64, o Options) ([]Info, error) {

	result := make([]Info, len(productIDs))
	getter := getSelectedInfo(o)

	infoGetter := make(map[int64][]Getter)
	cacheKeys := make(map[string][]string)

	// populate getter object base on registrar
	// generate cache key to append field into same parent key
	for _, pid := range productIDs {
		for _, g := range getter {
			g.prop = g.reg.New(pid)
			g.cacheKey = g.prop.GetCacheKey()
			infoGetter[pid] = append(infoGetter[pid], g)
			cacheKeys[g.cacheKey] = append(cacheKeys[g.cacheKey], g.prop.GetCacheFields()...)
		}
	}

	cacheResult, err := getFromCache(cacheKeys)
	if err != nil {
		fmt.Println(err)
	}

	queryGroups := make(map[string]*QueryGroup)

	// apply cache result data to every struct fields
	// if any failed process, fallback to database
	for _, pid := range productIDs {
		for i := range infoGetter[pid] {
			getter := infoGetter[pid][i]
			key := getter.cacheKey
			var isCompleted bool
			if res, ok := cacheResult[key]; ok {
				isCompleted = getter.prop.Apply(getter.prop, res)
			}

			if !isCompleted {
				// register incomplete redis result to queryGroup
				query := getter.prop.GetQuery()
				id := getter.prop.GetIdentifier()
				if _, exist := queryGroups[query]; exist {
					queryGroups[query].ids = append(queryGroups[query].ids, id)
					queryGroups[query].props[id] = getter.prop
				} else {
					prps := make(map[int64]DataProp)
					prps[id] = getter.prop
					queryGroups[query] = &QueryGroup{
						ids:   []int64{id},
						props: prps,
					}
				}
			}
		}
	}

	if len(queryGroups) > 0 {
		err = getFromDatabase(queryGroups)
		if err != nil {
			log.Println(err)
			return result, err
		}

		setCache(queryGroups)
	}

	// cast type to original and assign to result values
	i := 0
	for _, val := range infoGetter {
		info := Info{}
		for _, data := range val {
			switch data.infoType {
			case BasicOpt:
				if val, ok := data.prop.(*Basic); ok && val != nil {
					info.Basic = *val
				}
			case StatsOpt:
				if val, ok := data.prop.(*Stats); ok && val != nil {
					info.Stats = *val
				}
			case VariantOpt:
				if val, ok := data.prop.(*Variant); ok && val != nil {
					info.Variant = *val
				}
			}
		}

		result[i] = info
		i++
	}

	return result, nil
}

func getSelectedInfo(o Options) (registry []Getter) {

	if o&BasicOpt > 0 {
		registry = append(registry, Getter{
			reg:      new(BasicReg),
			infoType: BasicOpt,
		})
	}

	if o&VariantOpt > 0 {
		registry = append(registry, Getter{
			reg:      new(VariantReg),
			infoType: VariantOpt,
		})
	}

	if o&StatsOpt > 0 {
		registry = append(registry, Getter{
			reg:      new(StatsReg),
			infoType: StatsOpt,
		})
	}

	return
}

func getFromCache(mapKeys map[string][]string) (map[string]map[string]string, error) {

	result, err := rdsConn.MultiHashGetPipeline(mapKeys)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func getFromDatabase(queryGroup map[string]*QueryGroup) error {
	fmt.Println("start query", len(queryGroup))
	for query, mapper := range queryGroup {
		q, args, err := sqlx.In(query, mapper.ids)
		if err != nil {
			log.Println(err)
			return err
		}
		rows, err := conn.Queryx(conn.Rebind(q), args...)
		if err == nil && rows != nil {
			defer rows.Close()
			for rows.Next() {
				res, err := rows.SliceScan()
				if err != nil || len(res) < 1 {
					log.Println(err)
					return err
				}
				if id, ok := res[0].(int64); ok {
					if _, exist := mapper.props[id]; exist {
						target := mapper.props[id]
						err := rows.StructScan(target)
						if err != nil {
							log.Println("err scan", err)
							return err
						}

						err = target.PostQueryProcess()
						if err != nil {
							log.Println(err)
							return err
						}
					}
				}
			}
		} else if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func setCache(queryGroup map[string]*QueryGroup) {
	cacheData := make(map[string][]string)
	for _, v := range queryGroup {
		for _, prop := range v.props {
			cacheData[prop.GetCacheKey()] = append(cacheData[prop.GetCacheKey()], prop.GetCacheMap()...)
		}
	}

	errs := rdsConn.MultiHashSetPipeline(cacheData)
	if len(errs) > 0 {
		log.Println("Errors", errs)
	}
}
