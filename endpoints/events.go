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
        
        stmtIns.Exec(event.Name, attributesString)
        
        fmt.Fprintf(w, "Yay!")
    default:
        w.WriteHeader(404)
    }
}
