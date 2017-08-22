package endpoints

import (
    "encoding/json"
    "fmt"
    "net/http"
    
    "../models"
)

func EventsHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "POST":
        if key, ok := r.Header["Authorization"]; ok {
            if key[0] != Api_key {
                w.WriteHeader(403)
                return
            }
        } else {
            w.WriteHeader(403)
            return
        }
        
        decoder := json.NewDecoder(r.Body)
        var event models.Event
        decoder.Decode(&event)
        defer r.Body.Close()
        
        attributesString, _ := json.Marshal(event.Attributes)
        
        stmtIns, _ := Db.Prepare("INSERT INTO events (name, attributes) VALUES (?, ?)")
        defer stmtIns.Close()
        
        insertResult, err := stmtIns.Exec(event.Name, attributesString)
        if err != nil {
            fmt.Println(err)
        }
        
        id, _ := insertResult.LastInsertId()
        result := Db.QueryRow("SELECT * FROM events WHERE id = ? LIMIT 1", id)
        
        var resultAttributes string
        var resultEvent models.Event
        result.Scan(&resultEvent.Id, &resultEvent.Name, &resultAttributes, &resultEvent.Created_at)
        resultEvent.Attributes = event.Attributes
        
        resultJson, _ := json.Marshal(resultEvent)
        fmt.Fprintf(w, string(resultJson))
    default:
        w.WriteHeader(404)
    }
}
