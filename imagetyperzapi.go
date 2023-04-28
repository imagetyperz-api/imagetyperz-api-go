package imagetyperzapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	USER_AGENT                    = "goClient"
	TIMEOUT                       = 120
	ENDPOINT_BALANCE              = "https://captchatypers.com/Forms/RequestBalanceToken.ashx"
	ENDPOINT_IMAGE                = "https://captchatypers.com/Forms/UploadFileAndGetTextNEWToken.ashx"
	ENDPOINT_RECAPTCHA            = "https://captchatypers.com/captchaapi/UploadRecaptchaToken.ashx"
	ENDPOINT_RECAPTCHA_ENTERPRISE = "https://captchatypers.com/captchaapi/UploadRecaptchaEnt.ashx"
	ENDPOINT_CAPY                 = "https://captchatypers.com/captchaapi/UploadCapyCaptchaUser.ashx"
	ENDPOINT_FUNCAPTCHA           = "https://captchatypers.com/captchaapi/UploadFunCaptcha.ashx"
	ENDPOINT_GEETEST              = "https://captchatypers.com/captchaapi/UploadGeeTestToken.ashx"
	ENDPOINT_GEETESTV4            = "https://captchatypers.com/captchaapi/UploadGeeTestV4.ashx"
	ENDPOINT_HCAPTCHA             = "https://captchatypers.com/captchaapi/UploadHCaptchaUser.ashx"
	ENDPOINT_TIKTOK               = "https://captchatypers.com/captchaapi/UploadTikTokCaptchaUser.ashx"
	ENDPOINT_TURNSTILE            = "https://captchatypers.com/captchaapi/Uploadturnstile.ashx"
	ENDPOINT_TASK                 = "https://captchatypers.com/captchaapi/UploadCaptchaTask.ashx"
	ENDPOINT_TASK_PUSH_VARIABLES  = "https://captchatypers.com/CaptchaAPI/SaveCaptchaPush.ashx"
	ENDPOINT_RESULT               = "https://captchatypers.com/captchaapi/GetCaptchaResponseJson.ashx"
	ENDPOINT_BAD                  = "https://captchatypers.com/Forms/SetBadImageToken.ashx"
)

type Client struct {
	AccessToken string
	Timeout     int
	httpClient  *http.Client
}

// New - create a new API client
func New(accessToken string) *Client {
	return &Client{
		AccessToken: accessToken,
		httpClient: &http.Client{
			Timeout: time.Duration(TIMEOUT) * time.Second,
		},
	}
}

// GetRequest - make GET request to API server
func (c *Client) GetRequest(url string) (string, error) {
	// make HTTP request
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", USER_AGENT)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	// read response body
	defer func(Body io.ReadCloser) {
	}(resp.Body)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bodyString := string(bodyBytes)

	return bodyString, err
}

// PostRequest - make POST request to API server
func (c *Client) PostRequest(endpoint string, parameters map[string]string) (string, error) {
	data := url.Values{}

	for key, value := range parameters {
		data.Set(key, value)
	}
	// make HTTP request
	req, _ := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(data.Encode()))
	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	// read response body
	defer func(Body io.ReadCloser) {
	}(resp.Body)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bodyString := string(bodyBytes)

	// check for error
	if strings.HasPrefix(bodyString, "ERROR") {
		error := strings.Replace(bodyString, "ERROR:", "", 1)
		error = strings.TrimSpace(error)
		return "", errors.New(error)
	}
	if strings.Contains(bodyString, "Error") && !strings.Contains(bodyString, "\"Error\":\"\"") {
		return "", errors.New(bodyString)
	}

	return bodyString, err
}

// GetBalance - get the user balance
func (c *Client) GetBalance() (float64, error) {
	reqUrl := fmt.Sprintf("%s?action=REQUESTBALANCE&submit=Submit&token=%s", ENDPOINT_BALANCE, c.AccessToken)
	response, err := c.GetRequest(reqUrl)
	if err != nil {
		return 0.0, err
	}
	return strconv.ParseFloat(fmt.Sprint(response), 64)
}

// SolveImage - submit (classic) image captcha
func (c *Client) SolveImage(parameters map[string]string) (string, error) {
	reqUrl := ENDPOINT_IMAGE
	p := map[string]string{
		"token":  c.AccessToken,
		"action": "UPLOADCAPTCHA",
		"file":   parameters["image"],
	}

	if parameters["is_case"] != "" {
		p["iscase"] = "true"
	}
	if parameters["is_phrase"] != "" {
		p["isphrase"] = "true"
	}
	if parameters["is_math"] != "" {
		p["ismath"] = "true"
	}
	if parameters["digits_only"] != "" {
		p["alphanumeric"] = "1"
	}
	if parameters["letters_only"] != "" {
		p["alphanumeric"] = "2"
	}
	if parameters["min_length"] != "" {
		p["minlength"] = parameters["min_length"]
	}
	if parameters["max_length"] != "" {
		p["maxlength"] = parameters["max_length"]
	}

	if parameters["affiliate_id"] != "" {
		p["affiliate_id"] = parameters["affiliate_id"]
	}
	// make request to API server
	response, err := c.PostRequest(reqUrl, p)
	if err != nil {
		return "", err
	}
	return strings.Split(response, "|")[1], nil
}

