// getDomain CfV4

package main

import (
	"fmt"
	"log"
	"os"
	"context"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/cloudflare/cloudflare-go/v4/zones"

    util "github.com/prr123/utility/utilLib"
)

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
		option.WithAPIKey("d93d508d3125411dfab702c636a2adad63248"), // defaults to os.LookupEnv("CLOUDFLARE_API_KEY")
		option.WithAPIEmail("azulsoftwarevlc@gmail.com"),               // defaults to os.LookupEnv("CLOUDFLARE_EMAIL")
	)

	ctx := context.Background()
	query := zones.ZoneListParams{Name:cloudflare.F(domStr)}

	//  zones.ZoneListParams
	pgar, err := client.Zones.List(ctx, query)
	if err != nil {log.Fatalf("error -- Zones List: %v\n", err)}

	fmt.Printf("pg array: %d\n", len(pgar.Result))

	for i, zon := range pgar.Result {
		fmt.Printf("**** zone %d ****\n", i)
		fmt.Printf("name: %s\n", zon.Name)
		fmt.Printf("id:   %s\n", zon.ID)
		fmt.Printf("**** end zone ****\n")
	}
}

