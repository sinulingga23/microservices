package utils

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
)

func connectConsumer() (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, errNewConsumer := sarama.NewConsumer([]string{"kafka:9092"}, config)
	if errNewConsumer != nil {
		return nil, errNewConsumer
	}

	return consumer, nil
}

func ReceiveMessage(topic string) {
	consumer, errConnectConsumer := connectConsumer()
	if errConnectConsumer != nil {
		log.Printf("errConnectConsumer: %v", errConnectConsumer)
	}

	if errConnectConsumer == nil {
		partitionConsumer, errConsumerPartition := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
		if errConsumerPartition != nil {
			log.Printf("errConsumerPartition: %v", errConsumerPartition)
		}

		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

		// Get signal for finish
		doneCh := make(chan struct{})
		go func() {
			for {
				select {
				case errErrors := <-partitionConsumer.Errors():
					log.Printf("errErrors: %v", errErrors)
				case msg := <-partitionConsumer.Messages():
					log.Printf("msg: %s", string(msg.Value))
				case <-sigchan:
					doneCh <- struct{}{}
				}
			}
		}()

		<-doneCh

		if errClose := partitionConsumer.Close(); errClose != nil {
			log.Printf("errClose: %v", errClose)
		}
	}
}
