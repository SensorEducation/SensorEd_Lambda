package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
)

var (
	address, user, password string
	config                  *configStruct
)

type Request struct {
	MAC      string `json:\"mac\"`
	Date     string `json:\"date\"`
	Datatype string `json:\"datatype\"`
}

type Response struct {
	Base64     bool          `json:"isBase64Encoded"`
	StatusCode int           `json:"statusCode"`
	Headers    *headerStruct `json:"headers"`
	Body       string        `json:"body"`
}
type dataStruct struct {
	Data []float64 `json:"data"`
	Time []string  `json:"time"`
}

type headerStruct struct {
	Methods    string `json:"Access-Control-Allow-Methods"`
	CORS       string `json:"Access-Control-Allow-Origin"`
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

func queryMaria(req Request) (data []float64, time []string) {
	//ReadConfig()
	//select temperature from b827eb06efa4 where datetime like '08/03/2018%';
	selectData := `SELECT ` + req.Datatype + `, timeNow FROM ` + req.MAC + ` WHERE dateNow="` + req.Date + `";`
	//selectData := `SELECT temperature FROM b827eb06efa4 WHERE datetime LIKE "08/03/2018%";`
	user := os.Getenv("user")
	password := os.Getenv("password")
	address := os.Getenv("address")
	db, err := sql.Open("mysql", user+":"+password+"@tcp("+address+":3306)/SensorEdAWS")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()
	var Data []float64
	var Time []string
	rows, err := db.Query(selectData)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var data float64
		var time string
		err = rows.Scan(&data, &time)
		if err != nil {
			panic(err)
		}

		Data = append(Data, data)
		Time = append(Time, time)
	}
	return Data, Time
}

//Don't know if this needs to be exported or not, probably not though

func handleRequest(req events.APIGatewayProxyRequest) (Response, error) {
	var request Request
	err := json.Unmarshal([]byte(req.Body), &request)
	if err != nil {
		log.Fatal(err)
	}
	data, time := queryMaria(request)
	dataResp, err := json.Marshal(&dataStruct{Data: data, Time: time})
	if err != nil {
		log.Fatal(err)
	}
	//stringResp, err := json.Marshal(&dataStruct{Data: data})
	return Response{
			Base64:     false,
			StatusCode: 200,
			Headers:    &headerStruct{Headertype: "application/json", CORS: "*", Methods: "*"},
			Body:       string(dataResp),
		},
		nil
}

func main() {
	lambda.Start(handleRequest)
}
