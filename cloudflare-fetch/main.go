package cloudflarefetch

import (
	sr "ddns-updater/send-request"
	"fmt"
	"log"
)

func FetchCloudFlareRecord() (sr.Cloudflare_DNSRECORDS, error) {

	fmt.Println("Finding Zone Details..")
	zone, err := sr.GetZoneId()
	if err != nil {
		fmt.Println("[!] Error:", err)
		return sr.Cloudflare_DNSRECORDS{}, err
	}
	if len(zone.Result) == 0 {
		log.Fatal("[!] No Cloudflare zones detected.", err)
		return sr.Cloudflare_DNSRECORDS{}, err
	}

	target_zone := zone.Result[0] // We always want the first zone

	fmt.Println("\t[+] Zone ID:", target_zone.ID)
	fmt.Println("\t[+] Zone Name:", target_zone.Name)

	fmt.Println("Fetching Dns Records...")

	records, err := sr.ListDnsRecords(target_zone.ID)
	if err != nil {
		fmt.Println("Error", err)
		return sr.Cloudflare_DNSRECORDS{}, err
	}

	for i := 0; i < len(records.Result); i++ {
		fmt.Println("[+] Record ", i+1)
		fmt.Println("\t[+] ID:", records.Result[i].ID)
		fmt.Println("\t[+] Name:", records.Result[i].Name)
		fmt.Println("\t[+] Type:", records.Result[i].Type)
		fmt.Println("\t[+] IP:", records.Result[i].Content)
	}

	return records, nil
}
