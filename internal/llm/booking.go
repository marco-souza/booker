package llm

type BookingAttempt struct {
	VendorID         string   `json:"vendor_id"`
	UserID           string   `json:"user_id"`
	Services         []string `json:"services"`
	Date             string   `json:"date"`
	Hour             string   `json:"hour"`
	ExpectedDuration string   `json:"expected_duration"`
}

func (b *BookingAttempt) MissingFields() []string {
	fields := []string{}

	if b.VendorID == "" {
		fields = append(fields, "vendor_id")
	}

	if b.UserID == "" {
		fields = append(fields, "user_id")
	}

	if b.Services == nil || len(b.Services) == 0 {
		fields = append(fields, "services")
	}

	if b.Date == "" {
		fields = append(fields, "date")
	}

	if b.Hour == "" {
		fields = append(fields, "hour")
	}

	if b.ExpectedDuration == "" {
		fields = append(fields, "expected_duration")
	}

	return fields
}
