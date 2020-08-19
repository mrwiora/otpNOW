package main

import (
	"bufio"
	"bytes"
	"encoding/base32"
	"fmt"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// random string generator
// thanks https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890@#!?")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func display(key *otp.Key, data []byte) {
	fmt.Printf("Issuer:       %s\n", key.Issuer())
	fmt.Printf("Account Name: %s\n", key.AccountName())
	fmt.Printf("Secret:       %s\n", key.Secret())
	fmt.Println("Writing PNG to qr-code.png....")
	ioutil.WriteFile("qr-code.png", data, 0644)
	fmt.Println("")
	fmt.Println("Please add your TOTP to your OTP Application now!")
	fmt.Println("")
}

func promptForPasscode() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Passcode: ")
	text, _ := reader.ReadString('\n')
	return text
}

// Demo function, not used in main
// Generates Passcode using a UTF-8 (not base32) secret and custom paramters
func GeneratePassCode(utf8string string) string {
	secret := base32.StdEncoding.EncodeToString([]byte(utf8string))
	passcode, err := totp.GenerateCodeCustom(secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA512,
	})
	if err != nil {
		panic(err)
	}
	return passcode
}

func generateKey() *otp.Key {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Example.com",
		AccountName: "alice@example.com",
	})
	if err != nil {
		panic(err)
	}
	return key
}

func generateQR(key *otp.Key) bytes.Buffer {
	// Convert TOTP key into a PNG
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		panic(err)
	}
	png.Encode(&buf, img)
	return buf
}

func validHandlerTOTP(key *otp.Key, passcode string, SESSIONID string) string {
	valid := totp.Validate(passcode, key.Secret())
	if valid {
		fmt.Printf("%v validHandlerTOTP: valid \n", SESSIONID)
		return "valid"
	} else {
		fmt.Printf("%v validHandlerTOTP: invalid \n", SESSIONID)
		return "invalid"
	}
}

func main() {

	// start server
	key := generateKey()
	// generate QR
	qr := generateQR(key)

	display(key, qr.Bytes())

	// display the QR code to the user.
	// display(key, buf.Bytes())

	// Hello world, the web server
	verifyHandler := func(w http.ResponseWriter, r *http.Request) {
		// creating random string
		rand.Seed(time.Now().UnixNano())
		var SESSIONID = randSeq(10)
		// performing external query for working ssh connection
		passcode := r.URL.Query().Get("passcode")
		fmt.Printf("%v verifyHandler : passcode is %v \n", SESSIONID, passcode)
		out := validHandlerTOTP(key, passcode, SESSIONID)
		io.WriteString(w, string(out))
	}

	http.HandleFunc("/totp", verifyHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
