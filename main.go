package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
)

var (
	address, user, password string
	config                  *configStruct
)

type Request struct {
	MAC      string `json:"mac"`
	Date     string `json:"date"`
	Datatype string `json:"datatype"`
}

type Response struct {
	Data []string `json:"data"`
	Ok   bool
}

type configStruct struct {
	Address  string `json:"address"`
	User     string `json:"user"`
	Password string `json:"password"`
}

//go away
func ReadConfig() error {

	file, err := ioutil.ReadFile("./config.json")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = json.Unmarshal(file, &config)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	address = config.Address
	user = config.User
	password = config.Password

	return nil
}

func queryMaria(req Request) (array []string) {
	//ReadConfig()
	//select temperature from b827eb06efa4 where datetime like '08/03/2018%';
	selectData := "select " + req.Datatype + " from " + req.MAC + " where datetime like " + req.Date + " %"
	user := os.Getenv("user")
	password := os.Getenv("password")
	address := os.Getenv("address")
	db, err := sql.Open("mysql", user+":"+password+"@tcp("+address+":3306)/SensorEdAWS")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()
	var Data []string
	rows, err := db.Query(selectData)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err = rows.Scan(&req.Datatype)
		if err != nil {
			panic(err)
		}
		Data = append(Data, req.Datatype)
	}
	return Data
}

//Don't know if this needs to be exported or not, probably not though

func handleRequest(request Request) []byte {
	data := queryMaria(request)
	var response *Response
	response.Data = data
	response.Ok = true

	r, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	return r
}

func main() {
	lambda.Start(handleRequest)
}
