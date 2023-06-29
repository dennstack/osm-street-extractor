package main

import (
	"fmt"
	"log"

	"github.com/dschila/osm-street-extractor/internal/csv"
	"github.com/dschila/osm-street-extractor/internal/models"
	"github.com/dschila/osm-street-extractor/internal/osm"
)

var urls = []string{
	"https://download.geofabrik.de/europe/germany/baden-wuerttemberg-latest.osm.bz2",
	"https://download.geofabrik.de/europe/germany/bayern-latest.osm.bz2",
	"https://download.geofabrik.de/europe/germany/berlin-latest.osm.bz2",
	"https://download.geofabrik.de/europe/germany/brandenburg-latest.osm.bz2",
	"https://download.geofabrik.de/europe/germany/bremen-latest.osm.bz2",
	"https://download.geofabrik.de/europe/germany/hamburg-latest.osm.bz2",
	"https://download.geofabrik.de/europe/germany/hessen-latest.osm.bz2",
	"https://download.geofabrik.de/europe/germany/mecklenburg-vorpommern-latest.osm.bz2",
	"https://download.geofabrik.de/europe/germany/niedersachsen-latest.osm.bz2",
	"https://download.geofabrik.de/europe/germany/nordrhein-westfalen-latest.osm.bz2",
	"https://download.geofabrik.de/europe/germany/rheinland-pfalz-latest.osm.bz2",
	"https://download.geofabrik.de/europe/germany/saarland-latest.osm.bz2",
	"https://download.geofabrik.de/europe/germany/sachsen-latest.osm.bz2",
	"https://download.geofabrik.de/europe/germany/sachsen-anhalt-latest.osm.bz2",
	"https://download.geofabrik.de/europe/germany/schleswig-holstein-latest.osm.bz2",
	"https://download.geofabrik.de/europe/germany/thueringen-latest.osm.bz2",
}

func main() {

	writer, err := csv.CreateCSVWriter()
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

	addressStream := make(chan models.Address)
	go csv.WriteAddress(addressStream, writer)

	for _, url := range urls {
		fmt.Printf("Download from %s\n", url)
		osm.ParseFromUrl(url, addressStream)
		if err != nil {
			panic(err)
		}
	}

	close(addressStream)
}
