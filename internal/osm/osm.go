package osm

import (
	"compress/bzip2"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/dschila/osm-street-extractor/internal/models"
)

type Node struct {
	Tags []Tag `xml:"tag"`
}

type Tag struct {
	Key   string `xml:"k,attr"`
	Value string `xml:"v,attr"`
}

func hasAddress(tags []Tag) bool {
	hasCity := false
	hasStreet := false
	hasPostcode := false

	for _, tag := range tags {
		switch tag.Key {
		case "addr:city":
			hasCity = true
		case "addr:street":
			hasStreet = true
		case "addr:postcode":
			hasPostcode = true
		default:
		}
	}
	return hasCity && hasStreet && hasPostcode
}

func addressFromTags(tags []Tag) models.Address {
	var street, postcode, city string
	for _, tag := range tags {
		switch tag.Key {
		case "addr:city":
			city = tag.Value
		case "addr:street":
			street = tag.Value
		case "addr:postcode":
			postcode = tag.Value
		default:
		}
	}
	return models.Address{Street: street, City: city, Postcode: postcode}
}

func contains(arr []models.Address, a models.Address) bool {
	for _, e := range arr {
		if e == a {
			return true
		}
	}
	return false
}

func ParseFromUrl(url string, result chan<- models.Address) {

	response, err := http.Get(url)
	if err != nil {
		return
	}
	defer response.Body.Close()

	bzipReader := bzip2.NewReader(response.Body)
	xmlDecoder := xml.NewDecoder(bzipReader)

	addresses := []models.Address{}
	for {
		tok, err := xmlDecoder.Token()
		if tok == nil || err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Error decoding token: %s", err)
		}

		switch ty := tok.(type) {
		case xml.StartElement:
			if ty.Name.Local == "node" || ty.Name.Local == "way" {
				var node Node
				err := xmlDecoder.DecodeElement(&node, &ty)
				if err != nil {
					fmt.Println("Error: ", err)
				}
				if hasAddress(node.Tags) {
					address := addressFromTags(node.Tags)
					if !contains(addresses, address) {
						addresses = append(addresses, address)
						result <- address
					}
				}
			}
		default:
		}
	}
}
