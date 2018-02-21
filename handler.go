package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	address, user, password string
	config                  *configStruct
)

type dataRequest struct {
	Date     string `json:"date"`
	Datatype string `json:"datatype"`
}

type dataAnswer struct {
	Date     string   `json:"date"`
	Datatype string   `json:"datatype"`
	Data     []string `json:"data"`
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

func mariaConnect() {
	ReadConfig()
	db, err := sql.Open("mysql", user+":"+password+"@tcp("+address+":3306)/sensored")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()
	var version string
	db.QueryRow("SELECT VERSION()").Scan(&version)
	fmt.Println("Connected to:", version)
}

//Don't know if this needs to be exported or not, probably not though

func handleRequest(ctx context.Context, request dataRequest) (string, error) {
	return fmt.Sprintf("Data: %s %s", request.Date, request.Datatype), nil
}

func main() {
	mariaConnect()
	//lambda.Start(handleRequest)
}
