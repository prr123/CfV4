// addDnsRecord CfV4

package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"context"
	"strings"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/cloudflare/cloudflare-go/v4/zones"
	"github.com/cloudflare/cloudflare-go/v4/dns"

	"github.com/goccy/go-json"
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
	log.Printf("found domain dns rec: %s!\n", domStr)

	if dbg {
		for i, zon := range pgar.Result {
			fmt.Printf("**** zone %d ****\n", i)
			fmt.Printf("name: %s\n", zon.Name)
			fmt.Printf("id:   %s\n", zon.ID)
			fmt.Printf("**** end zone ****\n")
		}
	}

	zon := pgar.Result[0]
	// create dns record
	txtrec := dns.TXTRecordParam {
		Name: cloudflare.F("azultest"),
		Content: cloudflare.F("\"abcTest\""),
		Type: cloudflare.F(dns.TXTRecordTypeTXT),
	}

	params := dns.RecordNewParams {
		ZoneID: cloudflare.F(zon.ID),
		Record: txtrec,
	}
//	client.DNS.Records.New(ctx context.Context, params dns.RecordNewParams) (dns.RecordResponse, error)
	recresp, err := client.DNS.Records.New(ctx, params)
	if err !=nil {
		//parse error
		errStr := err.Error()
//		if dbg {fmt.Printf("error info -- %s\n", errStr)}
		idx := strings.Index(errStr, "{")
		jsonStr := string(errStr[idx:])
//		fmt.Printf("error json info -- %s\n", jsonStr)
		jsonErrObj := CfErrObj{}
		jserr := json.Unmarshal([]byte(jsonStr),&jsonErrObj)
 		if jserr != nil {fmt.Printf("error -- json conv: %v\n",jserr)}

		fmt.Println("***** add Dns Record Error *****")
		for i:=0; i < len(jsonErrObj.Errors); i++ {
			fmt.Printf("error[%d] code: %d message: %s\n", i+1, jsonErrObj.Errors[i].Code, jsonErrObj.Errors[i].Msg)
		}
		fmt.Println("*** End add Dns Record Error ***")
//		fmt.Printf("info -- jsonErrObj: %v\n", jsonErrObj)
		log.Fatalf("error -- new dns rec")
	}

	fmt.Printf("info -- recresp id: %s\n", recresp.ID)
	tim:= recresp.CreatedOn
	fmt.Printf("info -- recresp Created: %s\n", tim.Format(time.RFC1123))

}

/*

type RecordNewParams struct {
	// Identifier
	ZoneID param.Field[string] `path:"zone_id,required"`
	Record RecordUnionParam    `json:"record,required"`
}

type TXTRecordParam struct {
	// Comments or notes about the DNS record. This field has no effect on DNS
	// responses.
	Comment param.Field[string] `json:"comment"`
	// Text content for the record. The content must consist of quoted "character
	// strings" (RFC 1035), each with a length of up to 255 bytes. Strings exceeding
	// this allowed maximum length are automatically split.
	//
	// Learn more at
	// <https://www.cloudflare.com/learning/dns/dns-records/dns-txt-record/>.
	Content param.Field[string] `json:"content"`
	// DNS record name (or @ for the zone apex) in Punycode.
	Name param.Field[string] `json:"name"`
	// Whether the record is receiving the performance and security benefits of
	// Cloudflare.
	Proxied param.Field[bool] `json:"proxied"`
	// Settings for the DNS record.
	Settings param.Field[TXTRecordSettingsParam] `json:"settings"`
	// Custom tags for the DNS record. This field has no effect on DNS responses.
	Tags param.Field[[]RecordTagsParam] `json:"tags"`
	// Time To Live (TTL) of the DNS record in seconds. Setting to 1 means 'automatic'.
	// Value must be between 60 and 86400, with the minimum reduced to 30 for
	// Enterprise zones.
	TTL param.Field[TTL] `json:"ttl"`
	// Record type.
	Type param.Field[TXTRecordType] `json:"type"`
}

*/
