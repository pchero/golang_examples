package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	data := "1234뮻ㅇ!@#$%^&^%$#"

	sEnc := base64.StdEncoding.EncodeToString([]byte(data))
	fmt.Println(sEnc)
}
