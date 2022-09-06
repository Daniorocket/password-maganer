package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func prepareMenagerWindow() {
	Window.MenagerWindow.SetOnClosed(func() {
		Window.MenagerWindow = fyne.CurrentApp().NewWindow("Menedżer haseł")
		Window.LoginWindow.Show()
	})
	logoutItem := fyne.NewMenuItem("Wyloguj", func() {
		Window.MenagerWindow.Hide()
		Window.LoginWindow.Content().Refresh()
		Window.LoginWindow.Show()

	})

	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu("Program", generatePasswordItem(), fyne.NewMenuItemSeparator(), logoutItem),
		fyne.NewMenu("Zarządzaj", vaultPasswordItem()),
	)
	Window.MenagerWindow.SetMainMenu(mainMenu)
	Window.MenagerWindow.Resize(fyne.Size{Width: 600, Height: 300})
	Window.MenagerWindow.Show()
	label := widget.NewLabel("Witaj użytkowniku! Korzystaj z menu do nawigacji w programie.")
	Window.MenagerWindow.SetContent(label)
}
