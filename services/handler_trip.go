package services

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	cm "Hanif_Aulia_Sabri-MyTrip/git/order/common"

	_ "github.com/go-sql-driver/mysql"
)

func (PaymentService) TripsHandler(ctx context.Context, req cm.MyTripsrequest) (res cm.MytripsResponse) {

	msg := &cm.MyTripsrequest{
		Provinsi:      req.Provinsi,
		DepatureDate1: req.DepatureDate1,
		DepatureDate2: req.DepatureDate2,
	}

	reqBody, err := json.Marshal(msg)

	if err != nil {
		print(err)
	}

	resp, err := http.Post("http://35.186.147.192/travel/GetTripsSample.php", "application/json", bytes.NewBuffer(reqBody))

	if err != nil {
		print(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		print(err)
	}
	fmt.Println(string(body))
	var response cm.MytripsResponse
	json.Unmarshal(body, &response)

	res.Message = response.Message
	res.Status = response.Status
	res.TripDetail = response.TripDetail

	host := cm.Config.Connection.Host
	port := cm.Config.Connection.Port
	user := cm.Config.Connection.User
	pass := cm.Config.Connection.Password
	data := cm.Config.Connection.Database
	var mySQL = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", user, pass, host, port, data)
	db, err := sql.Open("mysql", mySQL)
	if err != nil {
		panic(err.Error())
	}

	for _, data := range response.TripDetail {
		AirlineName := data.AirlineName
		AirportName := data.AirportName
		CityName := data.CityName

		fmt.Println("AirlineName : ", AirlineName)
		fmt.Println("AirlineName : ", AirportName)
		fmt.Println("AirlineName : ", CityName)

		sql := "INSERT INTO `trip` (`AirlineName`, `AirportName`, `CityName`) VALUES (?, ?, ?)"
		stmt, err := db.Prepare(sql)
		if err != nil {
			panic(err.Error())
		}

		_, err = stmt.Exec(AirlineName, AirportName, CityName)

	}

	return
}
