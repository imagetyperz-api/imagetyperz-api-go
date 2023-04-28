package main

import (
	"encoding/base64"
	"fmt"
	"github.com/imagetyperz-api/imagetyperz-api-go"
	"io/ioutil"
	"log"
)

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func readImageB64(filePath string) string {
	// Read the entire file into a byte slice
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var base64Encoding string

	//// Determine the content type of the image file
	//mimeType := http.DetectContentType(bytes)
	//
	//// Prepend the appropriate URI scheme header depending
	//// on the MIME type
	//switch mimeType {
	//case "image/jpeg":
	//	base64Encoding += "data:image/jpeg;base64,"
	//case "image/png":
	//	base64Encoding += "data:image/png;base64,"
	//}

	// Append the base64 encoded output
	base64Encoding += toBase64(bytes)

	// Print the full base64 representation of the image
	return base64Encoding
}

func main() {
	api := imagetyperzapi.New("YOUR_ACCESS_TOKEN")
	// get the user balance
	balance, err := api.GetBalance()
	if err != nil {
		fmt.Printf("ERROR balance: %s\n", err.Error())
		return
	}
	fmt.Printf("Balance: $%f\n", balance)

	// submit captcha
	b64Image := readImageB64("captcha.jpg")
	parameters := map[string]string{
		"image": b64Image,
	}
	//parameters["is_case"] = "1"
	//parameters["is_phrase"] = "1"
	//parameters["is_math"] = "1"
	//parameters["digits_only"] = "1"
	//parameters["letters_only"] = "1"
	//parameters["min_length"] = "3"
	//parameters["max_length"] = "6"
	//parameters["affiliate_id"] = "affiliate_id from /account webpage"
	fmt.Println("Solving image captcha...")
	text, err := api.SolveImage(parameters)
	if err != nil {
		fmt.Printf("ERROR submit: %s\n", err.Error())
		return
	}
	fmt.Printf("Text: %s\n", text)

	// if captcha was solved incorrectly, set it as bad
	//err = api.SetCaptchaBad(captchaId)
	//fmt.Println(err)
}
