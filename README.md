imagetyperz-api-go - Imagetyperz API wrapper
=========================================

imagetyperz-api-go is a super easy to use bypass captcha API wrapper for imagetyperz.com captcha service

## Installation
```bash
go get github.com/imagetyperz-api/imagetyperz-api-go
```

## Usage

``` python
import (
	"github.com/imagetyperz-api/imagetyperz-api-go"
)
```
Initialize the API object with your access token and start using it

``` python
api := imagetyperzapi.New("YOUR_ACCESS_TOKEN")
```

**Get balance**

``` python
balance, err := api.GetBalance()
if err != nil {
    fmt.Printf("ERROR balance: %s\n", err.Error())
    return
}
fmt.Printf("Balance: $%f\n", balance)
```

## Solving
For solving a captcha, it's a two step process:
- **submit captcha** details - returns an ID
- use ID to check it's progress - and **get solution** when solved.

Each captcha type has it's own submission method.

For getting the response, same method is used for all types, except for the image (classic) captcha, which has only one method for both submitting and getting the text solution.

### Image captcha

``` go
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
```
### reCAPTCHA

For recaptcha submission there are two things that are required.
- page_url
- site_key
- type (optional, defaults to 1 if not given)
  - `1` - v2
  - `2` - invisible
  - `3` - v3
  - `4` - enterprise v2
  - `5` - enterprise v3
- domain - used in loading of reCAPTCHA interface, default: `www.google.com` (alternative: `recaptcha.net`) - `optional`
- v3_min_score - minimum score to target for v3 recaptcha `- optional`
- v3_action - action parameter to use for v3 recaptcha `- optional`
- proxy - proxy to use when solving recaptcha, eg. `12.34.56.78:1234` or `12.34.56.78:1234:user:password` `- optional`
- user_agent - useragent to use when solve recaptcha `- optional` 
- data-s - extra parameter used in solving recaptcha `- optional`
- cookie_input - cookies used in solving reCAPTCHA `- optional`

``` go
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
```
ID will be used to retrieve the g-response (solution), once workers have 
completed the captcha. This takes somewhere between 10-80 seconds. 

### GeeTest

GeeTest is a captcha that requires 3 parameters to be solved:
- domain
- challenge
- gt
- api_server (optional)
- user_agent (optional)
- proxy (optional)

The response of this captcha after completion are 3 codes:
- challenge
- validate
- seccode

**Important**
This captcha requires a **unique** challenge to be sent along with each captcha.

```go
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
```

Optionally, you can send proxy and user_agent along.

### GeeTestV4

GeeTesV4 is a new version of captcha from geetest that requires 2 parameters to be solved:

- domain
- geetestid (captchaID) - gather this from HTML source of page with captcha, inside the `<script>` tag you'll find a link that looks like this: https://i.imgur.com/XcZd47y.png

The response of this captcha after completion are 5 parameters:

- captcha_id
- lot_number
- pass_token
- gen_time
- captcha_output

```go
parameters := map[string]string{
    "domain":    "https://example.com",
    "geetestid": "647f5ed2ed8acb4be36784e01556bb71",
}
//parameters["user_agent"] = "user agent for solving captcha"
//parameters["proxy"] = "126.45.34.53:123 or 126.45.34.53:123:joe:password"
//parameters["affiliate_id"] = "affiliate_id from /account webpage"
captchaId, err := api.SubmitGeetestV4(parameters)
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
```

Optionally, you can send proxy and user_agent along.

### hCaptcha

Requires page_url and sitekey

```go
parameters := map[string]string{
    "page_url": "https://your-site.com",
    "sitekey":  "8c7062c7-cae6-4e12-96fb-303fbec7fe4f",
}
//parameters["invisible"] = "1"
//parameters["enterprise_payload"] = "{\"rqdata\": \"take value from web requests\"}"
//parameters["domain"] = "hcaptcha.com"
//parameters["user_agent"] = "user agent for solving captcha"
//parameters["proxy"] = "126.45.34.53:123 or 126.45.34.53:123:joe:password"
//parameters["affiliate_id"] = "affiliate_id from /account webpage"
captchaId, err := api.SubmitHcaptcha(parameters)
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
```

### Capy

Requires page_url and sitekey

```go
parameters := map[string]string{
    "page_url": "https://your-site.com",
    "sitekey":  "Fme6hZLjuCRMMC3uh15F52D3uNms5c",
}
//parameters["user_agent"] = "user agent for solving captcha"
//parameters["proxy"] = "126.45.34.53:123 or 126.45.34.53:123:joe:password"
//parameters["affiliate_id"] = "affiliate_id from /account webpage"
captchaId, err := api.SubmitCapy(parameters)
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
```

