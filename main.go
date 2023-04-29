package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sclevine/agouti"
)

type ReserveConfig struct {
	LastName        string `json:"lastName"`
	FirstName       string `json:"firstName"`
	LastNameKana    string `json:"lastNameKana"`
	FirstNameKana   string `json:"firstNameKana"`
	MailAddress     string `json:"mailAddress"`
	Tel             string `json:"tel"`
	ReservationHour string `json:"reservationHour"`
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
	if err := page.Navigate(fmt.Sprintf(`https://airrsv.net/amainz-kitasenju/calendar/menuDetail/?schdlId=s00004C101&bookingDate=%s`, calDate(cfg.ReservationHour))); err != nil {
		return err
	}
	if err := page.FindByName("resrcSchdlItemId").Select("A"); err != nil {
		return err
	}
	if err := page.FindByButton("予約する").Click(); err != nil {
		return err
	}
	if err := page.Find(`input[name="lastNm"]`).Fill(cfg.LastName); err != nil {
		return err
	}
	if err := page.Find(`input[name="firstNm"]`).Fill(cfg.FirstName); err != nil {
		return err
	}
	if err := page.Find(`input[name="lastNmKn"]`).Fill(cfg.LastNameKana); err != nil {
		return err
	}
	if err := page.Find(`input[name="firstNmKn"]`).Fill(cfg.FirstNameKana); err != nil {
		return err
	}
	if err := page.Find(`input[name="mailAddress1"]`).Fill(cfg.MailAddress); err != nil {
		return err
	}
	if err := page.Find(`input[name="mailAddress1ForCnfrm"]`).Fill(cfg.MailAddress); err != nil {
		return err
	}
	if err := page.Find(`input[name="tel1"]`).Fill(cfg.Tel); err != nil {
		return err
	}
	if err := page.FindByButton("確認へ進む").Click(); err != nil {
		return err
	}
	if err := page.FindByButton("上記に同意して予約を確定する").Click(); err != nil {
		return err
	}
	return nil
}

func calDate(h string) string {
	_, _ = time.LoadLocation("Asia/Tokyo")
	t := time.Now()
	dt := t.AddDate(0, 0, 20)
	fixedMonth := func(m time.Month) string {
		// 1月なら01に戻す
		switch m {
		case time.November, time.December:
			return fmt.Sprintf("%d", m)
		default:
			return fmt.Sprintf("0%d", m)
		}
	}
	return fmt.Sprintf("%d%s%d%s0000", dt.Year(), fixedMonth(dt.Month()), dt.Day(), h)
}

// screenShot デバッグ用
//func screenShot(page *agouti.Page) (string, error) {
//	if err := page.Screenshot("/tmp/test.png"); err != nil {
//		return "", err
//	}
//
//	buf, err := os.ReadFile("/tmp/test.png")
//	if err != nil {
//		return "", err
//	}
//
//	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf), nil
//}

func main() {
	lambda.Start(handler)
}
