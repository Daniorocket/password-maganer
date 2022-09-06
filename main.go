package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"github.com/Daniorocket/password-manager/gui"
	"github.com/Daniorocket/password-manager/user"
)

func main() {
	User := user.InitUser()
	application := app.New()
	loginWindow := application.NewWindow("Logowanie")
	messageWindow := fyne.CurrentApp().NewWindow("Powiadomienie")
	menagerWindow := fyne.CurrentApp().NewWindow("Menedżer haseł")
	detailsPasswordWindow := fyne.CurrentApp().NewWindow("Szczegóły wpisu")
	googleAuthWindow := fyne.CurrentApp().NewWindow("Uwierzytelnianie Google Authenticator")
	confirmOperationWindow := fyne.CurrentApp().NewWindow("Potwierdź operację")
	gui.Window = gui.Windows{
		DetailsPasswordWindow:  detailsPasswordWindow,
		LoginWindow:            loginWindow,
		MenagerWindow:          menagerWindow,
		MessageWindow:          messageWindow,
		GoogleAuthWindow:       googleAuthWindow,
		ConfirmOperationWindow: confirmOperationWindow,
	}
	gui.InitLoginGUI(application, User)
}
