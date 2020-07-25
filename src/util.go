package main

import (
	"strconv"
	"strings"
	"log"
	"net"
	"fmt"
	"math/big"
	"path/filepath"
	"os"
)


func parseFileMeta(req string) (string, int64, int, int) {
	res := strings.Split(req, msgDelimiter)
	fileSize, err := strconv.ParseInt(res[1], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	numChunks, err := strconv.Atoi(res[2])
	if err != nil {
		log.Fatal(err)
	}
	chunkSize, err := strconv.Atoi(res[3])
	if err != nil {
		log.Fatal(err)
	}
	return res[0], fileSize, numChunks, chunkSize
}


func makeFileMeta(name string, fileSize int64, numChunks int, chunkSize int) string {
	return name + msgDelimiter + strconv.FormatInt(fileSize, 10) + msgDelimiter + strconv.Itoa(numChunks) + msgDelimiter + strconv.Itoa(chunkSize)
}


func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}


func getHostIP(hostSocketInfo *string) {
	
	netIfs, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range netIfs {
		addrs, err := i.Addrs()
		if err != nil {
			log.Fatal(err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
					ip = v.IP
			case *net.IPAddr:
					ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip != nil {
				*hostSocketInfo = ip.String()
				return
			}
		}
	}

	/*
	// get public IP address
	url := "https://api.ipify.org?format=text"
	fmt.Printf("Starting...\n")
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("public IP address: %s\n", ip)
	*/

	log.Fatal("Can't find private IPv4 address!")
}


func createHostSocket(hostIP string) net.Listener {
	var curSocketInfo string
	for i := portStart; i < portEnd; i++ {
		curSocketInfo = hostIP + ":" + strconv.Itoa(i)
		hostSocket, err := net.Listen("tcp", curSocketInfo)
		if err == nil {
			return hostSocket
		} 	
	}
	log.Fatal("No port available!")
	return nil
}


func toHexInt(n *big.Int) string {
    return fmt.Sprintf("%x", n)
}


func formatDirPath(dirPath string) string{
	if dirPath[:len(dirPath) - 1] != "/"{
		return dirPath + "/"
	}
	return dirPath
}


func infoPrint(msg string) {
	msg = msg + "\n"
	fmt.Printf(infoColor, msg)
}


func warningPrint(msg string) {
	msg = msg + "\n"
	fmt.Printf(warningColor, msg)
}


func noticePrint(msg string) {
	msg = msg + "\n"
	fmt.Printf(noticeColor, msg)
}


func debugPrint(msg string) {
	msg = msg + "\n"
	fmt.Printf(debugColor, msg)
}


func RemoveContents(dir string) error {
    d, err := os.Open(dir)
    if err != nil {
        return err
    }
    defer d.Close()
    names, err := d.Readdirnames(-1)
    if err != nil {
        return err
    }
    for _, name := range names {
        err = os.RemoveAll(filepath.Join(dir, name))
        if err != nil {
            return err
        }
    }
    return nil
}