### Tiktok

Requires page_url cookie_input

```go
parameters := map[string]string{
    "page_url":     "https://tiktok.com",
    "cookie_input": "s_v_web_id:verify_kd6243o_fd449FX_FDGG_1x8E_8NiQ_fgrg9FEIJ3f;tt_webid:612465623570154;tt_webid_v2:7679206562717014313;SLARDAR_WEB_ID:d0314f-ce16-5e16-a066-71f19df1545f;",
}
//parameters["user_agent"] = "user agent for solving captcha"
//parameters["proxy"] = "126.45.34.53:123 or 126.45.34.53:123:joe:password"
//parameters["affiliate_id"] = "affiliate_id from /account webpage"
captchaId, err := api.SubmitTiktok(parameters)
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
```

### FunCaptcha

Requires page_url, sitekey and s_url (source URL)

```go
parameters := map[string]string{
    "page_url": "https://your-site.com",
    "sitekey":  "11111111-1111-1111-1111-111111111111",
    "s_url":    "https://api.arkoselabs.com",
}
//parameters["data"] = "{\"a\": \"b\"}"
//parameters["user_agent"] = "user agent for solving captcha"
//parameters["proxy"] = "126.45.34.53:123 or 126.45.34.53:123:joe:password"
//parameters["affiliate_id"] = "affiliate_id from /account webpage"
captchaId, err := api.SubmitFuncaptcha(parameters)
if err != nil {
    fmt.Printf("ERROR submit: %s\n", err.Error())
    return
}

// wait for captcha to be solved
fmt.Printf("Waiting for captcha #%d to be solved ...\n", captchaId)
solution, err := api.Solve(captchaId, 1)
if err != nil {
    fmt.Printf("ERROR solve: %s\n", err.Error())
    return
}
fmt.Printf("Solution: %s\n", solution)
```

### Turnstile (Cloudflare)

```go
parameters := map[string]string{
    "page_url": "https://your-site.com",
    "sitekey":  "0x4ABBBBAABrfvW5vKbx11FZ",
}
//parameters["domain"] = "challenges.cloudflare.com"
//parameters["action"] = "homepage"
//parameters["cdata"] = "cdata information"
//parameters["user_agent"] = "user agent for solving captcha"
//parameters["proxy"] = "126.45.34.53:123 or 126.45.34.53:123:joe:password"
//parameters["affiliate_id"] = "affiliate_id from /account webpage"
captchaId, err := api.SubmitTurnstile(parameters)
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
```

### Task

Requires template_name, page_url and usually variables

```go
parameters := map[string]string{
    "template_name": "Login test page",
    "page_url":      "https://imagetyperz.net/automation/login",
    "variables":     "{\"username\": 'abc', \"password\": 'paZZW0rd'}",
}
//parameters["user_agent"] = "user agent for solving captcha"
//parameters["proxy"] = "126.45.34.53:123 or 126.45.34.53:123:joe:password"
//parameters["affiliate_id"] = "affiliate_id from /account webpage"
captchaId, err := api.SubmitTask(parameters)
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
```

#### Task pushVariable
Update a variable value while task is running. Useful when dealing with 2FA authentication and similar situations.

When template reaches an action that uses a variable which wasn't provided with the submission of the task,
task (while running on worker machine) will wait for variable to be updated through push.

You can use the `PushTaskVariables` method as many times as you need, even overwriting variables that were set previously.

```go
PushTaskVariables(captchaId int64, pushVariables string) error
```

### Set captcha bad

When a captcha was incorrectly solved by our workers, you can notify the server with it's ID, so we know something went wrong.

``` go
SetCaptchaBad(captchaId int64) error
```

### Response (solution)

The response is stringified JSON that looks like this:

```json
{
  "CaptchaId": 176707908, 
  "Response": "03AGdBq24PBCbwiDRaS_MJ7Z...mYXMPiDwWUyEOsYpo97CZ3tVmWzrB", 
  "Cookie_OutPut": "", 
  "Proxy_reason": "",
  "Status": "Solved"
}
```

## Examples
Check `examples` folder, which contains an example for each type of captcha.

## License
API library is licensed under the MIT License

## More information
More details about the server-side API can be found [here](http://imagetyperz.com)


<sup><sub>captcha, bypasscaptcha, decaptcher, decaptcha, 2captcha, deathbycaptcha, anticaptcha, 
bypassrecaptchav2, bypassnocaptcharecaptcha, bypassinvisiblerecaptcha, captchaservicesforrecaptchav2, 
recaptchav2captchasolver, googlerecaptchasolver, recaptchasolverpython, recaptchabypassscript</sup></sub>

