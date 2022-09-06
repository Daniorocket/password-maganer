# password-manager

This app has been created as two-factor password manager.   




## Author

- [@daniorocket](https://www.github.com/Daniorocket)

  
## Tech Stack

- **Language** Go 1.17.11
- **Database** SQLite
- **GUI** fyne.io
  
## Features

- GUI by fyne.io library
- Two-factor authentication to the app by master password and by token from Google Authenticator
- Hashed database SQLite
- Implemented AES encrypt and decrypt alghorithm with PBKDF
- Generation new secure passwords, alphabet choosen by user



  
## Run Locally

Clone the project

```bash
  git clone https://github.com/Daniorocket/password-manager.git
```

Go to the project directory

```bash
  cd path-to-my-project
```

Compile the app

```bash
  go run .
```


  
## Environment Variables

To run this project, you will need to add the following constants in configFile.

`DbFilename` - Defines path to the encrypted database.

`QrFilename` -  Defines path to the img with qr image used in GoogleAuthenticator app.

`SaltFile` - Defines path to the salt file.

`IvFile` - Defines path to the IV vector.

`CipherFile` - Defines path to the encrypted master password.

`EmailGoogleAuth` - Defines email used in generate QR code in GoogleAuthenticator.

`CountIterationPBKDF` - Defines number of iterations PBKDF alghorithm.
  
`NameAppGoogleAuth` - Defines name of the app Google Authenticator.