package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

const (
	CLOCK_STEP             = 30 // in seconds
	ONE_TIME_PASSWD_DIGITS = 4
	PAYCODE_FMT            = "28%04d%08d%04d"
	UID                    = 10203405609
)

func toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func toUint32(bytes []byte) uint32 {
	return (uint32(bytes[0]) << 24) + (uint32(bytes[1]) << 16) +
		(uint32(bytes[2]) << 8) + uint32(bytes[3])
}

func oneTimePassword(key []byte, value []byte) uint32 {
	// sign the value using HMAC-SHA1
	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(value)
	hash := hmacSha1.Sum(nil)

	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	offset := hash[len(hash)-1] & 0x0F

	// get a 32-bit (4-byte) chunk from the hash starting at offset
	hashParts := hash[offset : offset+4]

	// ignore the most significant bit as per RFC 4226
	hashParts[0] = hashParts[0] & 0x7F

	number := toUint32(hashParts)

	// size to 6 digits
	// one million is the first number with 7 digits so the remainder
	// of the division will always return < 7 digits
	pwd := number % uint32(math.Pow10(ONE_TIME_PASSWD_DIGITS))
	return pwd
}

// all []byte in this program are treated as Big Endian
func generateOneTimePassword(input string) uint32 {
	// decode the key from the first argument
	inputNoSpaces := strings.Replace(input, " ", "", -1)
	inputNoSpacesUpper := strings.ToUpper(inputNoSpaces)
	key, err := base32.StdEncoding.DecodeString(inputNoSpacesUpper)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// generate a one-time password using the time at 30-second intervals
	// server and client must use the same timezone
	epochSeconds := time.Now().Unix()
	pwd := oneTimePassword(key, toBytes(epochSeconds/CLOCK_STEP))
	return pwd
}

// -- |--------|-----------|-------
// 28 |  x     |  y        |  z
// -- |--------|-----------|-------
// AI  4 digits  9 digits   3 digits  => 18 digits
func demoPaycode(key string) {
	const FACTOR = 5
	x := int(generateOneTimePassword(key))
	y := UID/x + FACTOR*x
	z := UID % x
	paycode := fmt.Sprintf(PAYCODE_FMT, x, y, z)

	// paycode is generated, now given the paycode, decode to uid and validate x

	// get the uid from the paycode
	var origX, origY, origZ int
	fmt.Sscanf(paycode, PAYCODE_FMT, &origX, &origY, &origZ)
	origY -= origX * FACTOR
	origUid := origX*origY + origZ

	// validate the x factor
	var xValid bool
	if int(generateOneTimePassword(key)) == origX {
		// TODO 加入时间不一致的容错机制，容许n minutes的误差
		// 如果客户端时间和服务器时间相差非常大(e,g. 1h)，支付宝的做法是把取到的uid发送
		// 到设备上，让用户自己进行支付确认
		xValid = true
	}

	if origUid != UID || !xValid {
		panic("invalid paycode:" + paycode)
	}

	fmt.Printf("paycode: %s, uid: %d\n", paycode, origUid)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "must specify key to use")
		os.Exit(1)
	}

	fmt.Printf("step: %ds\n", CLOCK_STEP)
	key := os.Args[1]
	for i := 0; i < 100; i++ {
		demoPaycode(key)
		time.Sleep(time.Second)
	}
}
