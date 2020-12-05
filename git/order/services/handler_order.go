package services

import (
	"context"
	"database/sql"
	"fmt"

	cm "Hanif_Aulia_Sabri-MyTrip/git/order/common"

	_ "github.com/go-sql-driver/mysql"
)

func (PaymentService) OrderHandler(ctx context.Context, req cm.Message) (res cm.Message) {

	var db *sql.DB
	var err error

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

	res.OrderID = req.OrderID

	var order cm.Orders
	var orderdet cm.OrdersDetail

	sql := `SELECT
				OrderID,
				IFNULL(CustomerID,'') CustomerID,
				IFNULL(EmployeeID,'') EmployeeID,
				IFNULL(OrderDate,'') OrderDate				
			FROM orders WHERE OrderID = ?`

	result, err := db.Query(sql, req.OrderID)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		err := result.Scan(&order.OrderID, &order.CustomerID, &order.EmployeeID, &order.OrderDate)

		if err != nil {
			panic(err.Error())
		}

		sqlDetial := `SELECT
						order_details.OrderID		
						, products.ProductID
						, products.ProductName
						, order_details.UnitPrice
						, order_details.Quantity
					FROM
						order_details
						INNER JOIN products 
							ON (order_details.ProductID = products.ProductID)
					WHERE order_details.OrderID	= ?`

		orderID := &order.OrderID
		fmt.Println(*orderID)
		resultDetail, errDet := db.Query(sqlDetial, *orderID)

		defer resultDetail.Close()

		if errDet != nil {
			panic(err.Error())
		}

		for resultDetail.Next() {

			err := resultDetail.Scan(&orderdet.OrderID, &orderdet.ProductID, &orderdet.ProductName, &orderdet.UnitPrice, &orderdet.Quantity)

			if err != nil {
				panic(err.Error())
			}

			order.OrdersDet = append(order.OrdersDet, orderdet)

		}

	}

	if &order != nil {
		res.Code = 100
		res.Remark = "Success"
	}

	res.Orders = &order

	return
}
