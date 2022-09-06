package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func alertMessage(window fyne.Window, mess string) {
	window.Show()
	label := widget.NewLabel(mess)
	button := widget.NewButton("OK", func() {
		window.Hide()
	})
	window.SetContent(widget.NewVBox(
		label,
		button,
	))
}
