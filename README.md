# cfV4
cloudflare go library version 4.2  
https://pkg.go.dev/github.com/cloudflare/cloudflare-go/v4  

The methods have changed, so I had to rewrite the code.

Basic Methods

## listDomainsCfV4
Programs that lists all domains for an account.  

## getDomainCfV4
Program that checks a whether a domain name is served by cloudflare dns servers.  

## getDomainCfV4save
Same as above and it save a yaml file in the cloud domain directory with the zone id.
Zone in CF speak are domains.  

## addDNSTxtRecCfV4
Adds Dns Text record to a zone (domain).  

## listDNSTxtRecCfV4
Lists all Dns Txt records of a domain.  

## delDNSTxtRecCfV4
Deletes a Dns Txt record with a specified name.  

