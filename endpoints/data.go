package endpoints

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "time"
    
    "../models"
)

func DataHandler(w http.ResponseWriter, r *http.Request) {
    switch query := r.URL.Query().Get("q"); query {
    case "content_view":
        rows, _ := Db.Query("SELECT * FROM events WHERE name = 'Content View' ORDER BY created_at ASC")
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
            var event models.Event
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
    case "maus":
        rows, _ := Db.Query("SELECT DATE(created_at) date, COUNT(DISTINCT id) count FROM events GROUP BY DATE(created_at) ORDER BY date ASC")
        defer rows.Close()
        
        dates := make([]string, 0)
        counts := make(map[string]int)
        
        for rows.Next() {
            var date time.Time
            var count int
            rows.Scan(&date, &count)
            
            date_string := date.Format("2006 Jan 02")
            counts[date_string] = count
            
            dates_contains := false
            
            for _, value := range dates {
                if value == date_string {
                    dates_contains = true
                    break
                }
            }
            
            if !dates_contains {
                dates = append(dates, date_string)
            }
        }
        
        var buffer bytes.Buffer
        buffer.WriteString("date,count\n")
        
        for _, date := range dates {
            buffer.WriteString(date)
            buffer.WriteString(",")
            buffer.WriteString(strconv.Itoa(counts[date]))
            buffer.WriteString("\n")
        }
        
        fmt.Fprintf(w, buffer.String())
    default:
        fmt.Fprintf(w, "lol")
    }
}
