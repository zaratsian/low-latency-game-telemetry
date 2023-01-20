package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"context"
	"reflect"
	"strconv"
	"strings"

	"cloud.google.com/go/spanner"
)

// Define the JSON schema
const schema = `{
    "type": "object",
    "properties": {
        "name": { "type": "string" },
        "age": { "type": "integer" },
        "address": { "type": "string" }
    },
    "required": ["name", "age", "address"]
}`

func main() {

	ctx := context.Background()

	udpAddr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		panic(err)
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		panic(err)
	}
	defer udpConn.Close()

	for {
		// Process incoming data
		data, err := preProcessData(udpConn)
		if err != nil {
			log.Println(err)
		}

		fmt.Printf("\nRawdata: %v, %T\n", data, data)

		// Write data to database (default = Spanner)
		/*
		err = spannerWriteDML(ctx, data)
		if err != nil {
			log.Print(err)
		} else {
			log.Print("Successfully wrote data to Spanner")
		}
		*/

		// Send data to other sidecar container
		/*
		   if payload.Data == "send" {
		       otherAddr, _ := net.ResolveUDPAddr("udp", "other-container:8080")
		       _, err := udpConn.WriteToUDP([]byte(payload.Data), otherAddr)
		       if err != nil {
		           fmt.Println(err)
		       }
		   }
		*/
	}
}

func handleJSON(data []byte) {
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

func preProcessData(conn *net.UDPConn) (interface{}, error) {
	// Buffer packets
	buffer := make([]byte, 2048)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Process UDP data as Interface{}
	var data map[string]interface{}
	err = json.Unmarshal(buffer[:n], &data)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Process UDP data as string
	//data := string(buffer[:n])

	//unmarshal the struct to json
	//result, _ := json.Marshal(data)

	return data, nil
}

/*
func dataAnalysis(data interface{}) (interface{}, err) {

}
*/

func spannerWriteDML(ctx context.Context, data interface{}) error {

	gcpProjectId    := os.Getenv("TF_VAR_GCP_PROJECT_ID")
	spannerInstance := os.Getenv("TF_VAR_SPANNER_INSTANCE")
	spannerDatabase := os.Getenv("TF_VAR_SPANNER_DATABASE")
	spannerTable    := os.Getenv("TF_VAR_SPANNER_TABLE_GAME_TELEMETRY")

	connectionStr := fmt.Sprintf("projects/%v/instances/%v/databases/%v", gcpProjectId, spannerInstance, spannerDatabase)

	spannerClient, err := spanner.NewClient(ctx, connectionStr)
	if err != nil {
		return err
	}
	defer spannerClient.Close()

	keyString, valueString := formatInterface(data)

	// Generate DML
	dml := fmt.Sprintf("INSERT %v (%v) VALUES (%v)", spannerTable, keyString, valueString)
	fmt.Printf("dml: %v\n", dml)

	_, err = spannerClient.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: dml,
		}
		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		log.Printf("%d record(s) inserted.\n", rowCount)
		return err
	})
	return err

}

func formatJSON(records []json.RawMessage) string {
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

func formatStruct(s interface{}) (string, string) {
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
			if fieldValueType == "int" {
				stringValue = strconv.Itoa(structFieldValue.(int))
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

func formatInterface(data interface{}) (string, string) {
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