// addDnsRecord CfV4
// https://pkg.go.dev/github.com/cloudflare/cloudflare-go/v4@v4.2.0/dns#RecordResponse

package main

import (
	"fmt"
	"log"
	"os"
//	"time"
	"context"
//	"strings"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/cloudflare/cloudflare-go/v4/zones"
	"github.com/cloudflare/cloudflare-go/v4/dns"

//	"github.com/goccy/go-json"
    util "github.com/prr123/utility/utilLib"
)

type CfErrors struct {
	Code int `json:"code"`
	Msg string `json:"message"`
	Msgs []string `json:"messages"`
}

type CfErrObj struct {
	Result *string `json:"result"`
	Success bool `json:"success"`
	Errors []CfErrors `json:"errors"`
}



func main() {

    numArgs := len(os.Args)

    flags:=[]string{"dbg","domain"}

    useStr := "/domain=name [/dbg]"
    helpStr := fmt.Sprintf("help: The program retrieves zones name and id from cloudflare\n")

    if numArgs > len(flags)+1 {
        fmt.Println("too many arguments in cl!")
        fmt.Println("usage: %s %s\n", os.Args[0], useStr)
        os.Exit(1)
    }


    if numArgs == 2 {
        if os.Args[1] == "help" {
            fmt.Printf("usage is: %s\n", useStr)
            fmt.Printf("%s\n", helpStr)
            os.Exit(1)
        }
    }

    flagMap, err := util.ParseFlags(os.Args, flags)
    if err != nil {log.Fatalf("util.ParseFlags: %v\n", err)}

    dbg := false
    _, ok := flagMap["dbg"]
    if ok {dbg = true}

    domStr := ""
    dval, ok := flagMap["domain"]
    if ok {
        if dval.(string) == "none" {log.Fatalf("error -- no domain name provided with /domain flag!")}
        domStr = dval.(string)
    } else {
		log.Fatalf("error -- no domain flag provided!")
	}

	if dbg {
		fmt.Printf("debug -- domain: %s\n", domStr)
	}

	client := cloudflare.NewClient(
		option.WithAPIKey("api key"), // defaults to os.LookupEnv("CLOUDFLARE_API_KEY")
		option.WithAPIEmail("azulsoftwarevlc@gmail.com"),               // defaults to os.LookupEnv("CLOUDFLARE_EMAIL")
	)

	ctx := context.Background()
	query := zones.ZoneListParams{Name:cloudflare.F(domStr)}

	//  zones.ZoneListParams
	pgar, err := client.Zones.List(ctx, query)
	if err != nil {log.Fatalf("error -- Zones List: %v\n", err)}

	fmt.Printf("pg array: %d\n", len(pgar.Result))
	if len(pgar.Result) <1 {log.Fatalf("error -- no zone found!")}

	for i, zon := range pgar.Result {
		fmt.Printf("**** zone %d ****\n", i)
		fmt.Printf("name: %s\n", zon.Name)
		fmt.Printf("id:   %s\n", zon.ID)
		fmt.Printf("**** end zone ****\n")
	}

	zon := pgar.Result[0]


	// first retrieve all dns TXT records
	// func (r *RecordService) List(ctx context.Context, params RecordListParams, opts ...option.RequestOption) (res *pagination.V4PagePaginationArray[RecordResponse], err error)
//	dnsQuery := zones.ZoneListParams{Name:cloudflare.F(domStr)}

	qlParams := dns.RecordListParams {
		ZoneID: cloudflare.F(zon.ID),
		Type:  cloudflare.F(dns.RecordListParamsTypeTXT),
	}

	recpgar, err := client.DNS.Records.List(ctx, qlParams)
	if err != nil {log.Fatalf("error -- list DnsRecs: %v\n", err)}

//	recs := recpgar.Result[0]

	fmt.Printf("pg array: %d\n", len(recpgar.Result))
	    for i, rec := range recpgar.Result {
        fmt.Printf("**** txt Record %d ****\n", i)
        fmt.Printf("id:   %s\n", rec.ID)
        fmt.Printf("name: %s\n", rec.Name)
        fmt.Printf("content: %s\n", rec.Content)
        fmt.Printf("**** end zone ****\n")
    }

}

