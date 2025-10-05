package osm

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

func CreateCSVWriter() (*csv.Writer, error) {
	now := time.Now()
	fn := "osm-germany-" + now.Format("2006-01-02") + ".csv"
	file, err := os.Create(fn)
	if err != nil {
		return nil, err
	}

	writer := csv.NewWriter(file)
	return writer, nil
}

func WriteAddress(stream <-chan Address, writer *csv.Writer) {
	for addr := range stream {
		err := writer.Write([]string{addr.Street, addr.Postcode, addr.City})
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}
