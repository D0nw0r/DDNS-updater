# DDNS-updater

This project is to facilitate Dynamic DNS by updating CloudFlare's DNS A records periodically with the public IP address found.

Built with Golang, it is cross-platform compatible by default.

It fetches the current public IP and compares it to CloudFlare first A record (record[0]), if they are different it will update the record accordingly.

# Requirements
- Create an .env file on the "readenv" folder with the following format:

```
EMAIL="my_email@email.com"
KEY=<your_cloudflare_key>
```

# How it works

DDNS-Updater fetches cloudflare with your API Key and Email on 3 endpoints:

- https://api.cloudflare.com/client/v4/zones 
- https://api.cloudflare.com/client/v4/zones/ZONE_ID/dns_records
- https://api.cloudflare.com/client/v4/zones/ZONE_ID/dns_records/RECORD_ID

The first 2 are to obtain the current IP address pointed by the record. The last one is to update it.

The program runs a loop of fetching the Public IP on ipinfo.io and comparing it to the stored records of cloudflare. 
When there is a mismatch, a request for change is made.


# How to run

Feel free to run the program in any way such as simply building and running with go:
```
go build -o ddns-updater
./ddns-updater
``` 
The program will run on a loop until cancelled by user interaction.

## Docker support

Running as a container is also possible with the included docker file.

```
docker build -t ddns-updater .
docker run ddns-updater
```

A `docker-compose` file is also provided for those who prefer it.

# TODO
- Support for the user to select which record (default just uses record[0])
- Db integration?

