package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/Daniorocket/password-manager/sqldb"
)

func DetailsPasswordWindow(window fyne.Window, messageWindow fyne.Window, menagerWindow fyne.Window, cred sqldb.Credentials) {
	window.Show()
	window.SetOnClosed(func() {
		Window.DetailsPasswordWindow = fyne.CurrentApp().NewWindow("Szczegóły wpisu")
	})
	nameEntry := widget.NewEntry()
	nameEntry.Text = cred.Name
	userEntry := widget.NewEntry()
	userEntry.Text = cred.Username
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.Text = cred.Password
	repeatPasswordEntry := widget.NewPasswordEntry()
	repeatPasswordEntry.PlaceHolder = "Powtórz hasło"
	btnDelete := widget.NewButton("Usuń wpis", func() {
		if err := confirmOperation(Window.ConfirmOperationWindow, "Na pewno chcesz usunąć rekord?"); err != nil {
			return
		}
		if err := sqldb.DeleteRow(cred.ID); err != nil {
			alertMessage(messageWindow, "Nie można usunąć rekordu.")
			return
		}
		alertMessage(messageWindow, "Pomyślnie usunięto wpis")
		Window.DetailsPasswordWindow.Hide()
		ListOfPasswords()
		menagerWindow.Content().Refresh()
	})
	btnCancel := widget.NewButton("Wróć", func() {
		window.Hide()
	})
	form := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("Nazwa witryny:", nameEntry),
			widget.NewFormItem("Użytkownik:", userEntry),
			widget.NewFormItem("Hasło:", passwordEntry),
			widget.NewFormItem("Powtórz hasło:", repeatPasswordEntry),
		},
		OnSubmit: func() {

			if err := validateForm(nameEntry.Text, userEntry.Text, passwordEntry.Text, repeatPasswordEntry.Text, messageWindow); err != nil {
				passwordEntry.Text = ""
				repeatPasswordEntry.Text = ""
				window.Content().Refresh()
				return
			}
			if err := sqldb.UpdateRow(cred.ID, nameEntry.Text, userEntry.Text, passwordEntry.Text); err != nil {
				alertMessage(messageWindow, "Nie można zaktualizować danych, wprowadzono niepoprawne informacje")
				return
			}
			alertMessage(messageWindow, "Pomyślnie zaktualizowano rekord")
			Window.DetailsPasswordWindow.Hide()
			ListOfPasswords()
			menagerWindow.Content().Refresh()
			repeatPasswordEntry.Text = ""
		},
		SubmitText: "Edytuj wpis",
	}
	window.SetContent(widget.NewVBox(
		form,
		btnDelete,
		btnCancel,
	))
}
