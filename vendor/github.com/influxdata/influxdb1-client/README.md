# influxdb1-clientv2
influxdb1-clientv2 is the current Go client API for InfluxDB 1.x. A Go client for the 2.0 API will be coming soon.

InfluxDB is an open-source distributed time series database, find more about [InfluxDB](https://www.influxdata.com/time-series-platform/influxdb/) at https://docs.influxdata.com/influxdb/latest

## Usage
To import into your Go project, run the following command in your terminal:
`go get github.com/influxdata/influxdb1-client/v2`
Then, in your import declaration section of your Go file, paste the following:
`import "github.com/influxdata/influxdb1-client/v2"`

If you get the error `build github.com/user/influx: cannot find module for path github.com/influxdata/influxdb1-client/v2` when trying to build:
change your import to:
```go
import(
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
)
```

## Example
The following example creates a new client to the InfluxDB host on localhost:8086 and runs a query for the measurement `cpu_load` from the `mydb` database. 
``` go
func ExampleClient_query() {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		fmt.Println("Error creating InfluxDB Client: ", err.Error())
	}
	defer c.Close()

	q := client.NewQuery("SELECT count(value) FROM cpu_load", "mydb", "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		fmt.Println(response.Results)
	}
}
```
