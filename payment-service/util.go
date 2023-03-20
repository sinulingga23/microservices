package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
)

func connectProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	conn, errNewSyncProducer := sarama.NewSyncProducer([]string{"kafka:9092"}, config)
	if errNewSyncProducer != nil {
		return nil, errNewSyncProducer
	}

	return conn, nil
}

func PublishMessage(topic string, message []byte) error {
	producer, errConnectProducer := connectProducer()
	if errConnectProducer != nil {
		return errConnectProducer
	}
	defer producer.Close()

	messageProducer := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, errSendMessage := producer.SendMessage(messageProducer)
	if errSendMessage != nil {
		return errSendMessage
	}

	log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}

func connectConsumer() (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	conn, errNewConsumer := sarama.NewConsumer([]string{"localhost:29092"}, config)
	if errNewConsumer != nil {
		return nil, errNewConsumer
	}

	return conn, nil
}

func ReceiveMessage() {
	topic := "deduct-qtty-product-for-order"
	worker, err := connectConsumer()
	if err != nil {
		panic(err)
	}

	// Calling ConsumePartition. It will open one connection per broker
	// and share it for all partitions that live on it.
	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}
	fmt.Println("Consumer started ")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	// Count how many message processed
	msgCount := 0

	// Get signal for finish
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				msgCount++
				fmt.Printf("Received message Count %d: | Topic(%s) | Message(%s) \n", msgCount, string(msg.Topic), string(msg.Value))
			case <-sigchan:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	fmt.Println("Processed", msgCount, "messages")

	if err := worker.Close(); err != nil {
		panic(err)
	}
}
