package services

import (
	"context"
	"database/sql"
	"fmt"

	cm "Hanif_Aulia_Sabri-MyTrip/git/order/common"

	_ "github.com/go-sql-driver/mysql"
)

func (PaymentService) FastPayHandler(ctx context.Context, req cm.FastPayRequest) (res cm.FastPayResponse) {

	defer panicRecovery()

	host := cm.Config.Connection.Host
	port := cm.Config.Connection.Port
	user := cm.Config.Connection.User
	pass := cm.Config.Connection.Password
	data := cm.Config.Connection.Database

	var mySQL = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", user, pass, host, port, data)

	db, err = sql.Open("mysql", mySQL)

	if err != nil {
		panic(err.Error())
	}

	var fasResponse cm.FastPayResponse
	var list cm.PaymentChannel

	sql := `SELECT
				pg_code,
				IFNULL(pg_name,'')
			FROM list_payment WHERE merchant_id = ?`

	result, err := db.Query(sql, req.MerchantID)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		err := result.Scan(&list.PgCode, &list.PgName)

		if err != nil {
			panic(err.Error())
		}

		fasResponse.PaymentChannel = append(fasResponse.PaymentChannel, list)

	}

	fasResponse.Merchant = req.Merchant
	fasResponse.MerchantID = req.MerchantID

	res = fasResponse

	return
}
