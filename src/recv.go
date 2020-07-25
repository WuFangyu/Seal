package main

import (
	"fmt"
	"log"
	"os"
	"net"
	"io"
	"strconv"
	"strings"
	"sync"
	bar "github.com/schollz/progressbar"
)


func recvFile(conn net.Conn, hostIP string, destDir string) string {
	defer conn.Close()

	destDir = formatDirPath(destDir)
	// fmt.Println("Connection established!")

	buf := make([]byte, bufSize)
	num , err := conn.Read(buf)
    if err != nil {
		log.Fatal(err)
	}

	fileName, fileSize, numChunks, chunkSize := parseFileMeta(string(buf[:num]))
	progess := bar.Default(fileSize)

	filePath := destDir + fileName
	file, err := os.Create(filePath)
    if err != nil {
		log.Fatal(err)
	}
	
	file.Close()

	// send over availble ports
	var socketCount int = 0
	var recvSockets []net.Listener
	var recvSocketAddrs [] string
	
	for i := portStart; i < portEnd; i++ {
		curSocketInfo := hostIP + ":" + strconv.Itoa(i)
		recvSocket, err := net.Listen("tcp", curSocketInfo)
		if err == nil {
			socketCount += 1
			recvSockets = append(recvSockets, recvSocket)
			recvSocketAddrs = append(recvSocketAddrs, curSocketInfo)
		}
		if socketCount == numChunks {
			break
		}
	}

	/* 
	// print out meta info
	fmt.Println(len(recvSocketAddrs))
	fmt.Println(numChunks)
	*/

	var wg sync.WaitGroup
	wg.Add(numChunks)

	for i := 0; i < numChunks; i++ {
		go recvFileChunk(recvSockets[i], chunkSize, &wg, filePath, progess)
	}
	
	addrsBytes := []byte(strings.Join(recvSocketAddrs, msgDelimiter))
	offset := 0
	lenAddrs := len(addrsBytes)

	for{
		conn.Write(addrsBytes[offset:min(offset + bufSize, lenAddrs)])
		offset += bufSize
		if offset > lenAddrs {
			break;
		}
	}

	numChunks = socketCount
	conn.Close()

	wg.Wait()

	return fileName
}


func recvFileChunk(recvSocket net.Listener, chunkSize int, wg *sync.WaitGroup, filePath string, progess *bar.ProgressBar) {
	defer wg.Done()
	conn, err := recvSocket.Accept()
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, bufSize)
	var num int
	for {
		lnum, err := conn.Read(buf)
		num = lnum
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			continue
		}
		break
	}

	begin , err := strconv.ParseInt(string(buf[:num]), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	conn.Write([]byte(okSignal))
	begin = begin * int64(chunkSize)

	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
    if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.Seek(begin, 0)

	for {
		num, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error Read Data")
				log.Fatal(err)
			}
			break
		}
		progess.Add(num)
		file.Write(buf[:num])
	}

	// fmt.Println(strconv.Itoa(int(begin)) + " is done!")
}
