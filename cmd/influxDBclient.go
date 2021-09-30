package main

import (
	"context"
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	// Create a client
	// You can generate a Token from the "Tokens Tab" in the UI
	client := influxdb2.NewClient("https://europe-west1-1.gcp.cloud2.influxdata.com", "0H50vNcIocetKxlsmX86RuQHgGWU4O4gFfTFDMRLXDhxj-3LwrraHcMKkd3tq27ahf_s2lnCroTauXJIu737xg==")

	// get non-blocking write client
	writeAPI := client.WriteAPI("raphapainterr@gmail.com", "raphapainterr's Bucket")

	// write line protocol
	writeAPI.WriteRecord(fmt.Sprintf("stat, unit=temperature avg=%f,max=%f", 23.5, 45.0))
	// Flush writes
	writeAPI.Flush()

	// Get query client
	queryAPI := client.QueryAPI("raphapainterr@gmail.com")

	query := `from(bucket:"raphapainterr's Bucket")|> range(start: -1h) |> filter(fn: (r) => r._measurement == "stat")`

	// get QueryTableResult
	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		panic(err)
	}

	// Iterate over query response
	for result.Next() {
		// Notice when group key has changed
		if result.TableChanged() {
			fmt.Printf("table: %s\n", result.TableMetadata().String())
		}
		// Access data
		fmt.Printf("value: %v\n", result.Record().Value())
	}
	// check for an error
	if result.Err() != nil {
		fmt.Printf("query parsing error: %\n", result.Err().Error())
	}

	// always close client at the end
	defer client.Close()
}
