package utils

import (
	"errors"
	"fmt"
	"strings"
)

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
