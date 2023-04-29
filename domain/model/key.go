package model

import "crypto/rsa"

type Key struct {
	ID        string
	PublicKey *rsa.PublicKey
}
