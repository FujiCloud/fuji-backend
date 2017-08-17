package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "net/http"
    
    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Event struct {
    Name string
    Attributes map[string]string
}

func handler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "POST":
        if api_key, ok := r.Header["Authorization"]; ok {
            if api_key[0] != "hello" {
                w.WriteHeader(403)
                return
            }
        } else {
            w.WriteHeader(403)
            return
        }
        
        decoder := json.NewDecoder(r.Body)
        var event Event
        decoder.Decode(&event)
        defer r.Body.Close()
        
        attributesString, _ := json.Marshal(event.Attributes)
        
        stmtIns, _ := db.Prepare("INSERT INTO events (name, attributes) VALUES (?, ?)")
        defer stmtIns.Close()
        
        stmtIns.Exec(event.Name, attributesString)
        
        fmt.Fprintf(w, "Yay!")
    default:
        w.WriteHeader(404)
    }
}

func main() {
    db, _ = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/fuji")
    
    http.HandleFunc("/events", handler)
    
    fmt.Printf("Running on localhost:8000...")
    http.ListenAndServe(":8000", nil)
}
