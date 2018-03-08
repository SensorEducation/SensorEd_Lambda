package main

import (
	"context"
	"database/sql"
	"fmt"
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
	Base64     bool          `json:"isBase64Encoded"`
	StatusCode int           `json:"statusCode"`
	Headers    *headerStruct `json:"headers"`
	Body       *dataStruct   `json:"body"`
}
type dataStruct struct {
	Data []string `json:"data"`
}

type headerStruct struct {
	Headertype string `json:"Content-Type"`
}

type configStruct struct {
	Address  string `json:"address"`
	User     string `json:"user"`
	Password string `json:"password"`
}

/*
func ReadConfig() error {

	file, err := ioutil.ReadFile("./config.json")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = json.Unmarshal(file, &config)

	if err != nil {
		fmt.Println(err.Error())
		return errfmt.Sprintf(
	}
	address = config.Address
	user = config.User
	password = config.Password

	return nil
}
*/

func queryMaria(req Request) (array []string) {
	//ReadConfig()
	//select temperature from b827eb06efa4 where datetime like '08/03/2018%';
	selectData := `SELECT ` + req.Datatype + ` FROM ` + req.MAC + ` WHERE dateTime LIKE "` + req.Date + `%";`
	//selectData := `SELECT temperature FROM b827eb06efa4 WHERE datetime LIKE "08/03/2018%";`
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

func handleRequest(ctx context.Context, request Request) (Response, error) {
	data := queryMaria(request)
	return Response{
			Base64:     false,
			StatusCode: 200,
			Headers:    &headerStruct{Headertype: "application/json"},
			Body:       &dataStruct{Data: data},
		},
		nil
}

func main() {
	lambda.Start(handleRequest)
}
