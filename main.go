package main

import (
	"compress/bzip2"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Node struct {
	Tags []Tag `xml:"tag"`
}

type Tag struct {
	Key   string `xml:"k,attr"`
	Value string `xml:"v,attr"`
}

type Address struct {
	Street   string
	City     string
	Postcode string
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

func addressFromTags(tags []Tag) Address {
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
	return Address{Street: street, City: city, Postcode: postcode}
}

func contains(arr []Address, a Address) bool {
	for _, e := range arr {
		if e == a {
			return true
		}
	}
	return false
}

func convertOSMFile(url, csvFileName string) error {
	csvFile, err := os.Create(csvFileName)
	if err != nil {
		log.Fatalf("failed creating csv file: %s", err)
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	csvWriter.Comma = ';'
	err = csvWriter.Write([]string{"street", "postcode", "city"})
	if err != nil {
		return err
	}

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	bzipReader := bzip2.NewReader(response.Body)
	xmlDecoder := xml.NewDecoder(bzipReader)

	addresses := []Address{}

	// https://eli.thegreenplace.net/2019/faster-xml-stream-processing-in-go/
	for {
		tok, err := xmlDecoder.Token()
		if tok == nil || err == io.EOF {
			// EOF means we're done.
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
					return err
				}
				if hasAddress(node.Tags) {
					address := addressFromTags(node.Tags)
					if !contains(addresses, address) {
						addresses = append(addresses, address)
						err = csvWriter.Write([]string{address.Street, address.Postcode, address.City})
						if err != nil {
							return err
						}
					}
				}
			}
		default:
		}
	}

	csvWriter.Flush()
	return nil
}

func main() {
	var url, destination string
	if len(os.Args) > 2 {
		url = os.Args[1]
		destination = os.Args[2]
	} else {
		log.Fatal("Missing parameters")
	}

	fmt.Printf("Download from %s\n", url)
	err := convertOSMFile(url, destination)
	if err != nil {
		panic(err)
	}
	fmt.Println("Save file: " + destination)
}
