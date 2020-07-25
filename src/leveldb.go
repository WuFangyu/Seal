package main

import (
	leveldb "github.com/syndtr/goleveldb/leveldb"
	"os"
	"log"
)


func initDB(dirPath string) *leveldb.DB{
	var err error
	if _, err = os.Stat(dirPath); os.IsNotExist(err) {
		os.Mkdir(dirPath, os.FileMode(0700))
	}
	dbClient, err := leveldb.OpenFile(dirPath, nil)
	if err != nil {
		log.Fatal(err)
	}
	return dbClient
}


func getKey(db *leveldb.DB, remoteAddr string) string{
	pubKey, err := db.Get([]byte(remoteAddr), nil)
	if err != nil{
		return ""
	}
	return string(pubKey)
}


func removeKey(db *leveldb.DB, remoteAddr string) {
	err := db.Delete([]byte(remoteAddr), nil)
	if err != nil{
		log.Panic(err)
	}
}


func storeKey(db *leveldb.DB, remoteAddr string, pubKey string) {
	err := db.Put([]byte(remoteAddr), []byte(pubKey),  nil)
	if err != nil{
		log.Panic(err)
	}
}


func clearDB(){
	RemoveContents(levelDBpath)
}


func printKeys(db *leveldb.DB){
	noticePrint("Cached addresses and public keys:")
	empty := true
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		infoPrint(string(key) + ": " + string(value))
		empty = false
	}

	iter.Release()
	err := iter.Error()
	if err != nil{
		log.Panic(err)
	}

	if empty{
		infoPrint("empty")
	}
}
