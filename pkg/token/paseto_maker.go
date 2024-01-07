package token

import (
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

type PasetoMaker struct {
	paseto               paseto.Token
	v4SymmetricSecretKey paseto.V4SymmetricKey
}

func NewPasetoMaker(v4SymmetricSecretKeyHex string) (Maker, error) {
	if len(v4SymmetricSecretKeyHex) != 64 {
		return nil, fmt.Errorf("invalid symmetric key : must be exactly %d character", 64)
	}

	v4SymmetricSecretKey, err := paseto.V4SymmetricKeyFromHex(v4SymmetricSecretKeyHex)
	if err != nil {
		return nil, err
	}

	maker := &PasetoMaker{
		paseto:               paseto.NewToken(),
		v4SymmetricSecretKey: v4SymmetricSecretKey,
	}

	return maker, nil
}

func (maker *PasetoMaker) CreateToken(email string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(email, duration)
	if err != nil {
		return "", nil, err
	}

	fmt.Println("Login -- Expired at : ", payload.ExpiredAt)
	maker.paseto.SetJti(payload.Jti.String()) // tokenID
	maker.paseto.SetExpiration(payload.ExpiredAt)
	maker.paseto.SetIssuedAt(payload.IssuedAt)
	maker.paseto.SetIssuer(payload.Issuer)

	return maker.paseto.V4Encrypt(maker.v4SymmetricSecretKey, nil), payload, nil
}

func (maker *PasetoMaker) VerifyToken(tokenString string) (*Payload, error) {
	parser := paseto.NewParser()

	var token *paseto.Token
	token, err := parser.ParseV4Local(maker.v4SymmetricSecretKey, tokenString, nil)
	if err != nil {
		fmt.Println("parse token : ", err.Error())
		return nil, ErrInvalidToken
	}

	jti, err := token.GetJti()
	if err != nil {
		fmt.Println("Jti : ", err.Error())
		return nil, ErrInvalidToken
	}

	issuer, err := token.GetIssuer()
	if err != nil {
		fmt.Println("Issuer : ", err.Error())
		return nil, ErrInvalidToken
	}

	expiredAt, err := token.GetExpiration()
	if err != nil {
		fmt.Println("Expired at : ", err.Error())
		return nil, ErrInvalidToken
	}

	issuedAt, err := token.GetIssuedAt()
	if err != nil {
		fmt.Println("Issued at : ", err.Error())
		return nil, ErrInvalidToken
	}

	uuid, err := uuid.Parse(jti)
	if err != nil {
		fmt.Println("UUID : ", err.Error())
		return nil, ErrInvalidToken
	}

	payload := &Payload{
		Jti:       uuid,
		Issuer:    issuer,
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
