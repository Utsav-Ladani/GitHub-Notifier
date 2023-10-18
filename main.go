package main

import (
	"context"
	"fmt"
	"image/color"
	"log"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/electricbubble/go-toast"
	"github.com/google/go-github/v55/github"
)

const APP_ID string = "org.mygithub.notification"
const REPEAT_TIME time.Duration = time.Second * 30

var notifierApp fyne.App
var window fyne.Window
var globalCtx context.Context
var ctxMap map[string]*context.CancelFunc
var notificationList []*github.Notification
var notificationListComponent *widget.List

// theme colors
var (
	ghGreenColor color.NRGBA = color.NRGBA{0x2c, 0xa4, 0x2c, 0xff}
	bgItemColor  color.NRGBA = color.NRGBA{0xd0, 0xd0, 0xd0, 0x20}
)

func main() {
	notifierApp = app.NewWithID(APP_ID)

	window = notifierApp.NewWindow("Github Notifications")

	notificationListComponent = addNotificationListUI()

	windowContentRefresh("Loading...")

	window.Resize(fyne.NewSize(400, 600))

	globalCtx = context.Background()
	ctxMap = make(map[string]*context.CancelFunc)

	github_token := notifierApp.Preferences().String("github_token")

	if github_token == "" {
		openSettingsPanel()
	} else {
		startAsyncProcess("githubNotifyLoop", githubNotifyLoop)
	}

	window.SetFixedSize(true)
	window.SetMaster()

	addSystemStrayMenu()
	window.SetCloseIntercept(func() {
		window.Hide()
	})

	window.ShowAndRun()
}

func fetchNotifications(ctx context.Context) ([]*github.Notification, error) {
	github_token := notifierApp.Preferences().String("github_token")

	client := github.NewClient(nil).WithAuthToken(github_token)

	opt := &github.NotificationListOptions{
		Since: time.Now().AddDate(0, 0, -10),
	}

	ctxTimeOut, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	notifications, _, err := client.Activity.ListNotifications(ctxTimeOut, opt)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return notifications, nil
}

func markAsReadNotification(notification *github.Notification) (bool, error) {
	github_token := notifierApp.Preferences().String("github_token")

	client := github.NewClient(nil).WithAuthToken(github_token)

	ctxTimeOut, cancel := context.WithTimeout(globalCtx, time.Second*10)
	defer cancel()

	_, err := client.Activity.MarkThreadRead(ctxTimeOut, notification.GetID())

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	log.Println("Mark as read success")

	return true, nil
}

func addNotifications(notifications []*github.Notification, err error) {
	log.Println("Add notifications")
	if err != nil {
		notificationList = nil
		log.Println(err)
		windowContentRefresh("Failed to fetch notifications")
		return
	}

	notificationsDiff := getNotificationListDiff(notifications)

	notificationList = notifications

	windowContentRefresh("No New Notifications")

	if len(notificationsDiff) != 0 {
		_ = toast.Push("Github Notifications",
			toast.WithTitle(fmt.Sprintf("You have %d new notifications", len(notificationsDiff))),
			toast.WithObjectiveC(true),
		)
	}
}

func getNotificationListDiff(notifications []*github.Notification) []*github.Notification {
	var diff []*github.Notification

	for _, notification := range notifications {
		if !isNotificationExist(notification) {
			diff = append(diff, notification)
		}
	}

	return diff
}

func isNotificationExist(notification *github.Notification) bool {
	for _, n := range notificationList {
		if n.GetID() == notification.GetID() {
			return true
		}
	}

	return false
}

func waitForProcess(ch *chan int) {
	<-*ch
}

func processEnd(ch *chan int) {
	*ch <- 1
}

func githubNotify(ch *chan int, ctx context.Context, callback func([]*github.Notification, error)) {
	defer processEnd(ch)

	notifications, err := fetchNotifications(ctx)

	select {
	case <-ctx.Done():
		log.Println("Context canceled")
		return
	default:
		callback(notifications, err)
	}

	// showNotifications(notifications)
}

func githubNotifyLoop(ctx context.Context) {
	ch := make(chan int)

	log.Println("Start github notification loop")

	for {
		go githubNotify(&ch, ctx, addNotifications)

		waitForProcess(&ch)

		// check if context is canceled
		select {
		case <-ctx.Done():
			log.Println("Context canceled")
			return
		default:
			log.Println("Wait for next loop")
			time.Sleep(REPEAT_TIME)
		}
	}
}

