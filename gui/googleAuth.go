package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/Daniorocket/password-manager/secondfa"
)

func prepareGoogleAuthWindow(firstLaunch bool, imageByte []byte) {
	Window.GoogleAuthWindow.SetOnClosed(func() {
		Window.GoogleAuthWindow = fyne.CurrentApp().NewWindow("Uwierzytelnianie Google Authenticator")
		Window.LoginWindow.Show()
	})
	img := fyne.NewStaticResource("QR", imageByte)

	img2 := canvas.NewImageFromResource(img)
	img2.FillMode = canvas.ImageFillOriginal
	entryToken := widget.NewEntry()
	form := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("Token:", entryToken),
		},
		OnSubmit: func() {
			if err := secondfa.AuthenticateByToken(entryToken.Text); err != nil {
				alertMessage(Window.MessageWindow, "Autoryzacja niepomyślna. Proszę spróbować ponownie.")
				return
			}
			Window.MenagerWindow.Show()
			prepareMenagerWindow()
			Window.GoogleAuthWindow.Hide()
		},
		SubmitText: "Zaloguj się",
	}
	var box *fyne.Container
	if firstLaunch == true {
		entryToken.SetPlaceHolder("Zeskanuj zdjęcie i wprowadź token z aplikacji")
		box = fyne.NewContainerWithLayout(layout.NewVBoxLayout(), img2, form)
	} else {
		entryToken.SetPlaceHolder("Wprowadź token z aplikacji")
		box = fyne.NewContainerWithLayout(layout.NewVBoxLayout(), form)
	}
	Window.GoogleAuthWindow.SetContent(box)
	Window.GoogleAuthWindow.Show()
}
