package util

import (
	"log"
	"net"
	"time"
)

// Exponentially backoff network errors
func RetriedNetworkFunc(f func() (interface{}, error)) (interface{}, error) {
	// maximum backoff
	maxBackoff := 30
	currentBackoff := 1
	currentBackoffCounter := 0
	res, err := f()
	for {
		if err == nil {
			return res, err
		}
		switch err.(type) {
		case net.Error:
			// fallthrough to backoff
		default:
			return res, err
		}
		// backoff
		log.Printf("Network error, waiting %d seconds before trying again (attempt %d)\n", currentBackoff, currentBackoffCounter)
		time.Sleep(time.Duration(currentBackoff) * time.Second)
		currentBackoffCounter++
		if currentBackoffCounter == 2 {
			if currentBackoff == maxBackoff {
				break
			}
			currentBackoffCounter = 0
			currentBackoff *= 2
			if currentBackoff > maxBackoff {
				currentBackoff = maxBackoff
			}
		}
		res, err = f()
	}
	log.Printf("Network error, given up on retries (backoff time %d, attempt %d)\n", currentBackoff, currentBackoffCounter)
	return res, err
}