// SubmitRecaptcha - submit a captcha and get back the id
func (c *Client) SubmitRecaptcha(parameters map[string]string) (int64, error) {
	reqUrl := ENDPOINT_RECAPTCHA
	p := map[string]string{
		"token":         c.AccessToken,
		"action":        "UPLOADCAPTCHA",
		"pageurl":       parameters["page_url"],
		"googlekey":     parameters["sitekey"],
		"recaptchatype": "0",
	}
	if parameters["proxy"] != "" {
		p["proxy"] = parameters["proxy"]
		p["proxytype"] = "HTTP"
	}
	if parameters["affiliate_id"] != "" {
		p["affiliate_id"] = parameters["affiliate_id"]
	}
	if parameters["user_agent"] != "" {
		p["useragent"] = parameters["user_agent"]
	}
	if parameters["domain"] != "" {
		p["domain"] = parameters["domain"]
	}
	t := parameters["type"]
	if t != "" {
		p["recaptchatype"] = t
		if t == "4" || t == "5" {
			reqUrl = ENDPOINT_RECAPTCHA_ENTERPRISE
		}
		if t == "5" {
			p["enterprise_type"] = "v3"
		}
	}
	if parameters["v3_action"] != "" {
		p["captchaaction"] = parameters["v3_action"]
	}
	if parameters["v3_min_score"] != "" {
		p["score"] = parameters["v3_min_score"]
	}
	if parameters["data-s"] != "" {
		p["data-s"] = parameters["data-s"]
	}
	if parameters["cookie_input"] != "" {
		p["cookie_input"] = parameters["cookie_input"]
	}

	// make request to API server
	response, err := c.PostRequest(reqUrl, p)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(response, 10, 64)
}

// SubmitCapy - submit capy captcha
func (c *Client) SubmitCapy(parameters map[string]string) (int64, error) {
	reqUrl := ENDPOINT_CAPY
	p := map[string]string{
		"token":       c.AccessToken,
		"action":      "UPLOADCAPTCHA",
		"pageurl":     parameters["page_url"],
		"sitekey":     parameters["sitekey"],
		"captchatype": "12",
	}
	if parameters["proxy"] != "" {
		p["proxy"] = parameters["proxy"]
		p["proxytype"] = "HTTP"
	}
	if parameters["affiliate_id"] != "" {
		p["affiliate_id"] = parameters["affiliate_id"]
	}
	if parameters["user_agent"] != "" {
		p["useragent"] = parameters["user_agent"]
	}

	// make request to API server
	response, err := c.PostRequest(reqUrl, p)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(response, 10, 64)
}

// SubmitFuncaptcha - submit funcaptcha
func (c *Client) SubmitFuncaptcha(parameters map[string]string) (int64, error) {
	reqUrl := ENDPOINT_FUNCAPTCHA
	p := map[string]string{
		"token":       c.AccessToken,
		"action":      "UPLOADCAPTCHA",
		"pageurl":     parameters["page_url"],
		"sitekey":     parameters["sitekey"],
		"captchatype": "13",
	}
	if parameters["data"] != "" {
		p["data"] = parameters["data"]
	}
	if parameters["s_url"] != "" {
		p["surl"] = parameters["s_url"]
	}
	if parameters["user_agent"] != "" {
		p["useragent"] = parameters["user_agent"]
	}
	if parameters["proxy"] != "" {
		p["proxy"] = parameters["proxy"]
		p["proxytype"] = "HTTP"
	}
	if parameters["affiliate_id"] != "" {
		p["affiliate_id"] = parameters["affiliate_id"]
	}
	if parameters["user_agent"] != "" {
		p["useragent"] = parameters["user_agent"]
	}

	// make request to API server
	response, err := c.PostRequest(reqUrl, p)
	if err != nil {
		return 0, err
	}
	response = strings.TrimLeft(response, "[")
	response = strings.TrimRight(response, "]")
	var jsonMap map[string]interface{}
	err = json.Unmarshal([]byte(response), &jsonMap)
	cid := jsonMap["CaptchaId"]
	return strconv.ParseInt(cid.(string), 10, 64)
}