func openSettingsPanel() {
	githubTokenEntry := widget.NewEntry()
	githubTokenEntry.SetPlaceHolder("Enter Github Token")
	githubTokenEntry.SetText(notifierApp.Preferences().String("github_token"))

	spacer := canvas.NewRectangle(color.NRGBA{0x00, 0x00, 0x00, 0x00})
	spacer.SetMinSize(fyne.NewSize(0, 10))

	dialog := dialog.NewForm(
		"App Settings",
		"Save",
		"Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Token", githubTokenEntry),
			widget.NewFormItem("", spacer),
		},
		func(isSave bool) {
			old_github_token := notifierApp.Preferences().String("github_token")

			if isSave {
				notifierApp.Preferences().SetString("github_token", githubTokenEntry.Text)
			}

			new_github_token := notifierApp.Preferences().String("github_token")

			if new_github_token == "" {
				notifierApp.Quit()
			}

			if old_github_token != new_github_token {
				startAsyncProcess("githubNotifyLoop", githubNotifyLoop)
			}
		},
		window,
	)

	dialog.Resize(fyne.NewSize(300, 150))
	dialog.Show()
}

func addNotificationListUI() *widget.List {
	return widget.NewList(
		func() int {
			return len(notificationList)
		},
		func() fyne.CanvasObject {
			title := canvas.NewText("Title", color.White)
			title.TextStyle.Bold = true
			title.TextSize = 16

			// github green color
			title.Color = ghGreenColor

			content := widget.NewRichTextWithText("Content")
			content.Truncation = fyne.TextTruncateEllipsis

			markAsReadButton := widget.NewButton("Mark As Read", func() {
				log.Println("Mark As Read")
			})
			markAsReadButton.Importance = widget.HighImportance

			openInGithubButton := widget.NewButton("Open In Github", func() {
				log.Println("Open In Github")
			})
			openInGithubButton.Importance = widget.SuccessImportance

			buttonContainer := container.NewHBox(markAsReadButton, openInGithubButton)

			iteamContainer := container.NewPadded(container.NewVBox(title, content, buttonContainer))

			bgRect := canvas.NewRectangle(bgItemColor)

			return container.NewStack(bgRect, iteamContainer)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			title := notificationList[id].GetRepository().GetFullName()
			content := notificationList[id].GetSubject().GetTitle()

			iteamContainer := item.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*fyne.Container)

			titleCompoenent := iteamContainer.Objects[0].(*canvas.Text)
			contentComponent := iteamContainer.Objects[1].(*widget.RichText)

			titleCompoenent.Text = title
			contentComponent.ParseMarkdown(content)

			buttonContainer := iteamContainer.Objects[2].(*fyne.Container)

			markAsReadButton := buttonContainer.Objects[0].(*widget.Button)
			openInGithubButton := buttonContainer.Objects[1].(*widget.Button)

			markAsReadButton.OnTapped = func() {
				markAsReadButton.Disable()

				isRead, _ := markAsReadNotification(notificationList[id])

				if isRead {
					markAsReadButton.Hide()
					startAsyncProcess("githubNotifyLoop", githubNotifyLoop)
				}

				markAsReadButton.Enable()
			}

			openInGithubButton.OnTapped = func() {
				openInGithubButton.Disable()

				if err := exec.Command("open", "https://github.com/notifications").Start(); err != nil {
					log.Println(err)
				}

				openInGithubButton.Enable()
			}
		},
	)
}

func addToolbarUI() fyne.CanvasObject {
	preference := widget.NewToolbarAction(
		resourceSettingsPng,
		func() {
			openSettingsPanel()
		},
	)

	toolbar := widget.NewToolbar(preference)
	toolbar.Resize(fyne.NewSize(400, 50))

	return toolbar
}

func wrapperContainer(list *widget.List, altMessage string) fyne.CanvasObject {
	if list.Length() == 0 {
		label := widget.NewLabel(altMessage)
		label.Alignment = fyne.TextAlignCenter
		label.TextStyle.Bold = true

		return container.NewCenter(label)
	}

	return list
}

func windowContentRefresh(altMessage string) {
	notificationListComponent = addNotificationListUI()

	mainContainer := container.NewBorder(
		addToolbarUI(),
		nil,
		nil,
		nil,
		wrapperContainer(notificationListComponent, altMessage),
	)

	window.SetContent(mainContainer)
}

func addSystemStrayMenu() {
	menu := fyne.NewMenu("GitHub Notify",
		fyne.NewMenuItem("Show", func() {
			window.Show()
		}),
		fyne.NewMenuItem("Quit", func() {
			notifierApp.Quit()
		}),
	)

	if desk, ok := notifierApp.(desktop.App); ok {
		desk.SetSystemTrayIcon(resourceIconPng)
		desk.SetSystemTrayMenu(menu)
	}
}

func startAsyncProcess(name string, process func(ctx context.Context)) {
	if ctxCancelFunc, ok := ctxMap[name]; ok {
		log.Println("Cancel process", name)
		(*ctxCancelFunc)()
	}

	ctx, cancel := context.WithCancel(globalCtx)

	ctxMap[name] = &cancel

	go process(ctx)
}
