package token

import (
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoMaker struct {
	paseto                paseto.Token
	v4AsymmetricSecretKey paseto.V4AsymmetricSecretKey
}

func NewPasetoMaker(v4AsymmetricSecretKeyHex string) (Maker, error) {
	if len(v4AsymmetricSecretKeyHex) != 128 {
		return nil, fmt.Errorf("invalid symmetric key : must be exactly %d character", chacha20poly1305.KeySize)
	}

	v4AsymmetricSecretKey, err := paseto.NewV4AsymmetricSecretKeyFromHex(v4AsymmetricSecretKeyHex)
	if err != nil {
		return nil, err
	}

	maker := &PasetoMaker{
		paseto:                paseto.NewToken(),
		v4AsymmetricSecretKey: v4AsymmetricSecretKey,
	}

	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	maker.paseto.SetJti(payload.ID.String()) // tokenID
	maker.paseto.SetExpiration(payload.ExpiredAt)
	maker.paseto.SetIssuedAt(payload.IssuedAt)
	maker.paseto.SetString("userID", username)

	return maker.paseto.V4Sign(maker.v4AsymmetricSecretKey, nil), payload, nil

}

func (maker *PasetoMaker) VerifyToken(tokenString string, v4AsymmetricPublicKeyHex string) (*Payload, error) {
	payload := &Payload{}
	publicKey, err := paseto.NewV4AsymmetricPublicKeyFromHex(v4AsymmetricPublicKeyHex)

	parser := paseto.NewParser()
	var token *paseto.Token
	token, err = parser.ParseV4Public(publicKey, tokenString, nil)
	if err != nil {
		fmt.Println(ErrInvalidToken, err)
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	fmt.Println(token.Claims())

	return payload, nil
}