// SubmitGeetest - submit geetest captcha
func (c *Client) SubmitGeetest(parameters map[string]string) (int64, error) {
	reqUrl := ENDPOINT_GEETEST
	p := map[string]string{
		"token":     c.AccessToken,
		"action":    "UPLOADCAPTCHA",
		"domain":    parameters["domain"],
		"challenge": parameters["challenge"],
		"gt":        parameters["gt"],
	}
	if parameters["api_server"] != "" {
		p["api_server"] = parameters["api_server"]
	}
	if parameters["proxy"] != "" {
		p["proxy"] = parameters["proxy"]
		p["proxytype"] = "HTTP"
	}
	if parameters["affiliate_id"] != "" {
		p["affiliate_id"] = parameters["affiliate_id"]
	}
	if parameters["user_agent"] != "" {
		p["useragent"] = parameters["user_agent"]
	}

	// make request to API server
	response, err := c.PostRequest(reqUrl, p)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(response, 10, 64)
}

// SubmitGeetestV4 - submit geetestv4 captcha
func (c *Client) SubmitGeetestV4(parameters map[string]string) (int64, error) {
	reqUrl := ENDPOINT_GEETESTV4
	p := map[string]string{
		"token":     c.AccessToken,
		"action":    "UPLOADCAPTCHA",
		"domain":    parameters["domain"],
		"geetestid": parameters["geetestid"],
	}
	if parameters["proxy"] != "" {
		p["proxy"] = parameters["proxy"]
		p["proxytype"] = "HTTP"
	}
	if parameters["affiliate_id"] != "" {
		p["affiliate_id"] = parameters["affiliate_id"]
	}
	if parameters["user_agent"] != "" {
		p["useragent"] = parameters["user_agent"]
	}

	// make request to API server
	response, err := c.PostRequest(reqUrl, p)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(response, 10, 64)
}

// SubmitHcaptcha - submit hcaptcha
func (c *Client) SubmitHcaptcha(parameters map[string]string) (int64, error) {
	reqUrl := ENDPOINT_HCAPTCHA
	p := map[string]string{
		"token":       c.AccessToken,
		"action":      "UPLOADCAPTCHA",
		"pageurl":     parameters["page_url"],
		"sitekey":     parameters["sitekey"],
		"captchatype": "11",
	}
	if parameters["invisible"] != "" {
		p["invisible"] = "1"
	}
	if parameters["domain"] != "" {
		p["apiEndpoint"] = parameters["domain"]
	}
	if parameters["enterprise_payload"] != "" {
		p["HcaptchaEnterprise"] = parameters["enterprise_payload"]
	}
	if parameters["proxy"] != "" {
		p["proxy"] = parameters["proxy"]
		p["proxytype"] = "HTTP"
	}
	if parameters["affiliate_id"] != "" {
		p["affiliate_id"] = parameters["affiliate_id"]
	}
	if parameters["user_agent"] != "" {
		p["useragent"] = parameters["user_agent"]
	}

	// make request to API server
	response, err := c.PostRequest(reqUrl, p)
	if err != nil {
		return 0, err
	}
	response = strings.TrimLeft(response, "[")
	response = strings.TrimRight(response, "]")
	var jsonMap map[string]interface{}
	err = json.Unmarshal([]byte(response), &jsonMap)
	cid := jsonMap["CaptchaId"]
	return strconv.ParseInt(cid.(string), 10, 64)
}

// SubmitTiktok - submit tiktok captcha
func (c *Client) SubmitTiktok(parameters map[string]string) (int64, error) {
	reqUrl := ENDPOINT_TIKTOK
	p := map[string]string{
		"token":       c.AccessToken,
		"action":      "UPLOADCAPTCHA",
		"pageurl":     parameters["page_url"],
		"captchatype": "10",
	}
	if parameters["cookie_input"] != "" {
		p["cookie_input"] = parameters["cookie_input"]
	}
	if parameters["proxy"] != "" {
		p["proxy"] = parameters["proxy"]
		p["proxytype"] = "HTTP"
	}
	if parameters["affiliate_id"] != "" {
		p["affiliate_id"] = parameters["affiliate_id"]
	}
	if parameters["user_agent"] != "" {
		p["useragent"] = parameters["user_agent"]
	}

	// make request to API server
	response, err := c.PostRequest(reqUrl, p)
	if err != nil {
		return 0, err
	}
	response = strings.TrimLeft(response, "[")
	response = strings.TrimRight(response, "]")
	var jsonMap map[string]interface{}
	err = json.Unmarshal([]byte(response), &jsonMap)
	cid := jsonMap["CaptchaId"]
	return strconv.ParseInt(cid.(string), 10, 64)
}

