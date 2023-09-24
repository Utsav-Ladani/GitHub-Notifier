package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"fyne.io/fyne/app"
	"github.com/google/go-github/v55/github"
	"github.com/joho/godotenv"
)

type AppState struct {
	Last_Updated time.Time
}

var appState AppState = AppState{}

const REPEAT_TIME time.Duration = time.Second * 2

func main() {
	loadEnv()

	appID := os.Getenv("APP_ID")

	app := app.NewWithID(appID)

	window := app.NewWindow("Github Notifications")

	window.ShowAndRun()
}

func loadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Failed to load env variables")
		log.Fatal(err)
		return
	}
}

func getState() *AppState {
	if appState.Last_Updated.IsZero() {
		appState.Last_Updated = time.Now().AddDate(0, 0, -5)
	}

	return &appState
}

func setState() {
	appState.Last_Updated = time.Now()
}

func fetchNotifications() []*github.Notification {
	client := github.NewClient(nil).WithAuthToken(os.Getenv("GITHUB_TOKEN"))

	appState := getState()

	opt := &github.NotificationListOptions{
		Since: appState.Last_Updated,
	}

	notifications, _, err := client.Activity.ListNotifications(context.TODO(), opt)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	setState()

	return notifications
}

func showNotification(notification *github.Notification) {
	fmt.Println("--------------------")
	fmt.Println(notification.GetRepository().GetName())
	fmt.Println(notification.GetReason())
	fmt.Println(notification.GetSubject().GetType())
	fmt.Println("--------------------")
}

func showNotifications(notifications []*github.Notification) {
	for _, notification := range notifications {
		showNotification(notification)
	}
}

func waitForProcess(ch *chan int) {
	<-*ch
}

func processEnd(ch *chan int) {
	*ch <- 1
}

func githubNotify(ch *chan int) {
	defer processEnd(ch)

	notifications := fetchNotifications()

	showNotifications(notifications)
}

func githubNotifyLoop() {
	ch := make(chan int)

	log.Println("Start github notification loop")

	for {
		go githubNotify(&ch)

		waitForProcess(&ch)

		log.Println("Wait for next loop")

		time.Sleep(REPEAT_TIME)
	}
}
