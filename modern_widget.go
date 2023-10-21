package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func NewModernUI() fyne.CanvasObject {
	modernUI := &ModernUI{
		Status:       true,
		ProfileImage: "github.png",
		ProfileName:  "GitHub",
		Type:         "Issue",
		Message:      "GitHub is how people build software. Millions of developers and companies build, ship, and maintain their software on GitHub—the largest and most advanced development platform in the world.",
		Time:         time.Now(),
		OpenCallback: func(btn *widget.Button) {
			fmt.Println("Open")
		},
		ReadCallback: func(btn *widget.Button) {
			fmt.Println("Read")
		},
	}

	modernUI.ExtendBaseWidget(modernUI)
	return modernUI
}

type ModernUI struct {
	widget.BaseWidget
	Status       bool
	ProfileImage string
	ProfileName  string
	Type         string
	Message      string
	Time         time.Time
	OpenCallback func(*widget.Button)
	ReadCallback func(*widget.Button)
}

func (m *ModernUI) SetStatus(status bool) {
	m.Status = status
}

func (m *ModernUI) SetProfileImage(profileImage string) {
	m.ProfileImage = profileImage
}

func (m *ModernUI) SetProfileName(profileName string) {
	m.ProfileName = profileName
}

func (m *ModernUI) SetType(t string) {
	m.Type = t
}

func (m *ModernUI) SetMessage(message string) {
	m.Message = message
}

func (m *ModernUI) SetTime(time time.Time) {
	m.Time = time
}

func (m *ModernUI) SetOpenCallback(openCallback func(*widget.Button)) {
	m.OpenCallback = openCallback
}

func (m *ModernUI) SetReadCallback(readCallback func(*widget.Button)) {
	m.ReadCallback = readCallback
}

func (m *ModernUI) MinSize() fyne.Size {
	return m.BaseWidget.MinSize()
}

func (m *ModernUI) CreateRenderer() fyne.WidgetRenderer {
	padding := theme.Padding()

	statusColor := fyne.CurrentApp().Settings().Theme().Color("StatusRead", theme.VariantLight)
	if !m.Status {
		statusColor = fyne.CurrentApp().Settings().Theme().Color("StatusUnread", theme.VariantLight)
	}

	status := canvas.NewCircle(statusColor)
	status.Resize(fyne.NewSize(8, 8))

	githubIcon := fyne.CurrentApp().Settings().Theme().Icon("GitHub")

	image := canvas.NewImageFromResource(githubIcon)
	image.FillMode = canvas.ImageFillContain
	image.Resize(fyne.NewSize(40, 40))

	name := canvas.NewText(m.ProfileName, theme.ForegroundColor())
	name.TextStyle.Bold = true
	name.Resize(name.MinSize())

	ntypeColor := fyne.CurrentApp().Settings().Theme().Color("NType", theme.VariantLight)
	ntype := canvas.NewText(m.Type, ntypeColor)
	ntype.Resize(ntype.MinSize())

	message := canvas.NewText(m.Message, theme.ForegroundColor())
	message.Resize(message.MinSize())

	timeColor := fyne.CurrentApp().Settings().Theme().Color("Time", theme.VariantLight)
	time := canvas.NewText(convertTimeToTimeAgo(m.Time), timeColor)
	time.Alignment = fyne.TextAlignTrailing
	time.TextStyle.Italic = true
	time.Resize(time.MinSize())

	readBtn := widget.NewButton("Read", nil)
	readBtn.OnTapped = func() {
		m.ReadCallback(readBtn)

		m.Status = true
		status.FillColor = fyne.CurrentApp().Settings().Theme().Color("StatusRead", theme.VariantLight)
		status.Refresh()
	}
	readBtn.Resize(fyne.NewSize(readBtn.MinSize().Width+padding, 7*padding))

	openBtn := widget.NewButton("Open", nil)
	openBtn.OnTapped = func() {
		m.OpenCallback(openBtn)
	}
	openBtn.Resize(fyne.NewSize(openBtn.MinSize().Width+padding, 7*padding))

	modernUIRendererObj := &modernUIRenderer{
		ModernUI: m,
		status:   status,
		image:    image,
		name:     name,
		ntype:    ntype,
		message:  message,
		time:     time,
		readBtn:  readBtn,
		openBtn:  openBtn,
	}

	return modernUIRendererObj
}

type modernUIRenderer struct {
	ModernUI *ModernUI
	status   *canvas.Circle
	image    *canvas.Image
	name     *canvas.Text
	ntype    *canvas.Text
	message  *canvas.Text
	time     *canvas.Text
	readBtn  *widget.Button
	openBtn  *widget.Button
}

