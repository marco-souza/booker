package main

import (
	"encoding/json"
	"fmt"
	"log"
	"marco-souza/booker/internal/llm"
)

func main() {
	ba := llm.FeedbackLoop()
	fmt.Println("Final Booking Attempt:", ba)

	msg1 := "123123123@whatsapp.us: Hola, necessito agendar una limpeza"
	bookingAttempt := llm.UpdateBooking(msg1, "{}")

	fmt.Println("Booking Attempt:", bookingAttempt)

	draftData, err := json.Marshal(bookingAttempt)
	if err != nil {
		log.Fatal(err)
	}

	msg2 := "123123123@whatsapp.us: para manhana a las 10 AM"
	bookingAttempt = llm.UpdateBooking(msg2, string(draftData))

	fmt.Println("Booking Attempt:", bookingAttempt)
}
