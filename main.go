package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gopkg.in/gomail.v2"
)

const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587
const CONFIG_SENDER_NAME = "Bro.Inc <william16.lim@gmail.com>"

type Applicant struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	EmailTo    string `json:"emailTo"`
	Occupation string `json:"occupation"`
	Company    string `json:"company"`
}

type ApplicantRecord struct {
	Name       string
	Address    string
	Phone      string
	Email      string
	EmailTo    string
	Occupation string
	Company    string
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var applicant Applicant
	err := json.Unmarshal([]byte(request.Body), &applicant)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	err = sendMail(applicant.EmailTo, "New applicant!", applicant)
	if err != nil {
		log.Fatal(err.Error())
	}
	queryDatabase(applicant)

	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
	}
	return response, nil
}

func queryDatabase(applicant Applicant) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("KEY_ID"), os.Getenv("ACCESS_KEY"), ""),
	})
	if err != nil {
		return err
	}
	svc := dynamodb.New(sess)

	applicantRequest := ApplicantRecord(applicant)
	av, err := dynamodbattribute.MarshalMap(applicantRequest)
	if err != nil {
		log.Fatalf("Got error marshalling new movie item: %s", err)
	}

	tableName := "Applicant"
	fmt.Println(av)

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	fmt.Println("Successfully added '" + applicant.Email + " to " + tableName)

	return nil
}

func sendMail(to string, subject string, applicant Applicant) error {
	message := "<h1>Hello this is an automated message! Check out this new applicant</h1>" +
		"<table>" +
		"<tr>" +
		"<td>" +
		"Name" +
		"</td>" +
		"<td>" +
		applicant.Name +
		"</td>" +
		"</tr>" +
		"<tr>" +
		"<td>" +
		"Phone" +
		"</td>" +
		"<td>" +
		applicant.Phone +
		"</td>" +
		"</tr>" +
		"<tr>" +
		"<td>" +
		"Email" +
		"</td>" +
		"<td>" +
		applicant.Email +
		"</td>" +
		"</tr>" +
		"<tr>" +
		"<td>" +
		"Occupation" +
		"</td>" +
		"<td>" +
		applicant.Occupation +
		"</td>" +
		"</tr>" +
		"<tr>" +
		"<td>" +
		"Company" +
		"</td>" +
		"<td>" +
		applicant.Company +
		"</td>" +
		"</tr>" +
		"</table>"

	m := gomail.NewMessage()
	m.SetHeader("From", "william16.lim@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)

	d := gomail.NewDialer(CONFIG_SMTP_HOST, 587, os.Getenv("EMAIL"), os.Getenv("PASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	return nil
}

func main() {
	lambda.Start(handler)
	//Do not delete for testing locally
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	// app := Applicant{
	// 	Name:    "bob2",
	// 	Address: "bob",
	// 	Phone:   "01283212",
	// 	Email:   "william16.lim@gmail.com",
	// 	EmailTo: "william16.lim@gmail.com",
	// }
	// queryDatabase(app)
}
