package main

import (
    "bytes"
    "database/sql"
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    
    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Event struct {
    Name string
    Attributes map[string]string
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "dashboard.html")
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
    rows, _ := db.Query("SELECT * FROM events WHERE name = 'Content View'")
    defer rows.Close()
    
    pages := make(map[string]int)
    
    for rows.Next() {
        var id int
        var name string
        var attributes string
        rows.Scan(&id, &name, &attributes)
        
        eventJson := fmt.Sprintf("{\"name\": \"%s\", \"attributes\": %s}", name, attributes)
        var event Event
        json.Unmarshal([]byte(eventJson), &event)
        page := event.Attributes["page"]
        
        if _, ok := pages[page]; ok {
            pages[page] += 1
        } else {
            pages[page] = 1
        }
    }
    
    var buffer bytes.Buffer
    buffer.WriteString("page,count\n")
    
    for k, v := range pages {
        buffer.WriteString(k)
        buffer.WriteString(",")
        buffer.WriteString(strconv.Itoa(v))
        buffer.WriteString("\n")
    }
    
    fmt.Fprintf(w, buffer.String())
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
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
    
    http.HandleFunc("/dashboard", dashboardHandler)
    http.HandleFunc("/data", dataHandler)
    http.HandleFunc("/events", eventsHandler)
    
    fmt.Printf("Running on localhost:8000...\n")
    http.ListenAndServe(":8000", nil)
}
