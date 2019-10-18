package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

type RepeatType int

const (
	EveryWeek    RepeatType = iota // 毎週
	Biweekly_1_3                   // 第1、第3
	Biweekly_2_4                   // 第2、第4

	Resources     string = "資源ごみ"
	Combustible   string = "可燃ごみ"
	InCombustible string = "不燃ごみ"
	PetBottol     string = "ペットボトル"
)

var (
	debug = flag.Bool("d", false, "Console debug mode")
)

type SlackPayload struct {
	Channel  string `json:"channel"`
	UserName string `json:"username"`
	Text     string `json:"text"`
}

var (
	d = []struct {
		Weekday    time.Weekday
		RepeatType RepeatType
		Value      string
	}{
		{
			Weekday:    1,
			RepeatType: EveryWeek,
			Value:      Combustible,
		},
		{
			Weekday:    3,
			RepeatType: EveryWeek,
			Value:      Resources,
		},
		{
			Weekday:    4,
			RepeatType: EveryWeek,
			Value:      Combustible,
		},
		{
			Weekday:    5,
			RepeatType: Biweekly_1_3,
			Value:      InCombustible,
		},
		{
			Weekday:    5,
			RepeatType: Biweekly_2_4,
			Value:      PetBottol,
		},
	}
)

func Handler(ctx context.Context) error {
	now := time.Now().In(getJST())
	tommorrow := now.Add(24 * time.Hour)
	w := getWeekCount(&tommorrow)

	var msg string
	for _, v := range d {
		if v.Weekday != tommorrow.Weekday() {
			continue
		}
		fmt.Println(tommorrow)
		fmt.Printf("WeekDay: %d\n", int(tommorrow.Weekday()))

		if v.RepeatType == Biweekly_1_3 && (w == 1 || w == 3) {
			fmt.Println("1_3")
			msg = fmt.Sprintf(os.Getenv("MESSAGE_TEMPLATE"), v.Value)
			break
		} else if v.RepeatType == Biweekly_2_4 && (w == 2 || w == 4) {
			fmt.Println("2_4")
			msg = fmt.Sprintf(os.Getenv("MESSAGE_TEMPLATE"), v.Value)
			break
		} else if v.RepeatType == EveryWeek {
			fmt.Println("every")
			msg = fmt.Sprintf(os.Getenv("MESSAGE_TEMPLATE"), v.Value)
			break
		}
	}

	if len(msg) == 0 {
		fmt.Println("no value")
		return nil
	}

	payload := &SlackPayload{
		Channel:  os.Getenv("SLACK_CHANNEL_ID"),
		UserName: os.Getenv("SLACK_BOT_NAME"),
		Text:     msg,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, os.Getenv("SLACK_HOOK_URL"), bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func getWeekCount(t *time.Time) int {
	_, w := t.ISOWeek()
	_, firstW := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).ISOWeek()
	return w - firstW + 1
}

func getJST() *time.Location {
	tz := time.FixedZone("JST", 9*60*60)
	return tz
}

func postSlack() error {
	return nil
}

func main() {
	flag.Parse()
	if !*debug {
		lambda.Start(Handler)
	}
	if err := Handler(context.Background()); err != nil {
		panic(err)
	}
}
