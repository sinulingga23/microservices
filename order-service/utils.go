package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Shopify/sarama"
)

func connectProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll

	syncProducer, errNewSyncProducer := sarama.NewSyncProducer([]string{
		"kafka:9092",
	}, config)
	if errNewSyncProducer != nil {
		return nil, errNewSyncProducer
	}

	return syncProducer, nil
}

func PublishMessage(topic string, message []byte) error {
	producer, errConnectProducer := connectProducer()
	if errConnectProducer != nil {
		return errConnectProducer
	}
	defer func() {
		if errClose := producer.Close(); errClose != nil {
			log.Printf("Error when close producer: %v", errClose)
		}
	}()

	messageTopic := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, errSendMessage := producer.SendMessage(messageTopic)
	if errSendMessage != nil {
		return errSendMessage
	}

	log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}

func GenerateOrderId(currentSize int) (string, error) {
	if currentSize <= -1 {
		return "", errors.New("currentSize should greaer than -1")
	}

	currentSize += 1
	tempCurrentSize := currentSize
	countDigit := 0
	for currentSize != 0 {
		currentSize /= 10
		countDigit += 1
	}

	if tempCurrentSize == 0 {
		tempCurrentSize = 1
		countDigit = 1
	}

	totalDigitZero := 6
	totalDigitZero -= countDigit
	return fmt.Sprintf("ORD%s%d", strings.Repeat("0", totalDigitZero), tempCurrentSize), nil
}

func GenerateOrderDetailId(currentSize int) (string, error) {
	if currentSize <= -1 {
		return "", errors.New("currentSize should greaer than -1")
	}

	currentSize += 1
	tempCurrentSize := currentSize
	countDigit := 0
	for currentSize != 0 {
		currentSize /= 10
		countDigit += 1
	}

	if tempCurrentSize == 0 {
		tempCurrentSize = 1
		countDigit = 1
	}

	totalDigitZero := 6
	totalDigitZero -= countDigit
	return fmt.Sprintf("ODD%s%d", strings.Repeat("0", totalDigitZero), tempCurrentSize), nil
}
