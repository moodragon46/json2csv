package jsoncsv

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"strconv"
	"strings"
)

// Returns if its an array, and the key.
func checkArr(src string) (bool, string) {
	actualIndex := src[4:(len(src) - 1)]
	if src[0] == 'm' {
		return false, actualIndex
	}
	return true, actualIndex
}

func decodeVal(val string) interface{} {
	if val[0] == 'f' {
		val, err := strconv.ParseFloat(val[6:(len(val)-1)], 64)
		if err != nil {
			log.Fatal(err)
		}
		return val
	}

	return val[7:(len(val) - 1)]
}

// Smart index figures out from the infokey if it is a array or map, and uses the appropriate indexing.
func smartIndexSet(obj interface{}, infokey string, value interface{}) interface{} {
	isArr, key := checkArr(infokey)
	if isArr {
		val, err := strconv.Atoi(key)
		if err != nil {
			log.Fatal("Trying to sample from a list with a string")
		}

		if obj == nil {
			obj = make([]interface{}, 0)
		}
		arrObj := obj.([]interface{})
		for val >= len(arrObj) {
			arrObj = append(arrObj, nil)
		}
		arrObj[val] = value

		return arrObj
	}

	if obj == nil {
		obj = make(map[string]interface{})
	}
	mapObj := obj.(map[string]interface{})
	mapObj[key] = value

	return mapObj
}

func smartIndexGet(obj interface{}, infokey string) interface{} {
	isArr, key := checkArr(infokey)
	if isArr {
		val, err := strconv.Atoi(key)
		if err != nil {
			log.Fatal("Trying to sample from a list with a string")
		}

		if obj == nil {
			obj = make([]interface{}, 0)
		}
		arrObj := obj.([]interface{})
		for val >= len(arrObj) {
			arrObj = append(arrObj, nil)
		}
		return arrObj[val]
	}

	if obj == nil {
		obj = make(map[string]interface{})
	}
	mapObj := obj.(map[string]interface{})
	return mapObj[key]
}

func setKeyToVal(tree interface{}, infokeys []string, val interface{}) interface{} {
	if len(infokeys) > 1 {
		// Here we want to be calling setKeyToVal with the subtree that the smartIndexSet creates
		return smartIndexSet(tree, infokeys[0], setKeyToVal(smartIndexGet(tree, infokeys[0]), infokeys[1:], val))
	}
	return smartIndexSet(tree, infokeys[0], val)
}

func setNodeAtPath(path string, val interface{}, tree *interface{}) {
	keys := strings.Split(path, "/")[1:]

	*tree = setKeyToVal(*tree, keys, val)
}

//CSV2JSON Takes in a csv bytes slice and outputs a json bytes slice. Will only work with csv encoded using the Json2Csv function.
func CSV2JSON(csvData []byte) ([]byte, error) {
	bytesReader := bytes.NewReader(csvData)

	reader := csv.NewReader(bytesReader)

	data := make([][]string, 0)

	for {
		line, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		data = append(data, line)
	}

	var reconstructed interface{}
	for i := range data[0] {
		key := data[0][i]
		val := decodeVal(data[1][i])

		setNodeAtPath(key, val, &reconstructed)
	}

	jsonData, err := json.Marshal(reconstructed)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
