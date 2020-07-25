package main


import (
	"log"
	"os"
	"net"
	"io"
	"strconv"
	"strings"
	"math"
	"time"
	bar "github.com/schollz/progressbar"
)


func sendFile(filePath string, conn net.Conn){

	defer conn.Close()

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
	}

	numChunks := min(int(math.Ceil(float64(fileInfo.Size())/float64(defaultChunkSize))), maxChunk)
	chunkSize := int(math.Ceil(float64(fileInfo.Size())/float64(numChunks)))

	progess := bar.Default(fileInfo.Size())

	/*
	var conn net.Conn
	for {
		conn, err = net.Dial("tcp", remoteAddr)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
		}else{
			break
		}
	}
	defer conn.Close()
	fmt.Println("Connection established!")
	*/

	/*  meta info */
	_, err = conn.Write([]byte(makeFileMeta(fileInfo.Name(), fileInfo.Size(), numChunks, chunkSize)))
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, bufSize)
	var recvAddrStrs []string
	for{
		num, err := conn.Read(buf)
		recvAddrStrs = append(recvAddrStrs, string(buf[:num]))
		if err != nil || num == 0 {
			if err != io.EOF{
				log.Fatal(err)
			}
			break
		}
	}

	recvSocketAddrs := strings.Split(strings.Join(recvAddrStrs, ""), msgDelimiter)

	numChunks = len(recvSocketAddrs)

	/*  meta info */
	c := make(chan int)
	var begin int64 = 0
	var size int64 = fileInfo.Size()

	// var wg sync.WaitGroup
	// wg.Add(numChunks)

	for i := 0; i < numChunks; i++ {
		if i == numChunks - 1{
			go sendChunk(recvSocketAddrs[i], c, begin, size, i, filePath, progess)
		}else{
			go sendChunk(recvSocketAddrs[i], c, begin, begin+int64(chunkSize), i, filePath, progess)
		}
		begin += int64(chunkSize)
	}

	for j := 0; j < numChunks; j++ {
		<-c
	}

	// fmt.Println("Finished")
}


func sendChunk(remoteAddr string, c chan int, begin int64, end int64, order int, filePath string,  progess *bar.ProgressBar){

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

	defer conn.Close()
	// fmt.Println("Connection established!: " + strconv.Itoa(order))
	buf := make([]byte, bufSize)
	_, err = conn.Write([]byte(strconv.Itoa(order)))
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("begin recved: " + string(buf[:num]))

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.Seek(begin, 0)

	for i := begin; i < end; i += bufSize {
		len, err:= file.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		if i + int64(len) > end {
			updateLen := int(end - i)
			_, err := conn.Write(buf[:updateLen])
			progess.Add(updateLen)
			if err != nil {
				log.Fatal(err)
			}
		}else{
			_, err := conn.Write(buf[:len])
			progess.Add(len)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	c <- order
	// fmt.Println("Finihsed!: " + strconv.Itoa(order))
}
