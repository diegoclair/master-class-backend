package token

import (
	"fmt"
	"log"
	"time"

	"github.com/o1egl/paseto"

	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

func (t *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	return t.paseto.Encrypt(t.symmetricKey, payload, nil)
}

func (t *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := t.paseto.Decrypt(token, t.symmetricKey, payload, nil)
	if err != nil {
		log.Println("error to decrypt token: ", err)
		return nil, ErrInvalidToken
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
