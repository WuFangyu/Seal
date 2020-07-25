package main

import (
	"log"
	"os"
	"github.com/urfave/cli"
)


func main() {
	app := cli.NewApp()
	app.Name = "Seal"
	app.Usage = "Seal is a tool for simple, fast and secure end-to-end file tranfer."
	var hostIP string
	getHostIP(&hostIP)
	hostSocket := createHostSocket(hostIP)
	levelDB := initDB(levelDBpath)
	defer levelDB.Close()

	if len(getKey(levelDB, hostPubKeyIdx)) == 0 {
		eccKey, err := genKeyECC()
		if err != nil {
			log.Fatal(err)
		}
		eccKeyString := pubKeyToHexString(eccKey.PublicKey)
		storeKey(levelDB, hostPubKeyIdx, eccKeyString)
		eccPrivKeyString := privKeyToHexString(eccKey)
		storeKey(levelDB, hostPrivKeyIdx, eccPrivKeyString)
	}

	eccKeyString := getKey(levelDB, hostPubKeyIdx)
	infoPrint("Host: " + hostSocket.Addr().String() + " Public Key: " + eccKeyString)

	sendFlag := []cli.Flag{
		&cli.StringFlag{
			Name: "file, f",
			Usage: "The path of the file to be sent.",
		},
		&cli.StringFlag{
			Name: "pubkey, publickey, key",
			Usage: "The public key from receiver.(Optional, check `-help key`).",
		},
		&cli.StringFlag{
			Name: "dest, dst, ip, IP, remote",
			Usage: "The destination's IPv4 address.",
		},
	}

	recvFlag := []cli.Flag{
		&cli.StringFlag{
			Name: "dir, d",
			Usage: "The directory to place the received file.",
		},
	}

	keyFlag := []cli.Flag{
		&cli.BoolFlag{
			Name: "list, l, ls",
			Usage: "List all stored keys.",
		},
		&cli.BoolFlag{
			Name: "update, u",
			Usage: "Update the given key record.",
		},
		&cli.BoolFlag{
			Name: "put, store, add",
			Usage: "Store the given key record.",
		},
		&cli.StringFlag{
			Name: "remove, rm, delete, d",
			Usage: "Remove the given key record.",
		},
		&cli.StringFlag{
			Name: "get, read, r, getkey",
			Usage: "Retrieve the key value for the given label.",
		},
		&cli.StringFlag{
			Name: "name, label, n",
			Usage: "The name or label for the key.",
		},
		&cli.StringFlag{
			Name: "pubkey, publickey, key, val",
			Usage: "The key value for the given pair.",
		},
		&cli.BoolFlag{
			Name: "genkey, newkey",
			Usage: "Update/Generate a new asymmetric key pair.",
		},
		&cli.BoolFlag{
			Name: "clear, clean",
			Usage: "Clean out all stored key pairs.",
		},
	}


	app.Commands = []*cli.Command{
		{
			Name: "send",
			Usage: "Send file to the destination.",
			Flags: sendFlag,
			Action: func(c *cli.Context) error {

				remoteAddr := c.String("dest")
				if len(remoteAddr) == 0 {
					warningPrint(sendCmdErrMsg)
					return nil
				}

				file := c.String("file")
				if len(file) == 0 {
					warningPrint(sendCmdErrMsg)
					return nil
				}

				pubKey := c.String("pubkey")
				if len(pubKey) == 0 {
					pubKey = getKey(levelDB, remoteAddr)
					if len(pubKey) == 0{
						warningPrint(sendCmdErrMsg)
						return nil
					}
				}

				aesKey := genKeyAES()
				debugPrint("Generated AES key: " + byteToHex(aesKey))

				fileInfo, err := os.Stat(file)
				if err != nil {
					warningPrint(file + " does not exist!")
					return nil
				}

				encryptAES(file, aesKey)
				conn, err := sendConnSetup(remoteAddr, pubKey, aesKey)
				if err != nil{
					warningPrint("Connection failed, please check if the public key is correct!")
					return nil
				}

				tempPath := formatDirPath(tmpDir) + fileInfo.Name() + encExt
				sendFile(tempPath, conn)
				err = os.Remove(tempPath)
				if err != nil {
					log.Fatal(tempPath + " does not exits!")
				}
				return nil
			},
		},
		{
			Name: "recv",
			Usage: "Receive file from the remote.",
			Flags: recvFlag,
			Action: func(c *cli.Context) error {

				dirPath := c.String("dir")

				if len(dirPath) == 0 {
					warningPrint(recvCmdErrMsg)
					return nil
				}

				if _, err := os.Stat(dirPath); os.IsNotExist(err) {
					warningPrint(dirPath + " doest not exits!")
					return nil
				}

				eccSecret := getKey(levelDB, hostPrivKeyIdx)
				conn, aesKey, err := recvConnSetup(hostSocket, eccKeyString, eccSecret)
				if err != nil{
					warningPrint("Connection failed!")
					return nil
				}

				debugPrint("Received AES key: " + byteToHex(aesKey))

				tempDir := formatDirPath(tmpDir)
				fileName := recvFile(conn, hostIP, tmpDir)

				filePath := tempDir + fileName
				decryptAES(filePath, dirPath, aesKey)

				err = os.Remove(filePath)
				if err != nil {
					log.Fatal(filePath + " does not exits!")
				}
				return nil
			},
		},
		{
			Name: "key",
			Flags: keyFlag,
			Usage: "Manage cached key infomation.",
			Action: func(c *cli.Context) error {

				if c.Bool("list") {
					printKeys(levelDB)
					return nil
				}

				if c.Bool("genkey") {
					eccKey, err := genKeyECC()
					if err != nil {
						log.Fatal(err)
					}
					eccPubKeyString := pubKeyToHexString(eccKey.PublicKey)
					storeKey(levelDB, hostPubKeyIdx, eccPubKeyString)
					eccPrivKeyString := privKeyToHexString(eccKey)
					storeKey(levelDB, hostPrivKeyIdx, eccPrivKeyString)
					noticePrint("New public key: " + eccPubKeyString)
					return nil
				}

				if c.Bool("clear") {
					clearDB()
					noticePrint("All cached key pairs have been removed.")
					return nil
				}

				queryKeyName := c.String("get")
				if len(queryKeyName) != 0 {
					res := getKey(levelDB, queryKeyName)
					if len(res) == 0 {
						noticePrint("No key info for the given name: " + queryKeyName)
						return nil
					}
					noticePrint("Name: " + queryKeyName +  " Key:" + res)
					return nil
				}

				queryKeyName = c.String("remove")
				if len(queryKeyName) != 0 {
					res := getKey(levelDB, queryKeyName)
					if len(res) == 0 {
						noticePrint("No public info for the given address: " + queryKeyName)
						return nil
					}
					removeKey(levelDB, queryKeyName)
					noticePrint("Removed key pair: " + queryKeyName + ", " + res)
					return nil
				}

				update := c.Bool("update")
				add := c.Bool("put")

				if add || update {
					givenName := c.String("name")
					givenKey := c.String("pubkey")
					if (len(givenName) == 0 || len(givenKey) == 0) {
							warningPrint(keyCmdErrMsg)
							return nil
					}

					storeKey(levelDB, givenName, givenKey)
					noticePrint("Saved key pair: " + givenName + ", " + givenKey)
					return nil
				}

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
