package main

import (
	cf "ddns-updater/cloudflare-fetch"
	sr "ddns-updater/send-request"
	"fmt"
	"log"
	"time"
)

func mainLoop() {

	fmt.Println("Fetching Cloudflare for the initial record.")
	var CF_Record sr.Cloudflare_DNSRECORDS
	CF_Record, err := cf.FetchCloudFlareRecord()
	if err != nil {
		log.Fatal("Failed to fetch Cloudflare.", err)
	}

	for {
		fmt.Println("\nFetching public IP....")
		public_ip, err := sr.GetPublicIp()
		if err != nil {
			return
		}
		fmt.Println("\t [+] Public IP:", public_ip)

		if public_ip != CF_Record.Result[0].Content {
			fmt.Println("[!] IP mismatch. Updating Cloudflare...")
			record := CF_Record.Result[0]
			new_records, err := sr.OverwritteDnsrecords(record.ID, record.Name, public_ip)
			if err != nil {
				fmt.Println("\t[!] ", err)
				// Uncoment below is optional for the program to exit
				//return
				fmt.Println("[!] Sleeping for 10 minutes...")
				time.Sleep(10 * time.Minute)
			} else {
				fmt.Println("[+] Record updated successfully!")
				fmt.Println("\t[+] ID:", new_records.Result.ID)
				fmt.Println("\t[+] Name:", new_records.Result.Name)
				fmt.Println("\t[+] Type:", new_records.Result.Type)
				fmt.Println("\t[+] Content:", new_records.Result.Content)

				// update our original records with the new ones
				CF_Record.Result[0] = new_records.Result
			}
		} else {
			fmt.Println("[+] IPs match. Sleeping for 10 minutes...")
			time.Sleep(10 * time.Minute)
		}
	}
}

func main() {

	mainLoop()
}
