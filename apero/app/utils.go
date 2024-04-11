package app

import "log"

func check(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %w", message, err)
	}
}
