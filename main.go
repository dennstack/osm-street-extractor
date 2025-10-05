package main

import (
	"fmt"
	"log"

	"github.com/dschila/osm-street-extractor/osm"
)

var urls = []string{
	"https://download.geofabrik.de/europe/germany/baden-wuerttemberg-latest.osm.pbf",
	"https://download.geofabrik.de/europe/germany/bayern-latest.osm.pbf",
	"https://download.geofabrik.de/europe/germany/berlin-latest.osm.pbf",
	"https://download.geofabrik.de/europe/germany/brandenburg-latest.osm.pbf",
	"https://download.geofabrik.de/europe/germany/bremen-latest.osm.pbf",
	"https://download.geofabrik.de/europe/germany/hamburg-latest.osm.pbf",
	"https://download.geofabrik.de/europe/germany/hessen-latest.osm.pbf",
	"https://download.geofabrik.de/europe/germany/mecklenburg-vorpommern-latest.osm.pbf",
	"https://download.geofabrik.de/europe/germany/niedersachsen-latest.osm.pbf",
	"https://download.geofabrik.de/europe/germany/nordrhein-westfalen-latest.osm.pbf",
	"https://download.geofabrik.de/europe/germany/rheinland-pfalz-latest.osm.pbf",
	"https://download.geofabrik.de/europe/germany/saarland-latest.osm.pbf",
	"https://download.geofabrik.de/europe/germany/sachsen-latest.osm.pbf",
	"https://download.geofabrik.de/europe/germany/sachsen-anhalt-latest.osm.pbf",
	"https://download.geofabrik.de/europe/germany/schleswig-holstein-latest.osm.pbf",
	"https://download.geofabrik.de/europe/germany/thueringen-latest.osm.pbf",
}

func main() {

	writer, err := osm.CreateCSVWriter()
	if err != nil {
		fmt.Println("create csv file failed:", err)
		return
	}
	defer writer.Flush()

	writer.Comma = ';'
	err = writer.Write([]string{"street", "postcode", "city"})
	if err != nil {
		log.Fatal(err)
	}

	addressStream := make(chan osm.Address)
	go osm.WriteAddress(addressStream, writer)

	for _, url := range urls {
		fmt.Printf("Download from %s\n", url)
		osm.ParseFromUrl(url, addressStream)
		if err != nil {
			panic(err)
		}
	}

	close(addressStream)
}
