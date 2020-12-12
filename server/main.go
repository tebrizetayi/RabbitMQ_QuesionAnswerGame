package main

import (
	"log"
	"time"

	"github.com/streadway/amqp"
	"github.com/tebrizetayi/rabbitmq/utility"
)

func main() {
	questions := utility.CreateExam()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	utility.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	utility.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		utility.AnswersQ, // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	utility.FailOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	msgs, err := ch.Consume(q.Name,
		"",
		false,
		false,
		false,
		false,
		nil)
	utility.FailOnError(err, "There was a problem on consuming")

	for msg := range msgs {
		go func(msg amqp.Delivery) {
			log.Println(string(msg.Body), " received")
			err = ch.Publish("",
				msg.ReplyTo,
				false,
				false,
				amqp.Publishing{
					CorrelationId: msg.CorrelationId,
					Body:          []byte("Let me think"),
					ContentType:   "text/plain",
				})
			utility.FailOnError(err, "Failed to declare a queue")

			time.Sleep(5 * time.Second)
			err = ch.Publish("",
				//Send message to the queue , from where it came.
				msg.ReplyTo,
				false,
				false,
				amqp.Publishing{
					//Used for filtering in client side
					CorrelationId: msg.CorrelationId,
					Body:          []byte(questions[string(msg.Body)]),
					ContentType:   "text/plain",
				})
			utility.FailOnError(err, "Failed to declare a queue")
			msg.Ack(false)
		}(msg)
	}
}
