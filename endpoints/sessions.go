package endpoints

import (
    "encoding/json"
    "fmt"
    "net/http"
    
    "../models"
)

func SessionsHandler(w http.ResponseWriter, r *http.Request) {
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
        var session models.Session
        decoder.Decode(&session)
        defer r.Body.Close()
        
        stmtIns, _ := Db.Prepare("INSERT INTO SESSIONS (user_id) VALUES (?)")
        defer stmtIns.Close()
        
        stmtIns.Exec(session.User_id)
        
        fmt.Fprintf(w, "Yay!")
    default:
        w.WriteHeader(404)
    }
}
