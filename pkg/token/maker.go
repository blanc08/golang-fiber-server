package token

import "time"

type Maker interface {
	CreateToken(username string, duration time.Duration) (string, *Payload, error)

	VerifyToken(token string, v4AsymmetricPublicKeyHex string) (*Payload, error)
}
