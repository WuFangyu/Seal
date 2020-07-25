package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"log"
	"encoding/hex"
	aes "github.com/WuFangyu/FileEncryption"
	ecies "github.com/ecies/go"
)


func genKeyECC() (*ecies.PrivateKey, error) {
	return ecies.GenerateKey()
}


func encryptECC(pubKey *ecies.PublicKey, msg []byte) []byte {
	ciphterText, err := ecies.Encrypt(pubKey, msg)
	if err != nil {
		log.Fatal(err)
	}
	return ciphterText
}


func decryptECC(privKey *ecies.PrivateKey, msg []byte) []byte {
	plainText, err := ecies.Decrypt(privKey, msg)
	if err != nil {
		log.Fatal(err)
	}
	return plainText
}


func hexStringToPubKey(keyString string) *ecies.PublicKey {
	pubKey, err := ecies.NewPublicKeyFromHex(keyString)
	if err != nil {
		log.Fatal(err)
	}
	return pubKey
}


func pubKeyToHexString(pubKey *ecies.PublicKey) string {
	return pubKey.Hex(true)
}


func hexStringToPrivKey(keyString string) *ecies.PrivateKey {
	privKey, err := ecies.NewPrivateKeyFromHex(keyString)
	if err != nil {
		log.Fatal(err)
	}
	return privKey
}


func privKeyToHexString(privKey *ecies.PrivateKey) string {
	return privKey.Hex()
}


func genKeyRSA(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		log.Fatal(err)
	}
	return privkey, &privkey.PublicKey
}


func encryptRSA(msg []byte, pubkey *rsa.PublicKey) []byte {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pubkey, msg, nil)
	if err != nil {
		log.Fatal(err)
	}
	return ciphertext
}


func decryptRSA(ciphertext []byte, privkey *rsa.PrivateKey) []byte {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, privkey, ciphertext, nil)
	if err != nil {
		log.Fatal(err)
	}
	return plaintext
}


func genKeyAES () [] byte{
	key := make([] byte, 32)
	_, err := rand.Read(key)
	if err != nil{
		log.Fatal(err)
	}
	return key
}


func byteToHex(key []byte) string{
	return hex.EncodeToString(key)
}


func hexToByte(key string) []byte{
	res, err := hex.DecodeString(key)
	if err != nil{
		log.Fatal(err)
	}
	return res
}


func encryptAES (filePath string, key[]byte) {
	aes.InitializeBlock(key)
	err := aes.Encrypter(filePath, formatDirPath(tmpDir))
	if err != nil {
	  log.Fatal(err)
	}
}


func decryptAES (filePath string, outDirPath string, key[]byte) {
	outDirPath = formatDirPath(outDirPath)
	aes.InitializeBlock(key)
	err := aes.Decrypter(filePath, outDirPath)
	if err != nil {
	  log.Fatal(err)
	}
}
