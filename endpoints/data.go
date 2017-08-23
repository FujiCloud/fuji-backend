package endpoints

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "time"
    
    "github.com/montanaflynn/stats"
    
    "../models"
)

func DataHandler(w http.ResponseWriter, r *http.Request) {
    switch query := r.URL.Query().Get("q"); query {
    case "accessibility":
        counts := Db.QueryRow("SELECT SUM(voiceover = 1), SUM(bold_text = 1), SUM(reduce_motion = 1), SUM(reduce_transparency = 1), COUNT(*) FROM USERS")
        
        var voiceover int
        var bold_text int
        var reduce_motion int
        var reduce_transparency int
        var total int
        counts.Scan(&voiceover, &bold_text, &reduce_motion, &reduce_transparency, &total)
        
        var buffer bytes.Buffer
        buffer.WriteString("setting,count\n")
        buffer.WriteString("voiceover,")
        buffer.WriteString(strconv.Itoa(voiceover))
        buffer.WriteString("\nbold_text,")
        buffer.WriteString(strconv.Itoa(bold_text))
        buffer.WriteString("\nreduce_motion,")
        buffer.WriteString(strconv.Itoa(reduce_motion))
        buffer.WriteString("\nreduce_transparency,")
        buffer.WriteString(strconv.Itoa(reduce_transparency))
        buffer.WriteString("\ntotal,")
        buffer.WriteString(strconv.Itoa(total))
        
        fmt.Fprintf(w, buffer.String())
    case "content_view":
        rows, _ := Db.Query("SELECT * FROM events WHERE name = 'Content View' ORDER BY created_at ASC")
        defer rows.Close()
        
        dates := make([]string, 0)
        events := make(map[string]map[string]int)
        event_names := make([]string, 0)
        
        for rows.Next() {
            var attributes string
            var event models.Event
            rows.Scan(&event.Id, &event.Name, &attributes, &event.Created_at)
            
            json.Unmarshal([]byte(attributes), &event.Attributes)
            page := event.Attributes["page"]
            
            time_string := event.Created_at.Format("2006 Jan 02")
            
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
    case "duration":
        rows, _ := Db.Query("SELECT duration, created_at FROM sessions WHERE duration IS NOT NULL ORDER BY created_at ASC")
        defer rows.Close()
        
        var buffer bytes.Buffer
        buffer.WriteString("date,pct05,pct25,pct50,pct75,pct95\n")
        
        save := func(date string, values []float64) {
            buffer.WriteString(date)
            buffer.WriteString(",")
            
            fifth_percentile, err := stats.Percentile(values, 5)
            if err == nil {
                buffer.WriteString(strconv.FormatFloat(fifth_percentile, 'f', 2, 64))
            } else {
                buffer.WriteString("0")
            }
            buffer.WriteString(",")
            
            twentyfifth_percentile, err := stats.Percentile(values, 25)
            if err == nil {
                buffer.WriteString(strconv.FormatFloat(twentyfifth_percentile, 'f', 2, 64))
            } else {
                buffer.WriteString("0")
            }
            buffer.WriteString(",")
            
            median, err := stats.Percentile(values, 50)
            if err == nil {
                buffer.WriteString(strconv.FormatFloat(median, 'f', 2, 64))
            } else {
                buffer.WriteString("0")
            }
            buffer.WriteString(",")
            
            seventyfifth_percentile, err := stats.Percentile(values, 75)
            if err == nil {
                buffer.WriteString(strconv.FormatFloat(seventyfifth_percentile, 'f', 2, 64))
            } else {
                buffer.WriteString("0")
            }
            buffer.WriteString(",")
            
            ninetyfifth_percentile, err := stats.Percentile(values, 95)
            if err == nil {
                buffer.WriteString(strconv.FormatFloat(ninetyfifth_percentile, 'f', 2, 64))
            } else {
                buffer.WriteString("0")
            }
            buffer.WriteString("\n")
        }
        
        current_date := ""
        values := make([]float64, 0)
        
        for rows.Next() {
            var duration int
            var date time.Time
            rows.Scan(&duration, &date)
            
            date_string := date.Format("2006 Jan 02")
            
            if current_date == "" {
                current_date = date_string
            }
            
            if date_string != current_date {
                save(current_date, values)
                
                values = make([]float64, 0)
                current_date = date_string
            }
            
            values = append(values, float64(duration))
        }
        
        save(current_date, values)
        
        fmt.Fprintf(w, buffer.String())
    case "maus":
        rows, _ := Db.Query("SELECT DATE(created_at) date, COUNT(DISTINCT user_id) count FROM sessions GROUP BY DATE(created_at) ORDER BY date ASC")
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
