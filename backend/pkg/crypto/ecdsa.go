package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
)

// HexToECDSAPrivateKey converts a hexadecimal string representing a private key
// into an *ecdsa.PrivateKey for the secp256k1 curve.
func HexToECDSAPrivateKey(hexKey string) (*ecdsa.PrivateKey, error) {
	trimmedHexKey := strings.TrimPrefix(hexKey, "0x")

	privKeyBytes, err := hex.DecodeString(trimmedHexKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex string: %v", err)
	}

	privKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)

	return privKey.ToECDSA(), nil
}

// HexToECDSAPublicKey converts a hexadecimal string representing a public key
// into an *ecdsa.PublicKey for the secp256k1 curve.
func HexToECDSAPublicKey(hexKey string) (*ecdsa.PublicKey, error) {
	trimmedHexKey := strings.TrimPrefix(hexKey, "0x")

	pubKeyBytes, err := hex.DecodeString(trimmedHexKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex string: %v", err)
	}

	pubKey, err := btcec.ParsePubKey(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	return pubKey.ToECDSA(), nil
}

// SignMessage signs a message using the given ECDSA private key.
func SignMessage(privKey *ecdsa.PrivateKey, message []byte) []byte {
	hash := sha256.Sum256(message)

	signature, err := ecdsa.SignASN1(rand.Reader, privKey, hash[:])
	if err != nil {
		log.Fatalf("failed to sign message: %v", err)
	}

	return signature
}

// SignMessage signs a message using the given ECDSA private key.
func SignMessageAsBase64(privKey *ecdsa.PrivateKey, message []byte) string {
	hash := sha256.Sum256(message)

	signature, err := ecdsa.SignASN1(rand.Reader, privKey, hash[:])
	if err != nil {
		log.Fatalf("failed to sign message: %v", err)
	}

	signatureB64 := base64.StdEncoding.EncodeToString(signature)

	return signatureB64
}
