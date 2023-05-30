package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"gopkg.in/gomail.v2"
)

const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587
const CONFIG_SENDER_NAME = "Bro.Inc <william16.lim@gmail.com>"

type Applicant struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	EmailTo string `json:"emailTo"`
}

func handler(applicant Applicant) (string, error) {

	to := applicant.EmailTo
	subject := "New applicant!"

	err := sendMail(to, subject, applicant)
	if err != nil {
		log.Fatal(err.Error())
	}

	return "success", nil
}

func queryDatabase() error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	svc := dynamodb.New(sess)
	input := &dynamodb.ListTablesInput{}

	fmt.Printf("Tables:\n")

	for {
		// Get the list of tables
		result, err := svc.ListTables(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeInternalServerError:
					fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return err
		}

		for _, n := range result.TableNames {
			fmt.Println(*n)
		}

		// assign the last read tablename as the start for our next call to the ListTables function
		// the maximum number of table names returned in a call is 100 (default), which requires us to make
		// multiple calls to the ListTables function to retrieve all table names
		input.ExclusiveStartTableName = result.LastEvaluatedTableName

		if result.LastEvaluatedTableName == nil {
			break
		}
	}
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
	// lambda.Start(handler)
	queryDatabase()
}
