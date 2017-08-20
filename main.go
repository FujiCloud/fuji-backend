package main

import (
    "bytes"
    "database/sql"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "strconv"
    "strings"
    "time"
    
    _ "github.com/go-sql-driver/mysql"
    "github.com/nu7hatch/gouuid"
)

var api_key string
var db *sql.DB

type ConfigFile struct {
    Apikey string `json:"api_key"`
}

type Event struct {
    Name string
    Attributes map[string]string
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "dashboard.html")
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
    rows, _ := db.Query("SELECT * FROM events WHERE name = 'Content View' ORDER BY created_at ASC")
    defer rows.Close()
    
    dates := make([]string, 0)
    events := make(map[string]map[string]int)
    event_names := make([]string, 0)
    
    for rows.Next() {
        var id int
        var name string
        var attributes string
        var created_at time.Time
        rows.Scan(&id, &name, &attributes, &created_at)
        
        eventJson := fmt.Sprintf("{\"name\": \"%s\", \"attributes\": %s}", name, attributes)
        var event Event
        json.Unmarshal([]byte(eventJson), &event)
        page := event.Attributes["page"]
        
        time_string := created_at.Format("2006 Jan 02")
        
        if _, ok := events[time_string]; ok {
            if _, pageOk := events[time_string][page]; pageOk {
                events[time_string][page] += 1
            } else {
                events[time_string][page] = 1
            }
        } else {
            events[time_string] = make(map[string]int)
            events[time_string][page] = 1
        }
        
        if _, ok := events[time_string]["_total"]; ok {
            events[time_string]["_total"] += 1
        } else {
            events[time_string]["_total"] = 1
        }
        
        names_contains := false
        
        for _, value := range event_names {
            if value == page {
                names_contains = true
                break
            }
        }
        
        if !names_contains {
            event_names = append(event_names, page)
        }
        
        dates_contains := false
        
        for _, value := range dates {
            if value == time_string {
                dates_contains = true
                break
            }
        }
        
        if !dates_contains {
            dates = append(dates, time_string)
        }
    }
    
    var buffer bytes.Buffer
    buffer.WriteString("date")
    
    for _, value := range event_names {
        buffer.WriteString(",")
        buffer.WriteString(value)
    }
    
    buffer.WriteString("\n")
    
    for _, date := range dates {
        buffer.WriteString(date)
        
        for i := range event_names {
            page := event_names[i]
            
            if page != "_total" {
                if _, ok := events[date][page]; ok {
                    buffer.WriteString(",")
                    buffer.WriteString(strconv.FormatFloat(100 * float64(events[date][page]) / float64(events[date]["_total"]), 'f', 2, 64))
                } else {
                    buffer.WriteString(",0")
                }
            }
        }
        
        buffer.WriteString("\n")
    }
    
    fmt.Fprintf(w, buffer.String())
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "POST":
        if key, ok := r.Header["Authorization"]; ok {
            if key[0] != api_key {
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
    raw, err := ioutil.ReadFile("config.json")
    
    if err == nil {
        var config ConfigFile
        json.Unmarshal(raw, &config)
        api_key = config.Apikey
    } else {
        uuid, _ := uuid.NewV4()
        api_key = strings.Replace(uuid.String(), "-", "", 4)
        
        config := ConfigFile{api_key}
        configJson, _ := json.Marshal(config)
        ioutil.WriteFile("config.json", configJson, 0644)
    }
    
    fmt.Printf("API Key: %s\n", api_key)
    
    db, _ = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/fuji?parseTime=true")
    
    http.Handle("/", http.FileServer(http.Dir("./frontend")))
    http.HandleFunc("/dashboard", dashboardHandler)
    http.HandleFunc("/data", dataHandler)
    http.HandleFunc("/events", eventsHandler)
    
    fmt.Printf("Running on localhost:8000...\n")
    http.ListenAndServe(":8000", nil)
}
