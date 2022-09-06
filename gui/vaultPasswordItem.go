package gui

import (
	"errors"
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/Daniorocket/password-manager/sqldb"
)

func validateForm(name string, user string, password string, repeatedPassword string, window fyne.Window) error {
	if name == "" || user == "" || password == "" || repeatedPassword == "" {
		alertMessage(window, "Wszystkie pola są wymagane")
		return errors.New("Wszystkie pola są wymagane")
	}
	if repeatedPassword != password {
		alertMessage(window, "Podane hasła muszą być identyczne")
		return errors.New("Podane hasła muszą być identyczne")
	}
	return nil
}
func ListOfPasswords() {
	credentials, err := sqldb.SelectAll()
	if err != nil {
		alertMessage(Window.MessageWindow, "Nie udało się pomyślnie połączyć programu z bazą danych")
		return
	}
	fmt.Println(credentials)
	box := widget.NewVBox()
	if len(credentials) == 0 {
		box.Append(widget.NewLabel("Brak wpisów w bazie danych, użyj menu do dodania nowych."))
		Window.MenagerWindow.SetContent(box)
		return
	}
	box.Append(widget.NewLabel("Nazwa rekordu:"))
	for i := 0; i < len(credentials); i++ {
		nameLabel := widget.NewLabel(credentials[i].Name)
		nameLabel.Alignment = fyne.TextAlignCenter
		j := i
		btnEntry := widget.NewButton("Zarządzaj", func() {
			DetailsPasswordWindow(Window.DetailsPasswordWindow, Window.MessageWindow, Window.MenagerWindow, credentials[j])
		})
		box.Append(widget.NewHBox(
			nameLabel,
			btnEntry,
		),
		)
	}
	scroll := widget.NewScrollContainer(box)
	Window.MenagerWindow.SetContent(scroll)
}
func vaultPasswordItem() *fyne.MenuItem {
	vaultPasswordItem := fyne.NewMenuItem("Twoje wpisy:", nil)
	vaultPasswordItem.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("Wyświetl listę", ListOfPasswords),
		fyne.NewMenuItem("Dodaj nowy", func() {
			nameEntry := widget.NewEntry()
			nameEntry.PlaceHolder = "Podaj adres witryny"
			userEntry := widget.NewEntry()
			userEntry.PlaceHolder = "Podaj nazwę użytkownika"
			passwordEntry := widget.NewPasswordEntry()
			passwordEntry.PlaceHolder = "Podaj hasło"
			repeatPasswordEntry := widget.NewPasswordEntry()
			repeatPasswordEntry.PlaceHolder = "Powtórz hasło"
			form := &widget.Form{
				Items: []*widget.FormItem{
					widget.NewFormItem("Nazwa witryny:", nameEntry),
					widget.NewFormItem("Użytkownik:", userEntry),
					widget.NewFormItem("Hasło:", passwordEntry),
					widget.NewFormItem("Powtórz hasło:", repeatPasswordEntry),
				},
				OnSubmit: func() {

					if err := validateForm(nameEntry.Text, userEntry.Text, passwordEntry.Text, repeatPasswordEntry.Text, Window.MessageWindow); err != nil {
						nameEntry.Text = ""
						userEntry.Text = ""
						passwordEntry.Text = ""
						repeatPasswordEntry.Text = ""
						Window.MenagerWindow.Content().Refresh()
						return
					}
					if err := sqldb.InsertRow(nameEntry.Text, userEntry.Text, passwordEntry.Text); err != nil {
						nameEntry.Text = ""
						userEntry.Text = ""
						passwordEntry.Text = ""
						repeatPasswordEntry.Text = ""
						alertMessage(Window.MessageWindow, "Nie udało się pomyślnie połączyć programu z bazą danych")
						return
					}
					nameEntry.Text = ""
					userEntry.Text = ""
					passwordEntry.Text = ""
					repeatPasswordEntry.Text = ""
					alertMessage(Window.MessageWindow, "Dodano pomyślnie wpis")
					Window.DetailsPasswordWindow.Hide()
					ListOfPasswords()
					Window.MenagerWindow.Content().Refresh()
				},
				SubmitText: "Dodaj nowy wpis",
			}
			Window.MenagerWindow.SetContent(widget.NewVBox(
				form,
			))
		}),
	)
	return vaultPasswordItem
}
