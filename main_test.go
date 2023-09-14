package main

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Main(t *testing.T) {
	ctx := context.Background()
	cfg := ReserveConfig{
		LastName:      "l",
		FirstName:     "f",
		LastNameKana:  "lk",
		FirstNameKana: "fl",
		MailAddress:   "sample@aa.jp",
		Tel:           "08012345678",
		Reservations: map[JWeekday]*ReservationDetail{
			Tue: {
				TrainerName: "A",
				Hour:        "14",
			},
		},
	}
	if err := handler(ctx, cfg); err != nil {
		t.Fatal(err)
	}
}

func Test_todayReservation(t *testing.T) {
	// todo timeをmockしたい
	tests := []struct {
		label                      string
		input                      map[JWeekday]*ReservationDetail
		wantDetail                 *ReservationDetail
		wantStartDate, wantEndDate string
	}{
		{
			label: "tmp",
			input: map[JWeekday]*ReservationDetail{
				Mon: {
					Hour: "09",
				},
				Tue: {
					Hour: "10",
				},
				Wed: {
					Hour: "14",
				},
			},
			wantDetail:    &ReservationDetail{},
			wantStartDate: "",
			wantEndDate:   "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.label, func(t *testing.T) {
			rd, sd, ed := todayReservation(tt.input)
			if d := cmp.Diff(rd, tt.wantDetail); d != "" {
				t.Fatal(d)
			}
			if d := cmp.Diff(sd, tt.wantStartDate); d != "" {
				t.Fatal(d)
			}
			if d := cmp.Diff(ed, tt.wantEndDate); d != "" {
				t.Fatal(d)
			}
		})
	}
}