func (m *modernUIRenderer) Destroy() {
}

func (m *modernUIRenderer) Layout(size fyne.Size) {
	m.Resize(size)
}

func (m *modernUIRenderer) MinSize() fyne.Size {
	padding := 2 * theme.Padding()

	var width float32 = 0
	var height float32 = 0

	width += 4 * padding

	width += m.status.Size().Width
	width += m.image.Size().Width
	width += m.name.Size().Width

	height += 3.5 * padding

	height += m.image.Size().Height
	height += m.message.Size().Height
	height += m.time.Size().Height

	return fyne.NewSize(width, height)
}

func (m *modernUIRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		m.status,
		m.image,
		m.name,
		m.ntype,
		m.message,
		m.time,
		m.readBtn,
		m.openBtn,
	}
}

func (m *modernUIRenderer) Refresh() {
	m.status.FillColor = fyne.CurrentApp().Settings().Theme().Color("StatusRead", theme.VariantLight)

	if !m.ModernUI.Status {
		m.status.FillColor = fyne.CurrentApp().Settings().Theme().Color("StatusUnread", theme.VariantLight)
	}

	m.status.Refresh()

	m.image.File = m.ModernUI.ProfileImage
	m.image.Refresh()

	m.name.Text = (m.ModernUI.ProfileName)
	m.name.Refresh()

	m.ntype.Text = (m.ModernUI.Type)
	m.ntype.Refresh()

	m.message.Text = (m.ModernUI.Message)
	m.message.Refresh()

	m.time.Text = (convertTimeToTimeAgo(m.ModernUI.Time))
	m.time.Refresh()
}

func (m *modernUIRenderer) Resize(size fyne.Size) {
	padding := 2 * theme.Padding()

	statusPosX := padding
	statusPosY := size.Height/2.0 - m.status.Size().Height/2.0

	m.status.Move(fyne.NewPos(statusPosX, statusPosY))

	imagePosX := float32(m.status.Position().X + m.status.Size().Width + padding)
	imagePosY := padding

	m.image.Move(fyne.NewPos(imagePosX, imagePosY))

	readBtnPosX := float32(size.Width - m.readBtn.Size().Width - padding)
	readBtnPosY := float32(m.image.Position().Y + m.image.Size().Height/2.0 - m.readBtn.Size().Height/2.0)

	m.readBtn.Move(fyne.NewPos(readBtnPosX, readBtnPosY))
	// m.readBtn.Resize(m.readBtn.Size())

	openBtnPosX := float32(m.readBtn.Position().X - m.openBtn.Size().Width - padding)
	openBtnPosY := float32(m.readBtn.Position().Y)

	m.openBtn.Move(fyne.NewPos(openBtnPosX, openBtnPosY))
	// m.openBtn.Resize(m.openBtn.Size())

	namePosX := float32(m.image.Position().X + m.image.Size().Width + padding)
	namePosY := float32(m.image.Position().Y + m.image.Size().Height/2.0 - m.name.Size().Height)

	m.name.Move(fyne.NewPos(namePosX, namePosY))
	m.name.Resize(fyne.NewSize(m.openBtn.Position().X-namePosX-padding, m.name.MinSize().Height))

	ntypePosX := float32(m.name.Position().X)
	ntypePosY := float32(m.image.Position().Y + m.image.Size().Height/2.0)

	m.ntype.Move(fyne.NewPos(ntypePosX, ntypePosY))
	m.ntype.Resize(fyne.NewSize(m.openBtn.Position().X-ntypePosX-padding, m.ntype.MinSize().Height))

	messagePosX := float32(m.image.Position().X)
	messagePosY := float32(m.image.Position().Y + m.image.Size().Height + padding)

	m.message.Move(fyne.NewPos(messagePosX, messagePosY))
	m.message.Resize(fyne.NewSize(size.Width-messagePosX-padding, m.message.MinSize().Height))

	timePosX := float32(m.message.Position().X)
	timePosY := float32(m.message.Position().Y + m.message.Size().Height + padding/2.0)

	m.time.Move(fyne.NewPos(timePosX, timePosY))
	m.time.Resize(fyne.NewSize(size.Width-timePosX-padding, m.time.MinSize().Height))
}

func convertTimeToTimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "Just now"
	}

	if diff < time.Hour {
		return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
	}

	if diff < time.Hour*24 {
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	}

	if diff < time.Hour*24*7 {
		return fmt.Sprintf("%d days ago", int(diff.Hours()/24))
	}

	return t.Format("02 Jan 2006")
}
