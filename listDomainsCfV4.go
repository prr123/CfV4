// test CFV4

package main

import (
	"fmt"
	"log"
	"context"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/cloudflare/cloudflare-go/v4/zones"
)

func main() {
	client := cloudflare.NewClient(
		option.WithAPIKey("d93d508d3125411dfab702c636a2adad63248"), // defaults to os.LookupEnv("CLOUDFLARE_API_KEY")
		option.WithAPIEmail("azulsoftwarevlc@gmail.com"),               // defaults to os.LookupEnv("CLOUDFLARE_EMAIL")
	)

	ctx := context.Background()
	query := zones.ZoneListParams{}
	//  zones.ZoneListParams
	pgar, err := client.Zones.List(ctx, query)
//(pagination.V4PagePaginationArray[zones.Zone], error)
	if err != nil {log.Fatalf("error -- Zones List: %v\n", err)}

	fmt.Printf("pg array: %d\n", len(pgar.Result))

	for i, zon := range pgar.Result {
		fmt.Printf("**** zone %d ****\n", i)
		fmt.Printf("name: %s\n", zon.Name)
		fmt.Printf("id:   %s\n", zon.ID)
		fmt.Printf("**** end zone ****\n")
	}
}

