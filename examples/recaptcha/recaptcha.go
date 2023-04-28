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
		"page_url": "https://your-site.com",
		"sitekey":  "6LrGJmcUABBAALFtIb_FxC0LXm_GwOLyJAfbbUCL",
	}
	//reCAPTCHA type(s) - optional, defaults to 1
	//---------------------------------------------
	//1 - v2
	//2 - invisible
	//3 - v3
	//4 - enterprise v2
	//5 - enterprise v3
	//parameters["type"] = "1"
	//parameters["v3_action"] = "v3 recaptcha action"
	//parameters["v3_min_score"] = "0.3"
	//parameters["domain"] = "www.google.com"
	//parameters["data-s"] = "recaptcha data-s parameter used in loading reCAPTCHA"
	//parameters["cookie_input"] = "a=b;c=d"
	//parameters["user_agent"] = "user agent for solving captcha"
	//parameters["proxy"] = "126.45.34.53:123 or 126.45.34.53:123:joe:password"
	//parameters["affiliate_id"] = "affiliate_id from /account webpage"
	captchaId, err := api.SubmitRecaptcha(parameters)
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
