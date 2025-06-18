package sendrequest

import (
	"bytes"
	"ddns-updater/readenv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Cloudflare_ZONE struct {
	Result []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"result"`
}

type Cloudflare_DNSRECORDS struct {
	Result []struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Type    string `json:"type"`
		Content string `json:"content"`
	}
	Success bool `json:"success"`
}

type CloudFlare_UPDATERECORDS_errors struct {
	Code    uint32 `json:"code"`
	Message string `json:"message"`
}

type CloudFlare_UPDATERECORDS struct {
	Result struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Type    string `json:"type"`
		Content string `json:"content"`
	}
	Success bool                              `json:"success"`
	Errors  []CloudFlare_UPDATERECORDS_errors `json:"errors"`
}

type MyIpInfo_err struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type MyIpInfo struct {
	Ip     string       `json:"ip"`
	Status int16        `json:"status"`
	Error  MyIpInfo_err `json:"error"`
}

func SendPutRequest(url string, body []byte) (resp *http.Response, err error) {
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Request creation failed:", err)
		return nil, err
	}

	key, email := readenv.ReadEnvKeys()

	req.Header.Add("X-Auth-Email", email)
	req.Header.Add("X-Auth-Key", key)
	req.Header.Add("Content-type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}

	return resp, nil
}

func SendGetRequest(url string, sendHeaders bool) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Request creation failed:", err)
		return nil, err
	}

	if sendHeaders {
		key, email := readenv.ReadEnvKeys()
		req.Header.Add("X-Auth-Email", email)
		req.Header.Add("X-Auth-Key", key)
		req.Header.Add("Content-type", "application/json")
	}

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}

	return resp, nil
}

func GetZoneId() (Cloudflare_ZONE, error) {

	url := "https://api.cloudflare.com/client/v4/zones"

	resp, err := SendGetRequest(url, true)
	if err != nil {
		return Cloudflare_ZONE{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return Cloudflare_ZONE{}, err
	}

	var data Cloudflare_ZONE
	if err := json.Unmarshal(body, &data); err != nil {
		return Cloudflare_ZONE{}, err
	}

	if len(data.Result) == 0 {
		return Cloudflare_ZONE{}, fmt.Errorf("no zones found")
	}

	return data, nil
}

func ListDnsRecords(zone_id string) (Cloudflare_DNSRECORDS, error) {

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zone_id)

	resp, err := SendGetRequest(url, true)
	if err != nil {
		return Cloudflare_DNSRECORDS{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return Cloudflare_DNSRECORDS{}, err
	}

	var data Cloudflare_DNSRECORDS
	if err := json.Unmarshal(body, &data); err != nil {
		return Cloudflare_DNSRECORDS{}, err
	}

	return data, nil
}

func GetPublicIp() (string, error) {
	url := "https://ipinfo.io"

	resp, err := SendGetRequest(url, false)
	if err != nil {
		fmt.Println("Request creation failed:", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", err
	}

	// // Debug helper
	// fmt.Println("Raw IP Info response:", string(body))

	var data MyIpInfo
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("JSON unmarshal error:", err)
		return "", err
	}

	// If we rate limit on one API, try the alternative
	if data.Status == 429 {
		// fmt.Printf("[!] Rate limit exceeded: %s - %s\n", data.Error.Title, data.Error.Message)
		fmt.Println("[!] Rate limit exceeded on ipinfo.io, using alternative provider!")
		url := "https://api64.ipify.org?format=json"
		resp, err := SendGetRequest(url, false)
		if err != nil {
			fmt.Println("Request creation failed:", err)
			return "", err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return "", err
		}

		var data MyIpInfo
		if err := json.Unmarshal(body, &data); err != nil {
			fmt.Println("JSON unmarshal error:", err)
			return "", err
		}

		return data.Ip, nil
	}

	return data.Ip, nil
}

func OverwritteDnsrecords(dns_record_id, dns_record_name, new_ip string) (CloudFlare_UPDATERECORDS, error) {

	zones, err := GetZoneId()
	if err != nil {
		log.Fatal("Failed to fetch zone", err)
		return CloudFlare_UPDATERECORDS{}, err
	}
	zone := zones.Result[0]

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zone.ID, dns_record_id)

	var json_body = []byte(`{
		"comment": "Domain verification record",
		"content": "` + new_ip + `",
		"name": "` + dns_record_name + `",
		"proxied": false,
		"type": "A"
	  }`)

	resp, err := SendPutRequest(url, json_body)
	if err != nil {
		fmt.Println("Request failed", err)
		return CloudFlare_UPDATERECORDS{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return CloudFlare_UPDATERECORDS{}, err
	}

	var data CloudFlare_UPDATERECORDS
	if err := json.Unmarshal(body, &data); err != nil {
		return CloudFlare_UPDATERECORDS{}, err
	}
	success := data.Success

	if !success {
		fmt.Println("Records did not update, listing errors:")
		return CloudFlare_UPDATERECORDS{}, errors.New(data.Errors[0].Message)
	}

	return data, nil

}
