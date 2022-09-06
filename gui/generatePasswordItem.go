package gui

import (
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/Daniorocket/password-manager/crypt"
)

func generatePasswordItem() *fyne.MenuItem {
	generatePasswordItem := fyne.NewMenuItem("Generuj hasło", func() {
		wordsAlplabet := ""
		smallLettersAlphabet := ""
		bigLettersAlphabet := ""
		digitsAlphabet := ""
		specialCharactersAlphabet := ""
		label := widget.NewLabel("Wybierz zestawy znaków i naciśnij przycisk do generowania hasła. W rezultacie oprócz hasła otrzymasz jego siłę wyrażoną w bitach.")
		label2 := widget.NewLabel("Oto Twoje hasło:")
		password := widget.NewPasswordEntry()
		lenPassword := widget.NewEntry()
		lenPassword.SetPlaceHolder("Wprowadź liczbę całkowitą")
		entropyLabel := widget.NewLabel("")
		smallLetters := widget.NewCheck("[a-z]", func(b bool) {
			smallLettersAlphabet = crypt.PrepareAlphabet(b, smallLettersAlphabet, "abcdefghijklmnopqrstuvwxyz")
		})
		bigLetters := widget.NewCheck("[A-Z]", func(b bool) {
			bigLettersAlphabet = crypt.PrepareAlphabet(b, bigLettersAlphabet, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		})
		digits := widget.NewCheck("[0-9]", func(b bool) {
			digitsAlphabet = crypt.PrepareAlphabet(b, digitsAlphabet, "0123456789")
		})
		specialCharacters := widget.NewCheck("!@#", func(b bool) {
			specialCharactersAlphabet = crypt.PrepareAlphabet(b, specialCharactersAlphabet, "!@#$%^&*(){}[]\\|:\";'<>?,./")
		})
		btnToClipboard := widget.NewButton("Kopiuj", func() {
			clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
			clipboard.SetContent(password.Text)
		})
		boxwidget := widget.NewHBox(label2, password, btnToClipboard)
		form := &widget.Form{
			Items: []*widget.FormItem{
				widget.NewFormItem("Długość hasła:", lenPassword),
				widget.NewFormItem("", entropyLabel),
			},
			OnSubmit: func() {
				wordsAlplabet = smallLettersAlphabet + bigLettersAlphabet + digitsAlphabet + specialCharactersAlphabet
				lenPass, err := strconv.ParseInt(lenPassword.Text, 10, 64)
				if err != nil {
					alertMessage(Window.MessageWindow, "Należy wprowadzić długość hasła jako liczbę całkowitą.")
					return
				}
				if lenPass <= 0 || lenPass > 50 {
					alertMessage(Window.MessageWindow, "Długość hasła musi być dodatnią liczbą całkowitą nie większą od 50.")
					return
				}
				if len(wordsAlplabet) == 0 {
					alertMessage(Window.MessageWindow, "Należy zaznaczyć przynajmniej jeden alfabet znaków w generowanym haśle.")
					return
				}
				password.Text = crypt.GeneratePassword(lenPass, wordsAlplabet)
				if password.Text == "" {
					alertMessage(Window.MessageWindow, "Hasło zostało źle wygenerowane, proszę spróbować ponownie")
				}
				entropyLabel.Text = crypt.CalculateEntropy(password.Text, wordsAlplabet)
				boxwidget.Show()
				Window.MenagerWindow.Content().Refresh()
			},
			SubmitText: "Generuj hasło",
		}
		boxwidget.Hide()
		Window.MenagerWindow.SetContent(widget.NewVBox(
			label,
			widget.NewHBox(
				smallLetters,
				bigLetters,
				digits,
				specialCharacters,
			),
			form,
			boxwidget,
		))
	})
	return generatePasswordItem
}
