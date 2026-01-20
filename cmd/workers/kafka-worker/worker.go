package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"paytm-project/boot"
	"paytm-project/internal/clients"
	kafkaclient "paytm-project/internal/clients/kafka_client"
	"paytm-project/internal/db"
	"paytm-project/redis"

	"github.com/joho/godotenv"
)

const (
	HTML = "html"
	TEXT = "text"
)

type EmailEvent struct {
	To        string `json:"to"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
	EmailType string `json:"emailType"`
}

func handleNotification(ctx context.Context, clientRegistery *clients.ClientRegistry) func(key, value []byte) {
	return func(key, value []byte) {
		var event EmailEvent
		sender := os.Getenv("GOOGLE_EMAIL_ADDRESS")
		err := json.Unmarshal(value, &event)
		if err != nil {
			log.Println("❌ Error unmarshalling message:", err)
			return
		}
		switch event.EmailType {
		case HTML:
			err = clientRegistery.EmailClient.SendHtmlEmail(ctx, sender, event.To, event.Subject, event.Body)
			if err != nil {
				log.Println("❌ Error sending email:", err)
				return
			}
		case TEXT:
			err = clientRegistery.EmailClient.SendEmail(ctx, sender, event.To, event.Subject, event.Body)
			if err != nil {
				log.Println("❌ Error sending email:", err)
				return
			}
		default:
			err = clientRegistery.EmailClient.SendEmail(ctx, sender, event.To, event.Subject, event.Body)
			if err != nil {
				log.Println("❌ Error sending email:", err)
				return
			}
		}
	}
}

func main() {
	if os.Getenv("APP_ENV") != "prod" {
		err := godotenv.Load(".env")
		if err != nil {
			panic("env file not found")
		}
	}
	db.InitDb()
	redis.InitRedis()
	clientRegistery := boot.InitClientRegistery(db.GetDbConnection(), redis.GetRedisClientConnection())
	consumer := kafkaclient.GetKafkaClient("localhost:9092", "notifications", "notification-group")
	defer consumer.Stop()

	log.Println("👂 Listening for Kafka messages...")
	log.Println("yes new build ;)")
	ctx := context.Background()
	consumer.Listen(ctx, handleNotification(ctx, clientRegistery))
}
