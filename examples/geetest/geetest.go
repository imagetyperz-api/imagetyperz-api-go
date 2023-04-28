package main

import (
	"fmt"
	"github.com/imagetyperz-api/imagetyperz-api-go"
)

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
	parameters := map[string]string{
		"domain":    "https://your-site.com",
		"challenge": "eea8d7d1bd1a933d72a9eda8af6d15d3",
		"gt":        "1a761081b1114c388092c8e2fd7f58bc",
	}
	//parameters["api_server"] = "api.geetest.com"
	//parameters["user_agent"] = "user agent for solving captcha"
	//parameters["proxy"] = "126.45.34.53:123 or 126.45.34.53:123:joe:password"
	//parameters["affiliate_id"] = "affiliate_id from /account webpage"
	captchaId, err := api.SubmitGeetest(parameters)
	if err != nil {
		fmt.Printf("ERROR submit: %s\n", err.Error())
		return
	}

	// wait for captcha to be solved
	fmt.Printf("Waiting for captcha #%d to be solved ...\n", captchaId)
	solution, err := api.Solve(captchaId, 10)
	if err != nil {
		fmt.Printf("ERROR solve: %s\n", err.Error())
		return
	}
	fmt.Printf("Solution: %s\n", solution)

	// if captcha was solved incorrectly, set it as bad
	//err = api.SetCaptchaBad(captchaId)
	//fmt.Println(err)
}
