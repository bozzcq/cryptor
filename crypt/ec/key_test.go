package ec_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/thee-engineer/cryptor/crypt"

	"github.com/thee-engineer/cryptor/crypt/ec"
)

func TestECDSAKeysImportExport(t *testing.T) {
	t.Parallel()

	// Generate Go ECDSA key pair
	ecdsaKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Error(err)
	}

	// Import Go ECDSA key as custom ec key
	prv := ec.Import(ecdsaKey)
	pub := prv.PublicKey

	// Export custom ec key back to Go ECDSA
	ecdsaKeyExport := prv.Export()
	ecdsaPubKeyExport := pub.Export()

	// Compare private keys
	if ecdsaKey.D != ecdsaKeyExport.D {
		t.Error("ecdsa: exported key does not match")
	}

	// Compare public keys
	if ecdsaPubKeyExport.X != ecdsaKey.PublicKey.X &&
		ecdsaPubKeyExport.Y != ecdsaKey.PublicKey.Y {
		t.Error("ecdsa: exported key does not match")
	}
}

func TestECDSACompare(t *testing.T) {
	t.Parallel()

	// Generate key pairs
	key0, key1, err := generateKeyParis()
	if err != nil {
		t.Error(err)
	}
	// Clone one of the keys
	key0Clone := ec.Import(key0.Export())

	// Check for equal keys on different keys
	if key0.IsEqual(key1) || key1.IsEqual(key0) {
		t.Error("ecdsa: unexpected key equality")
	}

	// Compare two equal keys
	if !key0.IsEqual(key0Clone) {
		t.Error("ecdsa: failed to find equal keys")
	}

	// Compare two equal public keys
	if !key0.PublicKey.IsEqual(&key0Clone.PublicKey) {
		t.Error("ecdsa: failed to find public equal keys")
	}

	// Compare two different public keys
	if key0.PublicKey.IsEqual(&key1.PublicKey) {
		t.Error("ecdsa: unexpected public key equality")
	}
}

func TestKeyEncoding(t *testing.T) {
	t.Parallel()

	key, err := ec.GenerateKey()
	if err != nil {
		t.Error(err)
	}

	outByte := key.Encode()
	outString := key.EncodeString()

	if crypt.EncodeString(outByte) != outString {
		t.Log(crypt.EncodeString(outByte))
		t.Log(outString)
		t.Error("ec key | mismatch key encodings")
	}

	decodedKey, err := ec.Decode(outByte)
	if err != nil {
		t.Error(err)
	}

	if !key.IsEqual(decodedKey) {
		t.Log("decoded: ", decodedKey)
		t.Log("original:", key)
		t.Errorf("ec key | mismatch original with decoded key")
	}

	decodedKey, err = ec.DecodeString(outString)
	if err != nil {
		t.Error(err)
	}

	if !key.IsEqual(decodedKey) {
		t.Log("decoded: ", decodedKey)
		t.Log("original:", key)
		t.Errorf("ec key | mismatch original with decoded key")
	}

	if _, err := ec.Decode([]byte{10, 20, 30}); err == nil {
		t.Error("ec key | decoded invalid []byte")
	}

	if _, err := ec.Decode(crypt.RandomData(65)); err == nil {
		t.Error("ec key | decoded invalid []byte")
	}

	if _, err := ec.DecodeString("testing"); err == nil {
		t.Error("ec key | decoded invalid string")
	}
}
