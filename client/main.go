package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/streadway/amqp"
	"github.com/tebrizetayi/rabbitmq/utility"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
func main() {
	questions := utility.CreateExam()
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	ch.Confirm(false)

	q, err := ch.QueueDeclare(
		utility.QuestionQ, // name
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	failOnError(err, "Failed to declare a queue")

	/*
		err = ch.Qos(
			1,     // prefetch count
			0,     // prefetch size
			false, // global
		)*/
	//It is a listener, in order to get confirmation message , whether message is received by message broker
	confirm := ch.NotifyPublish(make(chan amqp.Confirmation))

	//Correlation Id is used in order to filter the incoming messages.
	corrID := rand.Intn(time.Now().Nanosecond())
	go func() {
		for {
			body := utility.PickRandomKey(questions)
			err = ch.Publish(
				"",               // exchange
				utility.AnswersQ, // routing key
				false,            // mandatory
				false,            // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					Body:          []byte(fmt.Sprintf("%v", body)),
					CorrelationId: strconv.Itoa(corrID),
					ReplyTo:       q.Name,
				})
			failOnError(err, "Failed to publish a message")

			//Waiting for the response whether message is received by Message Broker
			confirmed := <-confirm
			if confirmed.Ack {
				log.Println(body)
			} else {
				log.Println("Message could not sent to broker failed", body)
			}
			time.Sleep(10 * time.Second)
		}
	}()

	msgs, err := ch.Consume(q.Name,
		"",
		false,
		false,
		false,
		false,
		nil)
	failOnError(err, "There was a problem on consuming")

	for msg := range msgs {
		//Message Filtering
		//log.Println(msg.CorrelationId, " ", strconv.Itoa(int(corrID)), string(msg.Body))
		if msg.CorrelationId == strconv.Itoa(int(corrID)) {
			log.Println(string(msg.Body))
			msg.Ack(false)
		}
	}
}
