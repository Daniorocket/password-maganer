package gui

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/Daniorocket/password-manager/config"
	"github.com/Daniorocket/password-manager/crypt"
	"github.com/Daniorocket/password-manager/secondfa"
	"github.com/Daniorocket/password-manager/sqldb"
	"github.com/Daniorocket/password-manager/user"
	"golang.org/x/crypto/pbkdf2"
)

func checkFilesExists() error {
	var _, err = os.Stat(config.SaltFile)
	if os.IsNotExist(err) {
		return err
	}
	_, err = os.Stat(config.CipherFile)
	if os.IsNotExist(err) {
		return err
	}
	_, err = os.Stat(config.IvFile)
	if os.IsNotExist(err) {
		return err
	}
	return nil
}
func createFile(filepath string) (*os.File, error) {
	file, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func InitLoginGUI(app fyne.App, User *user.User) {
	entryPassword := widget.NewPasswordEntry()
	progress := widget.NewProgressBar()
	progress.Hidden = true
	label := widget.NewLabel("Jeśli to pierwsze logowanie to Twoje  hasło zostanie ustawione")
	entryPassword.PlaceHolder = "Wprowadź hasło główne do  programu"
	form := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("Hasło główne:", entryPassword),
		},

		OnSubmit: func() {
			if entryPassword.Text == "" {
				alertMessage(Window.MessageWindow, "Dane do logowania nie mogą być puste.")
				return
			}
			_, err := os.Stat("files")

			if os.IsNotExist(err) {
				time.Sleep(50 * time.Millisecond)
				errDir := os.MkdirAll("files", 0777)
				if errDir != nil {
					log.Fatal(err)
				}
			}
			if err := checkFilesExists(); err != nil {
				saltFile, err := createFile(config.SaltFile)
				if err != nil {
					alertMessage(Window.MessageWindow, "Nie udało się stworzyć pliku z danymi")
					return
				}
				salt := make([]byte, 32)
				_, err = io.ReadFull(rand.Reader, salt)
				if err != nil {

					return
				}
				if _, err := saltFile.Write(salt); err != nil {
					alertMessage(Window.MessageWindow, "Nie udało się zapisać do pliku z sola")
					return
				}
				defer saltFile.Close()
				iv := make([]byte, aes.BlockSize)
				if _, err = io.ReadFull(rand.Reader, iv); err != nil {
					alertMessage(Window.MessageWindow, "Nie udało się wygenerowac wektora IV")
					return
				}
				ivFile, err := createFile(config.IvFile)
				defer ivFile.Close()
				if err != nil {
					alertMessage(Window.MessageWindow, "Nie udało się stworzyć pliku z wektorem IV")
					return
				}
				if _, err := ivFile.Write(iv); err != nil {
					alertMessage(Window.MessageWindow, "Nie udało się zapisać do pliku z wektorem IV")
					return
				}
				keyMasterPassword := pbkdf2.Key([]byte(entryPassword.Text), salt, config.CountIterationPBKDF, 32, sha256.New)
				keyDatabase := pbkdf2.Key(keyMasterPassword, salt, config.CountIterationPBKDF, 32, sha256.New)
				encPass, err := crypt.AesEncrypt(keyMasterPassword, entryPassword.Text, iv)
				if err != nil {
					log.Println(err)
				}
				encFile, err := createFile(config.CipherFile)
				defer encFile.Close()
				if err != nil {
					alertMessage(Window.MessageWindow, "Nie udało się stworzyć pliku z szyfrogramem")
					return
				}
				if _, err := encFile.Write([]byte(encPass)); err != nil {
					alertMessage(Window.MessageWindow, "Nie udało się zapisać do pliku z szyfrogramem")
					return
				}
				go func() {
					progress.Hidden = false
					for i := 0.0; i <= 1.0; i += 0.01 {
						time.Sleep(time.Millisecond * 10)
						progress.SetValue(i)
					}
					err = sqldb.CreateDb(keyDatabase)
					if err != nil {
						alertMessage(Window.MessageWindow, "Nie udało się nawiązać połączenia z bazą danych")
						return
					}
					secret, imageByte, err := secondfa.GenerateQRAndSecret(config.EmailGoogleAuth)
					if err != nil {
						alertMessage(Window.MessageWindow, "Nie udało się wygenerować pliku QR")
						return
					}
					if err := sqldb.GoogleInsertRow(secret); err != nil {
						alertMessage(Window.MessageWindow, "Nie udało się zapisać sekretu do bazy")
						return
					}
					entryPassword.Text = ""
					Window.LoginWindow.Content().Refresh()
					Window.LoginWindow.Hide()
					prepareGoogleAuthWindow(true, imageByte)
					User.Key = keyDatabase
					User.IsLogged = true
					progress.Hidden = true
					progress.SetValue(0.0)
				}()
			} else {
				salt, err := ioutil.ReadFile(config.SaltFile)
				if err != nil {
					alertMessage(Window.MessageWindow, "Nie udało się otworzyć pliku z solą")
					return
				}
				iv, err := ioutil.ReadFile(config.IvFile)
				if err != nil {
					alertMessage(Window.MessageWindow, "Nie udało się otworzyć pliku z wektorem początkowym")
					return
				}
				encFile, err := ioutil.ReadFile(config.CipherFile)
				if err != nil {
					alertMessage(Window.MessageWindow, "Nie udało się otworzyć pliku z szyfrogramem")
					return
				}
				keyMasterPassword := pbkdf2.Key([]byte(entryPassword.Text), salt, config.CountIterationPBKDF, 32, sha256.New)
				keyDatabase := pbkdf2.Key(keyMasterPassword, salt, config.CountIterationPBKDF, 32, sha256.New)
				encmess, err := crypt.AesEncrypt(keyMasterPassword, entryPassword.Text, iv)
				if err != nil {
					alertMessage(Window.MessageWindow, "Nie udało się zaszyfrować danych")
					return
				}
				if res := bytes.Compare(encFile, []byte(encmess)); res == 0 {
					progress.Hidden = false
					go func() {

						err = sqldb.CreateDb(keyDatabase)
						if err != nil {
							alertMessage(Window.MessageWindow, "Nie udało się nawiązać połączenia z bazą danych")
							return
						}
						for i := 0.0; i <= 1.0; i += 0.01 {
							time.Sleep(time.Millisecond * 10)
							progress.SetValue(i)
						}
						Window.LoginWindow.Content().Refresh()
						Window.LoginWindow.Hide()
						User.Key = keyDatabase
						User.IsLogged = true
						prepareGoogleAuthWindow(true, nil)
						progress.Hidden = true
						progress.SetValue(0.0)
						entryPassword.Text = ""
					}()
				} else {
					//	fmt.Println("Logowanie niepomyślne!")
					entryPassword.Text = ""
					Window.LoginWindow.Content().Refresh()
					alertMessage(Window.MessageWindow, "Nie udało Ci się zalogować, spróbuj ponownie.")
				}
			}
		},
		SubmitText: "Zaloguj się",
	}
	Window.LoginWindow.SetContent(widget.NewVBox(form, label, progress))
	Window.LoginWindow.Resize(fyne.Size{Width: 300, Height: 150})
	Window.LoginWindow.ShowAndRun()
	Window.LoginWindow.SetOnClosed(func() {
		os.Exit(0)
	})
}
