package osm

import (
	"fmt"
	"io"
	"net/http"

	"github.com/qedus/osmpbf"
)

type Address struct {
	Street   string
	City     string
	Postcode string
}

func hasAddressTags(tags map[string]string) bool {
	hasCity := false
	hasStreet := false
	hasPostcode := false

	for key := range tags {
		switch key {
		case "addr:city":
			hasCity = true
		case "addr:street":
			hasStreet = true
		case "addr:postcode":
			hasPostcode = true
		}
	}
	return hasCity && hasStreet && hasPostcode
}

func addressFromOSMTags(tags map[string]string) Address {
	return Address{
		Street:   tags["addr:street"],
		City:     tags["addr:city"],
		Postcode: tags["addr:postcode"],
	}
}

func ParseFromUrl(url string, result chan<- Address) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error downloading file: %v\n", err)
		return
	}
	defer response.Body.Close()

	d := osmpbf.NewDecoder(response.Body)
	err = d.Start(4) // 4 worker goroutines
	if err != nil {
		fmt.Printf("Error starting decoder: %v\n", err)
		return
	}

	addressMap := make(map[Address]bool)

	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("Error decoding: %v\n", err)
			break
		} else {
			switch obj := v.(type) {
			case *osmpbf.Node:
				if hasAddressTags(obj.Tags) {
					address := addressFromOSMTags(obj.Tags)
					if _, exists := addressMap[address]; !exists {
						addressMap[address] = true
						result <- address
					}
				}
			case *osmpbf.Way:
				if hasAddressTags(obj.Tags) {
					address := addressFromOSMTags(obj.Tags)
					if _, exists := addressMap[address]; !exists {
						addressMap[address] = true
						result <- address
					}
				}
			}
		}
	}
}
