package info

import (
	"fmt"
	"strconv"
)

type StatsReg struct{}

func (r *StatsReg) New(productID int64) DataProp {
	s := new(Stats)
	s.Key = fmt.Sprintf("stats:pid:%d", productID)
	s.Fields = []string{"transaction_success", "transaction_reject", "count_sold"}

	s.Query = `
		SELECT product_id, transaction_success, transaction_reject, count_sold
		FROM db_product_stats
		WHERE product_id IN (?)
	`
	s.Identifier = productID
	return s
}

type Stats struct {
	CacheProp
	DBProp

	TransactionSuccess int64 `db:"transaction_success" cache:"transaction_success" setter:"SetTransactionSuccess"`
	TransactionReject  int64 `db:"transaction_reject" cache:"transaction_reject" setter:"SetTransactionReject"`
	CountSold          int64 `db:"count_sold" cache:"count_sold" setter:"SetCountSold"`
}

func (s *Stats) SetTransactionSuccess(value string) (err error) {
	s.TransactionSuccess, err = strconv.ParseInt(value, 10, 64)
	return err
}

func (s *Stats) SetTransactionReject(value string) (err error) {
	s.TransactionReject, err = strconv.ParseInt(value, 10, 64)
	return err
}

func (s *Stats) SetCountSold(value string) (err error) {
	s.CountSold, err = strconv.ParseInt(value, 10, 64)
	return err
}

func (s *Stats) GetCacheMap() []string {
	return []string{
		"transaction_success", strconv.FormatInt(s.TransactionSuccess, 10),
		"transaction_reject", strconv.FormatInt(s.TransactionReject, 10),
		"count_sold", strconv.FormatInt(s.CountSold, 10),
	}
}
