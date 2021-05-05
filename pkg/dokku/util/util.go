package util

import (
	"log"
	"net"
	"time"
)

// Linearly retry network errors forever
func RetriedNetworkFunc(f func() (interface{}, error)) (interface{}, error) {
	res, err := f()
	wait := time.Duration(2500) * time.Millisecond
	attempt := 1
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
		log.Printf("Network error, waiting 2.5s before trying again (attempt %d)\n", attempt)
		time.Sleep(wait)
		attempt++
		res, err = f()
	}
	// unreachable
}
