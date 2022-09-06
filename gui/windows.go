package gui

import "fyne.io/fyne"

type Windows struct {
	MessageWindow          fyne.Window
	DetailsPasswordWindow  fyne.Window
	MenagerWindow          fyne.Window
	LoginWindow            fyne.Window
	GoogleAuthWindow       fyne.Window
	ConfirmOperationWindow fyne.Window
}

var Window Windows
