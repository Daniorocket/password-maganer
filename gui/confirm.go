package gui

import (
	"errors"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func confirmOperation(window fyne.Window, mess string) error {
	window.Show()
	var err error
	err = nil
	label := widget.NewLabel(mess)
	btnYes := widget.NewButton("TAK", func() {
		window.Hide()
	})
	btnNo := widget.NewButton("NIE", func() {
		window.Hide()
		err = errors.New("Uncaccept")
	})
	window.SetContent(widget.NewVBox(
		label,
		widget.NewHBox(
			btnYes,
			btnNo,
		),
	))
	return err
}
