package keycache

import (
	"crypto/ed25519"
	"encoding/base64"
	"log"

	"github.com/devlup-labs/Libr/core/crypto/cryptoutils"
	"github.com/devlup-labs/Libr/core/mod_client/logger"
	"github.com/libp2p/go-libp2p/core/crypto"
)

var (
	PubKey  ed25519.PublicKey
	PrivKey ed25519.PrivateKey
)

func InitKeys() {
	var err error
	PubKey, PrivKey, err = cryptoutils.LoadKeys()
	if err != nil {
		log.Fatalf("Failed to load keys: %v", err)
	}
	logger.LogToFile("private key: " + base64.StdEncoding.EncodeToString(PrivKey))
	logger.LogToFile("public key: " + base64.StdEncoding.EncodeToString(PubKey))
}

func InitTempKeys() {
	var err error
	PubKey, PrivKey, err = cryptoutils.LoadTempKeys()
	if err != nil {
		log.Fatalf("Failed to load keys: %v", err)
	}
	logger.LogToFile("private key: " + base64.StdEncoding.EncodeToString(PrivKey))
	logger.LogToFile("public key: " + base64.StdEncoding.EncodeToString(PubKey))
}

func LoadPubKey() string {
	pub, _, _ := cryptoutils.LoadKeys()
	return base64.StdEncoding.EncodeToString(pub)
}

func LoadTempPubKey() string {
	pub, _, _ := cryptoutils.LoadTempKeys()
	return base64.StdEncoding.EncodeToString(pub)
}

func LoadPrivKey() crypto.PrivKey {
	_, priv, _ := cryptoutils.LoadKeys() // priv is ed25519.PrivateKey (a []byte)

	libp2pPriv, err := crypto.UnmarshalEd25519PrivateKey(priv)
	if err != nil {
		panic(err)
	}
	return libp2pPriv
}
