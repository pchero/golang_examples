package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func die(msg string, s error) {
	fmt.Printf("%s creshed: %v\n", msg, s)
}

func main() {
	fd, err := os.Open("./test.tar.gz")
	if err != nil {
		fmt.Printf("Could not open the file. err: %s\n", err)
		return
	}

	read, err := ioutil.ReadAll(fd)
	if err != nil {
		fmt.Printf("Could not open the file. err: %s\n", err)
		return
	}
	fmt.Print("Read: %s\n", read)

	// fd, err := net.Dial("tcp", "google.com:80")
	// if err != nil {
	// 	die("dial", err)
	// }

	// // req := []byte("GET /intl/en/privacy/ HTTP/1.0\r\nHost: www.google.com\r\n\r\n")
	// req := []byte("GET HTTP/1.0\r\nHost: https://www.naver.com/\r\n\r\n")

	// _, err = fd.Write(req)
	// if err != nil {
	// 	die("dial write", err)
	// }

	// // buf := make([]byte, 1024)
	// // var nr int64
	// // nr = 1

	// // for nr != 0 {

	// res, err := ioutil.ReadAll(fd)
	// if err != nil {
	// 	die("dial read", err)
	// }
	// fmt.Printf("Read: %s\n", res)
	// // nr, err = io.ReadFull(fd, buf)
	// // var buf bytes.Buffer
	// // nr, err = io.Copy(&buf, fd)
	// // // nr, err = ioutil.ReadAll(bufio.NewReader(fd))
	// // if err != nil {
	// // 	die("dial read", err)
	// // }
	// // fmt.Printf("read %d\n", nr)
	// // }
}