// SubmitTurnstile - submit turnstile (cloudflare) captcha
func (c *Client) SubmitTurnstile(parameters map[string]string) (int64, error) {
	reqUrl := ENDPOINT_TURNSTILE
	p := map[string]string{
		"token":       c.AccessToken,
		"action":      "UPLOADCAPTCHA",
		"pageurl":     parameters["page_url"],
		"sitekey":     parameters["sitekey"],
		"captchatype": "10",
	}
	if parameters["domain"] != "" {
		p["apiEndpoint"] = parameters["domain"]
	}
	if parameters["action"] != "" {
		p["taction"] = parameters["action"]
	}
	if parameters["cdata"] != "" {
		p["data"] = parameters["cdata"]
	}
	if parameters["proxy"] != "" {
		p["proxy"] = parameters["proxy"]
		p["proxytype"] = "HTTP"
	}
	if parameters["affiliate_id"] != "" {
		p["affiliate_id"] = parameters["affiliate_id"]
	}
	if parameters["user_agent"] != "" {
		p["useragent"] = parameters["user_agent"]
	}

	// make request to API server
	response, err := c.PostRequest(reqUrl, p)
	if err != nil {
		return 0, err
	}
	response = strings.TrimLeft(response, "[")
	response = strings.TrimRight(response, "]")
	var jsonMap map[string]interface{}
	err = json.Unmarshal([]byte(response), &jsonMap)
	cid := jsonMap["CaptchaId"]
	return strconv.ParseInt(cid.(string), 10, 64)
}

// SubmitTask - submit task
func (c *Client) SubmitTask(parameters map[string]string) (int64, error) {
	reqUrl := ENDPOINT_TASK
	p := map[string]string{
		"token":         c.AccessToken,
		"action":        "UPLOADCAPTCHA",
		"pageurl":       parameters["page_url"],
		"template_name": parameters["template_name"],
		"captchatype":   "16",
	}
	if parameters["variables"] != "" {
		p["variables"] = parameters["variables"]
	}
	if parameters["proxy"] != "" {
		p["proxy"] = parameters["proxy"]
		p["proxytype"] = "HTTP"
	}
	if parameters["affiliate_id"] != "" {
		p["affiliate_id"] = parameters["affiliate_id"]
	}
	if parameters["user_agent"] != "" {
		p["useragent"] = parameters["user_agent"]
	}

	// make request to API server
	response, err := c.PostRequest(reqUrl, p)
	if err != nil {
		return 0, err
	}
	response = strings.TrimLeft(response, "[")
	response = strings.TrimRight(response, "]")
	var jsonMap map[string]interface{}
	err = json.Unmarshal([]byte(response), &jsonMap)
	cid := jsonMap["CaptchaId"]
	return strconv.ParseInt(cid.(string), 10, 64)
}

// GetResult - get the captcha solution using the captcha ID returned on submission
func (c *Client) GetResult(captchaId int64) (string, error) {
	parameters := map[string]string{
		"action":    "GETTEXT",
		"captchaid": strconv.FormatInt(captchaId, 10),
		"token":     c.AccessToken,
	}
	response, err := c.PostRequest(ENDPOINT_RESULT, parameters)
	if err != nil {
		return response, err
	}
	if strings.Contains(response, "Pending") {
		return "", err
	}
	return response, nil
}

// Solve - waits for captcha to be solved and returns solution
func (c *Client) Solve(captchaId int64, pollingInterval int) (string, error) {
	counter := 360 / pollingInterval
	for counter > 0 {
		time.Sleep(time.Duration(pollingInterval) * time.Second)

		// get the result
		solution, err := c.GetResult(captchaId)

		if err != nil {
			return "", err
		}

		if solution != "" {
			if strings.Contains(solution, "IMAGE_TIMED_OUT") {
				return "", errors.New("IMAGE TIMED OUT")
			}
			return solution, err
		}

		counter -= 1
	}

	return "", errors.New("captcha could not be solved in time")
}

// PushTaskVariables - sets the captcha as bad captcha
func (c *Client) PushTaskVariables(captchaId int64, pushVariables string) error {
	parameters := make(map[string]string)
	parameters["action"] = "GETTEXT"
	parameters["token"] = c.AccessToken
	parameters["captchaid"] = strconv.FormatInt(captchaId, 10)
	parameters["pushVariables"] = pushVariables

	response, err := c.PostRequest(ENDPOINT_TASK_PUSH_VARIABLES, parameters)
	if err != nil {
		return err
	}

	if !strings.Contains(response, "Push Variables Added") {
		return errors.New(fmt.Sprintf("unknown response when pushing task variables: %s", response))
	}

	return nil
}

// SetCaptchaBad - sets the captcha as bad captcha
func (c *Client) SetCaptchaBad(captchaId int64) error {
	parameters := make(map[string]string)
	parameters["action"] = "SETBADIMAGE"
	parameters["token"] = c.AccessToken
	parameters["imageid"] = strconv.FormatInt(captchaId, 10)
	parameters["submit"] = "Submissssst"

	response, err := c.PostRequest(ENDPOINT_BAD, parameters)
	if err != nil {
		return err
	}

	if response != "SUCCESS" {
		return errors.New(response)
	}

	return nil
}
