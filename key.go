package main

import (
	"encoding/base32"
	"fmt"
	"os"
)

func main() {
	secret := os.Args[1]
	key := base32.StdEncoding.EncodeToString([]byte(secret))
	fmt.Println(key)
}
