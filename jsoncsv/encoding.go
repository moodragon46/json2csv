package jsoncsv

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
)

func interfaceToString(node interface{}) string {
	floatVal, isFloat := node.(float64)
	if isFloat {
		return "float<" + fmt.Sprint(floatVal) + ">"
	}

	return "string<" + node.(string) + ">"
}

func concatMap(a *map[string]string, b map[string]string) {
	for key, val := range b {
		(*a)[key] = val
	}
}

func traverseTree(currPath string, node interface{}) map[string]string {
	newKeyvalPairs := make(map[string]string)

	list, isList := node.([]interface{})

	if isList {
		for i, val := range list {
			concatMap(&newKeyvalPairs, traverseTree(currPath+"/arr<"+strconv.Itoa(i)+">", val))
		}
		return newKeyvalPairs
	}

	mapVal, isMap := node.(map[string]interface{})

	if isMap {
		for key, val := range mapVal {
			concatMap(&newKeyvalPairs, traverseTree(currPath+"/map<"+key+">", val))
		}
		return newKeyvalPairs
	}

	newKeyvalPairs[currPath] = interfaceToString(node)
	return newKeyvalPairs
}

func transpose(slice [][]string) [][]string {
	xl := len(slice[0])
	yl := len(slice)
	result := make([][]string, xl)
	for i := range result {
		result[i] = make([]string, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = slice[j][i]
		}
	}
	return result
}

func mapTo2DArr(keyvalPairs map[string]string) [][]string {
	arr := make([][]string, 0)

	for key, val := range keyvalPairs {
		arr = append(arr, []string{key, val})
	}

	return transpose(arr)
}

//JSON2CSV Reads a json byte slice and returns a csv byte slice.
func JSON2CSV(jsonData []byte) ([]byte, error) {
	var result interface{}
	json.Unmarshal(jsonData, &result)

	keyvalPairs := traverseTree("", result)
	data := mapTo2DArr(keyvalPairs)

	buf := new(bytes.Buffer)

	writer := csv.NewWriter(buf)

	for _, val := range data {
		err := writer.Write(val)
		if err != nil {
			return nil, err
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
