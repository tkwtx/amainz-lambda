package main

import (
	"context"
	"testing"
)

func Test_Main(t *testing.T) {
	ctx := context.Background()
	cfg := ReserveConfig{
		LastName:        "l",
		FirstName:       "f",
		LastNameKana:    "lk",
		FirstNameKana:   "fl",
		MailAddress:     "sample@aa.jp",
		Tel:             "08012345678",
		ReservationHour: "10",
	}
	if err := handler(ctx, cfg); err != nil {
		t.Fatal(err)
	}
}
