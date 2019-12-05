package main

import (
	"./jsoncsv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func readBytes(fromFile string) ([]byte, error) {
	f, err := os.Open(fromFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	byteValue, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return byteValue, nil
}

func writeBytes(toFile string, data []byte) error {
	f, err := os.Create(toFile)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func testWrite() {
	jsonBytes, err := readBytes("example.json")
	if err != nil {
		log.Fatal(err)
	}

	csvData, err := jsoncsv.JSON2CSV(jsonBytes)
	if err != nil {
		log.Fatal(err)
	}

	writeBytes("csvfile.csv", csvData)

	fmt.Println("json to csv written!")
}

func testRead() {
	csvBytes, err := readBytes("csvfile.csv")
	if err != nil {
		log.Fatal(err)
	}

	jsonData, err := jsoncsv.CSV2JSON(csvBytes)
	if err != nil {
		log.Fatal(err)
	}

	writeBytes("decoded.json", jsonData)

	fmt.Println("csv to json read!")
}

func main() {
	testWrite()
	testRead()
}
