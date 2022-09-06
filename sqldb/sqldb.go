package sqldb

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Daniorocket/password-manager/config"
	_ "github.com/mutecomm/go-sqlcipher/v4"
)

type Secret struct {
	key              []byte
	secretGoogleAuth string
}

var secret Secret

type Credentials struct {
	ID       int
	Name     string
	Username string
	Password string
}

func CreateDb(key []byte) error {
	secret.key = key
	dbname := fmt.Sprintf("%s?_pragma_key=x'%s'&_pragma_cipher_page_size=4096", config.DbFilename, secret.key)
	fmt.Println(dbname)
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS credentials (id INTEGER PRIMARY KEY AUTOINCREMENT," +
		"name varchar(64) NULL,username varchar(64) NULL,password varchar(64) NULL)")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(); err != nil {
		return err
	}
	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS google (id INTEGER PRIMARY KEY AUTOINCREMENT," +
		"secret varchar(64) NULL)")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(); err != nil {
		return err
	}
	return nil
}
func InsertRow(name string, username string, password string) error {
	dbname := fmt.Sprintf("%s?_pragma_key=x'%s'&_pragma_cipher_page_size=4096", config.DbFilename, secret.key)
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO credentials(name, username,password) values(?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(name, username, password)
	if err != nil {
		return err
	}
	return nil

}
func GoogleInsertRow(sec string) error {
	dbname := fmt.Sprintf("%s?_pragma_key=x'%s'&_pragma_cipher_page_size=4096", config.DbFilename, secret.key)
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO google(secret) values(?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(sec)
	if err != nil {
		return err
	}
	return nil
}
func GetSecretGoogleAuth() (string, error) {
	var sec string
	dbname := fmt.Sprintf("%s?_pragma_key=x'%s'&_pragma_cipher_page_size=4096", config.DbFilename, secret.key)
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return "", err
	}
	defer db.Close()
	rows, err := db.Query("SELECT secret FROM google where id=1")
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&sec)
		if err != nil {
			return "", err
		}
	}
	return sec, nil
}
func UpdateRow(id int, name string, username string, password string) error {
	dbname := fmt.Sprintf("%s?_pragma_key=x'%s'&_pragma_cipher_page_size=4096", config.DbFilename, secret.key)
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare("update credentials set name=?, username=?, password=? where id=?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(name, username, password, id)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affect == 0 {
		return errors.New("Row not found.")
	}
	fmt.Println("Updated: ", affect, " record.")
	return nil

}
func SelectAll() ([]Credentials, error) {
	var credentials []Credentials
	dbname := fmt.Sprintf("%s?_pragma_key=x'%s'&_pragma_cipher_page_size=4096", config.DbFilename, secret.key)
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT id,name,username,password FROM credentials")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var credential Credentials
		err = rows.Scan(
			&credential.ID, &credential.Name, &credential.Username, &credential.Password)
		if err != nil {
			return nil, err
		}
		credentials = append(credentials, credential)
	}
	return credentials, nil
}
func DeleteRow(id int) error {
	dbname := fmt.Sprintf("%s?_pragma_key=x'%s'&_pragma_cipher_page_size=4096", config.DbFilename, secret.key)
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare("delete from credentials where id=?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affect == 0 {
		return errors.New("Row not found.")
	}
	fmt.Println("Deleted: ", affect, " record.")
	return nil
}
