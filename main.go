package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
    
    _ "github.com/go-sql-driver/mysql"
    "github.com/nu7hatch/gouuid"
    
    "./endpoints"
    "./models"
)

var Api_key string
var Db *sql.DB

func main() {
    raw, err := ioutil.ReadFile("config.json")
    
    if err == nil {
        var config models.ConfigFile
        json.Unmarshal(raw, &config)
        Api_key = config.Apikey
    } else {
        uuid, _ := uuid.NewV4()
        Api_key = strings.Replace(uuid.String(), "-", "", 4)
        
        config := models.ConfigFile{Api_key}
        configJson, _ := json.Marshal(config)
        ioutil.WriteFile("config.json", configJson, 0644)
    }
    
    fmt.Printf("API Key: %s\n", Api_key)
    
    Db, _ = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/fuji?parseTime=true")
    
    http.Handle("/", http.FileServer(http.Dir("./frontend")))
    
    endpoints.Api_key = Api_key
    endpoints.Db = Db
    
    http.HandleFunc("/dashboard", endpoints.DashboardHandler)
    http.HandleFunc("/data", endpoints.DataHandler)
    http.HandleFunc("/events", endpoints.EventsHandler)
    
    fmt.Printf("Running on localhost:8000...\n")
    http.ListenAndServe(":8000", nil)
}
