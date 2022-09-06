package secondfa

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"net/url"

	"github.com/Daniorocket/password-manager/config"
	"github.com/Daniorocket/password-manager/sqldb"
	dgoogauth "github.com/dgryski/dgoogauth"
	qr "rsc.io/qr"
)

var otpc *dgoogauth.OTPConfig

func GenerateQRAndSecret(email string) (string, []byte, error) {
	sec := make([]byte, 10)
	_, err := rand.Read(sec)
	if err != nil {
		return "", nil, err
	}
	secretBase32 := base32.StdEncoding.EncodeToString(sec)
	issuer := config.NameAppGoogleAuth
	URL, err := url.Parse("otpauth://totp")
	if err != nil {
		return "", nil, err
	}
	URL.Path += "/" + url.PathEscape(issuer) + ":" + url.PathEscape(email)
	params := url.Values{}
	params.Add("secret", secretBase32)
	params.Add("issuer", issuer)
	URL.RawQuery = params.Encode()
	code, err := qr.Encode(URL.String(), qr.Q)
	if err != nil {
		return "", nil, err
	}
	b := code.PNG()
	otpc = &dgoogauth.OTPConfig{
		Secret:      secretBase32,
		WindowSize:  3,
		HotpCounter: 0,
		// UTC:         true,
	}
	return secretBase32, b, nil
}
func AuthenticateByToken(token string) error {
	if otpc == nil {
		sec, err := sqldb.GetSecretGoogleAuth()
		if err != nil {
			return err
		}
		otpc = &dgoogauth.OTPConfig{
			Secret:      sec,
			WindowSize:  3,
			HotpCounter: 0,
		}
	}
	val, err := otpc.Authenticate(token)
	if err != nil {
		return err
	}

	if !val {
		return errors.New("Not authenticated")
	}
	return nil
}
