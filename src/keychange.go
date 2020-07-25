package main

import (
	"net"
	"log"
	"errors"
	"time"
	"fmt"
)


func sendConnSetup(remoteAddr string, pubKey string, aesKey []byte) (net.Conn, error) {
	var conn net.Conn
	var err error
	for {
		conn, err = net.Dial("tcp", remoteAddr)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
		}else{
			break
		}
	}

	pubKeyByte := []byte(pubKey)

	pubKeyECC :=  hexStringToPubKey(pubKey)
	aesKeyCiphter := encryptECC(pubKeyECC, aesKey)
	buf := append(pubKeyByte, aesKeyCiphter...)

	conn.Write(buf)
	num , err := conn.Read(buf)
    if err != nil || string(buf[:num]) != okSignal {
		conn.Close()
		return nil, errors.New("Incorrect message format!")
	}

	return conn, nil
}


func recvConnSetup(hostSocket net.Listener, pubKey string, privKey string) (net.Conn, []byte, error) {
	// listening for tcp traffic
	conn, err := hostSocket.Accept()
    if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, bufSize)
	num , err := conn.Read(buf)
    if err != nil {
		log.Fatal(err)
	}

	if num  <= keylenECC { // 66 bytes for ECC key and 32 bytes for AES key
		conn.Write([]byte(errSignal))
		conn.Close()
		fmt.Println(num)
		return nil, nil, errors.New("Incorrect message format!")
	}
	recvPubKey := string(buf[:keylenECC])
	if recvPubKey != pubKey {
		fmt.Println(recvPubKey)
		fmt.Println(pubKey)
		conn.Close()
		return nil, nil, errors.New("Incorrect public key!")
	}

	aesKeyCiphter := buf[keylenECC:num]

	privKeyECC := hexStringToPrivKey(privKey)

	aesKey:= decryptECC(privKeyECC, aesKeyCiphter)

	conn.Write([]byte(okSignal))

	return conn, aesKey, nil
}