package llm

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

var llm *ollama.LLM

func FeedbackLoop() *BookingAttempt {
	fmt.Println("Ola, como posso ajudar?")

	bookingAttempt := &BookingAttempt{}

	for {
		input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Println("Input:", input)

		data, _ := json.Marshal(bookingAttempt)
		previousDraft := string(data)
		bookingAttempt = UpdateBooking(input, previousDraft)

		missingFields := bookingAttempt.MissingFields()
		if len(missingFields) == 0 {
			// if no missing field, return
			break
		}

		ask := RequestMissingData(missingFields)
		fmt.Println(ask)
	}

	return bookingAttempt
}

func RequestMissingData(missingFields []string) string {
	ctx := context.Background()
	ask, err := llm.Call(
		ctx,
		fmt.Sprintf("You are a booking assisntent! Ask the user to provive %s information in a friendly manner", missingFields[0]),
		llms.WithTemperature(0.8),
	)
	if err != nil {
		log.Fatal(err)
	}

	return ask
}

func UpdateBooking(message string, draftData string) *BookingAttempt {
	ctx := context.Background()
	bookingDraft, err := llm.Call(
		ctx,
		fmt.Sprintf("%s\nBookingDraft: %s\n Human: %s\nAssistent: ", systemPrompt, draftData, message),
		llms.WithTemperature(0.8),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Booking Draft:", bookingDraft)

	var bookingAttempt BookingAttempt
	if err := json.Unmarshal([]byte(bookingDraft), &bookingAttempt); err != nil {
		log.Fatal(err)
	}

	return &bookingAttempt
}

func init() {
	l, err := ollama.New(ollama.WithModel(MODEL))
	if err != nil {
		log.Fatal(err)
	}

	llm = l
}
