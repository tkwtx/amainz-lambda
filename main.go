package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/pkg/errors"
	"github.com/sclevine/agouti"
)

type ReserveConfig struct {
	LastName      string                          `json:"lastName"`
	FirstName     string                          `json:"firstName"`
	LastNameKana  string                          `json:"lastNameKana"`
	FirstNameKana string                          `json:"firstNameKana"`
	MailAddress   string                          `json:"mailAddress"`
	Tel           string                          `json:"tel"`
	Reservations  map[JWeekday]*ReservationDetail `json:"reservations"`
}

type JWeekday string

type ReservationDetail struct {
	TrainerName string `json:"trainerName"`
	Hour        string `json:"hour"`
}

const (
	Sun JWeekday = "日"
	Mon          = "月"
	Tue          = "火"
	Wed          = "水"
	Thu          = "木"
	Fri          = "金"
	Sat          = "土"
)

func adaptJWeekday(w time.Weekday) JWeekday {
	switch w {
	case 0:
		return Sun
	case 1:
		return Mon
	case 2:
		return Tue
	case 3:
		return Wed
	case 4:
		return Thu
	case 5:
		return Fri
	case 6:
		return Sat
	default:
		return Sun
	}
}

func handler(_ context.Context, cfg ReserveConfig) error {
	if err := os.Setenv("HOME", "/opt/"); err != nil {
		return err
	}

	opts := []agouti.Option{
		agouti.Browser("chrome"),
		agouti.ChromeOptions(
			"args", []string{
				"--headless",
				"--no-sandbox",
				"--disable-gpu",
				"--single-process",
				"--window-size=1920,1080",
			},
		),
		agouti.ChromeOptions(
			"binary", "/opt/headless-chromium",
		),
	}
	driver := agouti.NewWebDriver(
		"http://{{.Address}}",
		[]string{"/opt/chromedriver", "--port={{.Port}}"},
		opts...)
	if err := driver.Start(); err != nil {
		return err
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		return err
	}

	rd, d := todayReservation(cfg.Reservations)
	if err := page.Navigate(fmt.Sprintf(`https://airrsv.net/amainz-kitasenju/calendar/menuDetail/?schdlId=s00004C101&bookingDate=%s`, d)); err != nil {
		return handleErrorWithScreenShot(page, err)
	}
	if err := page.FindByName("resrcSchdlItemId").Select(rd.TrainerName); err != nil {
		return handleErrorWithScreenShot(page, err)
	}
	if err := page.FindByButton("予約する").Click(); err != nil {
		return handleErrorWithScreenShot(page, err)
	}
	if err := page.Find(`input[name="lastNm"]`).Fill(cfg.LastName); err != nil {
		return handleErrorWithScreenShot(page, err)
	}
	if err := page.Find(`input[name="firstNm"]`).Fill(cfg.FirstName); err != nil {
		return handleErrorWithScreenShot(page, err)
	}
	if err := page.Find(`input[name="lastNmKn"]`).Fill(cfg.LastNameKana); err != nil {
		return handleErrorWithScreenShot(page, err)
	}
	if err := page.Find(`input[name="firstNmKn"]`).Fill(cfg.FirstNameKana); err != nil {
		return handleErrorWithScreenShot(page, err)
	}
	if err := page.Find(`input[name="mailAddress1"]`).Fill(cfg.MailAddress); err != nil {
		return handleErrorWithScreenShot(page, err)
	}
	if err := page.Find(`input[name="mailAddress1ForCnfrm"]`).Fill(cfg.MailAddress); err != nil {
		return handleErrorWithScreenShot(page, err)
	}
	if err := page.Find(`input[name="tel1"]`).Fill(cfg.Tel); err != nil {
		return handleErrorWithScreenShot(page, err)
	}
	if err := page.FindByButton("確認へ進む").Click(); err != nil {
		return handleErrorWithScreenShot(page, err)
	}
	if err := page.FindByButton("上記に同意して予約を確定する").Click(); err != nil {
		return handleErrorWithScreenShot(page, err)
	}
	return nil
}

func todayReservation(rs map[JWeekday]*ReservationDetail) (*ReservationDetail, string) {
	_, _ = time.LoadLocation("Asia/Tokyo")
	t := time.Now()
	adt := t.AddDate(0, 0, 20)
	rd := rs[adaptJWeekday(adt.Weekday())]
	adaptMonth := func(m time.Month) string {
		// 1月なら01に戻す
		switch m {
		case time.November, time.December:
			return fmt.Sprintf("%d", m)
		default:
			return fmt.Sprintf("0%d", m)
		}
	}
	return rd, fmt.Sprintf("%d%s%d%s0000", adt.Year(), adaptMonth(adt.Month()), adt.Day(), rd.Hour)
}

// screenShot デバッグ用
func handleErrorWithScreenShot(page *agouti.Page, err error) error {
	page.Screenshot("/tmp/test.png")
	buf, _ := os.ReadFile("/tmp/test.png")
	fmt.Printf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(buf))
	return err
}

func main() {
	lambda.Start(handler)
}
