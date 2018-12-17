package southxchange

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"
)

type client struct {
	apiKey      string
	apiSecret   string
	httpClient  *http.Client
	httpTimeout time.Duration
	useragent   string
	debug       bool
}

// NewClient return a new SouthXchange HTTP client
func NewClient(apiKey, apiSecret, userAgent string) (c *client) {
	return &client{apiKey, apiSecret, &http.Client{}, 30 * time.Second, userAgent, false}
}

// NewClientWithCustomHttpConfig returns a new SouthXchange HTTP client using the predefined http client
func NewClientWithCustomHttpConfig(apiKey, apiSecret, userAgent string, httpClient *http.Client) (c *client) {
	timeout := httpClient.Timeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &client{apiKey, apiSecret, httpClient, timeout, userAgent, false}
}

// NewClient returns a new SouthXchange HTTP client with custom timeout
func NewClientWithCustomTimeout(apiKey, apiSecret, userAgent string, timeout time.Duration) (c *client) {
	return &client{apiKey, apiSecret, &http.Client{}, timeout, userAgent, false}
}

func (c client) dumpRequest(r *http.Request) {
	if r == nil {
		log.Print("dumpReq ok: <nil>")
		return
	}
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Print("dumpReq err:", err)
	} else {
		log.Print("dumpReq ok:", string(dump))
	}
}

func (c client) dumpResponse(r *http.Response) {
	if r == nil {
		log.Print("dumpResponse ok: <nil>")
		return
	}
	dump, err := httputil.DumpResponse(r, true)
	if err != nil {
		log.Print("dumpResponse err:", err)
	} else {
		log.Print("dumpResponse ok:", string(dump))
	}
}

func DecodeData(body []byte, indexes []int) (ret []byte, err error) {
	ret = make([]byte, hex.DecodedLen(indexes[3]-indexes[2]))
	_, err = hex.Decode(ret, body[indexes[2]:indexes[3]])
	return ret, err
}

// doTimeoutRequest do a HTTP request with timeout
func (c *client) doTimeoutRequest(timer *time.Timer, req *http.Request) (*http.Response, error) {
	// Do the request in the background so we can check the timeout
	type result struct {
		resp *http.Response
		err  error
	}
	if c.useragent != "" {
		req.Header.Set("User-Agent", c.useragent)
	}

	done := make(chan result, 1)
	go func() {
		if c.debug {
			c.dumpRequest(req)
		}
		resp, err := c.httpClient.Do(req)
		if err != nil {
			done <- result{nil, err}
			return
		}
		if c.debug {
			c.dumpResponse(resp)
		}
		done <- result{resp, err}
	}()
	// Wait for the read or the timeout
	select {
	case r := <-done:
		return r.resp, r.err
	case <-timer.C:
		return nil, errors.New("timeout on reading data from SouthXchange API")
	}
}

// do prepare and process HTTP request to SouthXchange API
func (c *client) do(method string, ressource string, payload map[string]string, authNeeded bool) (response []byte, err error) {
	connectTimer := time.NewTimer(c.httpTimeout)

	var rawurl string
	if strings.HasPrefix(ressource, "http") {
		rawurl = ressource
	} else {
		rawurl = fmt.Sprintf("%s/%s", API_BASE, ressource)
	}
	if payload == nil {
		payload = make(map[string]string)
	}
	formData := []byte("")
	if method == "POST" {
		payload["key"] = c.apiKey
		payload["nonce"] = strconv.FormatInt(time.Now().UnixNano(), 10)
		formData, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, rawurl, strings.NewReader(string(formData)))
	if err != nil {
		return
	}
	if method == "POST" || method == "PUT" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
	}
	//req.Header.Add("Accept", "application/json") // cloudflare protected api doesnt accept this, i got captcha page
	req.Header.Add("Accept", "*/*")

	// Auth
	if authNeeded {
		if len(c.apiKey) == 0 || len(c.apiSecret) == 0 {
			err = errors.New("You need to set API Key and API Secret to call this method")
			return
		}
		mac := hmac.New(sha512.New, []byte(c.apiSecret))
		_, err = mac.Write(formData)
		if err != nil {
			return
		}
		sig := hex.EncodeToString(mac.Sum(nil))
		req.Header.Add("Hash", sig)
	}

	resp, err := c.doTimeoutRequest(connectTimer, req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	/*
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}
	response, err = ioutil.ReadAll(reader)
	*/
	response, err = ioutil.ReadAll(resp.Body)
	//fmt.Println(fmt.Sprintf("reponse %s", response), err)
	if err != nil {
		return response, err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 401 {
		err = errors.New(resp.Status + ": "+strings.Trim(string(response), "\""))
	}
	return response, err
}
