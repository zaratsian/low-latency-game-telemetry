package main

import (
    "encoding/json"
    "net/http"
    "fmt"
)

func main() {

	server := http.NewServeMux()

    // Health endpoint
    server.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("%v", "success")
		w.Write([]byte("success"))
    })

    server.HandleFunc("/get_endpoint", handleGetData)
    server.HandleFunc("/post_endpoint", handlePostData)
    
    err := http.ListenAndServe(":8080", server)
    if err != nil {
        fmt.Printf("[ ERROR ] http.ListenAndServe: %v", err)
    }
}

func handleGetData(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Extract GET parameters
    queryValues := r.URL.Query()
    k1 := queryValues.Get("k1")

    data := map[string]string{
        "k1": k1,
    }

	go processData(data)

    jsonData, _ := json.Marshal(data)
    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonData)
}

func handlePostData(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var data map[string]string
    json.NewDecoder(r.Body).Decode(&data)
    fmt.Println(data)

    go processData(data)

    w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprint(data)))
}

func processData(data map[string]string) {
    // TODO - Analyze data
    fmt.Println("Data:", data)
}
