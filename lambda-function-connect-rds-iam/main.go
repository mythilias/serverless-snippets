package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	_ "github.com/go-sql-driver/mysql"
)

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(event *MyEvent) (map[string]interface{}, error) {

	var dbName string = "DatabaseName"
	var dbUser string = "DatabaseUser"
	var dbHost string = "mysqldb.123456789012.us-east-1.rds.amazonaws.com"
	var dbPort int = 3306
	var dbEndpoint string = fmt.Sprintf("%s:%d", dbHost, dbPort)
	var region string = "us-east-1"

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error: " + err.Error())
	}

	authenticationToken, err := auth.BuildAuthToken(
		context.TODO(), dbEndpoint, region, dbUser, cfg.Credentials)
	if err != nil {
		panic("failed to create authentication token: " + err.Error())
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=true&allowCleartextPasswords=true",
		dbUser, authenticationToken, dbEndpoint, dbName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	var sum int
	err = db.QueryRow("SELECT ?+? AS sum", 3, 2).Scan(&sum)
	if err != nil {
		panic(err)
	}
	s := fmt.Sprint(sum)
	message := fmt.Sprintf("The selected sum is: %s", s)

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	messageString := string(messageBytes)
	return map[string]interface{}{
		"statusCode": 200,
		"headers":    map[string]string{"Content-Type": "application/json"},
		"body":       messageString,
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
