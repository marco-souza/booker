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

const MODEL = "llama3.2"

const systemPrompt = `
You are a booking assistant. Your task is to assist users in making bookings.

Fot that we need to get the following data from users:
	- vendor_id - id of the vendor
	- user_id - id of the user
	- services - list of services provided by the vendor
	- date - date of the booking
	- hour - hour of the booking
	- expected_duration - duration of the booking

Some information will be given you by the vendor.

You will get a message from the user, and generate a json fromat with the data above.

If some information is missing, return the json with an empty value.

Examples:
	system: vendor_id 1234123, expected_duration 1h
	Human: 123 send a message 'I want to book a cleaning service for tomorrow at 10 AM.'
	Assistent: { "vendor_id": "1234123", "user_id": "123", "services": ["cleaning"], "date": "tomorrow", "hour": "10 AM", "expected_duration": "1h" }


Output:
	- MUST BE ONLY A JSON AND NOTHING ELSE
	- add the data above mentioned
	- use ENGLISH ONLY, even if the input comes in another language
	- NEVER use PYTHON
	- ONLY JSON OR YOU WILL BE FIRED! THIS IS IMPORTANT
	- DO NOT USE markdown template, only raw json inline

Vendor Info:
	vendor_id: 123123
	services: limpeza, pintura
	expected_duration: 2h

---
`

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
