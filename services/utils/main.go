package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func HandleJSON(data []byte, schema string) {
    var jsonData map[string]interface{}
    json.Unmarshal(data, &jsonData)

    // Check for schema integrity
    err := json.Unmarshal([]byte(schema), &jsonData)
    if err != nil {
        log.Println("Error: Invalid JSON schema")
        return
    }

    // Convert JSON data to human-readable sentence
    sentence := "The name is " + jsonData["name"].(string) + ", age is " +
        jsonData["age"].(string) + " and the address is " + jsonData["address"].(string)
    fmt.Println(sentence)
}

func FormatJSON(records []json.RawMessage) string {
	// Decode the raw JSON data into a map of key-value pairs
	var data []map[string]interface{}
	for _, record := range records {
		var recordMap map[string]interface{}
		err := json.Unmarshal(record, &recordMap)
		if err != nil {
			return err.Error()
		}
		data = append(data, recordMap)
	}

	// Format data as a list kv pairs
	formattedData := make([]string, len(data))
	for i, recordMap := range data {
		formattedData[i] = fmt.Sprintf("%s: %v", "key", recordMap["key"])
	}

	formattedString := strings.Join(formattedData, ", ")
	return formattedString
}

func FormatStruct(s interface{}) (string, string) {
	// Use reflection to get the fields of the struct
	st := reflect.TypeOf(s)
	sv := reflect.ValueOf(s)

	structNames := []string{}
	structValues := []string{}

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)

		fieldValue := sv.FieldByName(field.Name)
		fieldValueType := fieldValue.Type().String()
		structFieldValue := fieldValue.Interface()

		// Convert the interface to string
		stringValue, ok := structFieldValue.(string)
		if !ok {
			fmt.Printf("Field Type: %v\n", fieldValueType)
			if fieldValueType == "int" {
				stringValue = strconv.Itoa(structFieldValue.(int))
			} else if fieldValueType == "float64" {
				stringValue = fmt.Sprintf("%f", structFieldValue)
			} else if fieldValueType == "bool" {
				stringValue = strconv.FormatBool(structFieldValue.(bool))
			}
		} else {
			stringValue = "\"" + stringValue + "\""
		}

		// Append items to list
		structNames = append(structNames, field.Name)
		structValues = append(structValues, stringValue)

	}

	keyString := strings.Join(structNames, ", ")
	valueString := strings.Join(structValues, ", ")
	return keyString, valueString
}

func FormatInterface(data interface{}) (string, string) {
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Map {
		return "",""
	}

	structNames := []string{}
	structValues := []string{}

	keys := val.MapKeys()

	for _, key := range keys {
		// Process Key
		structNames = append(structNames, key.String())
		
		// Process Value               val.MapIndex(key).Type()
		fmt.Printf("Value: %v, %T\n", val.MapIndex(key), val.MapIndex(key))
		/*
		if val.MapIndex(key).Type() == "int" {
			stringValue = strconv.Itoa(structFieldValue.(int))
		} else if val.MapIndex(key).Type() == "bool" {
			stringValue = strconv.FormatBool(structFieldValue.(bool))
		} else {
			stringValue = "\"" + stringValue + "\""
		}
		*/
		//structValues = append(structValues, fmt.Sprint(val.MapIndex(key)))
	}

	keyString := strings.Join(structNames, ", ")
	valueString := strings.Join(structValues, ", ")
	return keyString, valueString
}
