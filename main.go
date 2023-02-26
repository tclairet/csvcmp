package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	nameOri := os.Args[1]
	fileOri, err := os.Open(nameOri)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	nameNew := os.Args[2]
	fileNew, err := os.Open(nameNew)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	rOri := csv.NewReader(fileOri)
	rNew := csv.NewReader(fileNew)

	rOri.Comma = ';'
	rNew.Comma = ';'

	res := Check(rOri, rNew)

	file, err := os.Create(fmt.Sprintf("from_%s_to_%s.csv", strings.TrimSuffix(nameOri, ".csv"), strings.TrimSuffix(nameNew, ".csv")))
	defer file.Close()
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	w := csv.NewWriter(file)
	defer w.Flush()

	var data [][]string
	for _, record := range res {
		row := []string{record}
		data = append(data, row)
	}
	w.WriteAll(data)
}

func Check(original *csv.Reader, newVersion *csv.Reader) []string {
	old := make(map[string]record)
	iterate(original, makeCompareFunc(old))

	new := make(map[string]record)
	iterate(newVersion, makeCompareFunc(new))

	var newRDV []string
	for key, value := range new {
		if _, exist := old[key]; exist {
			continue
		}
		newRDV = append(newRDV, value.student)
	}
	return newRDV
}

type Reader interface {
	Read() (record []string, err error)
}

func iterate(reader Reader, fn func([]string)) {
	for {
		records, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fn(records)
	}
}

func makeCompareFunc(matching map[string]record) func([]string) {
	return func(records []string) {
		if records[9] != "accepted" {
			return
		}
		matching[fmt.Sprintf("%s:%s", records[1], records[4])] = record{
			student: records[4],
		}
	}
}

type record struct {
	student string
}
