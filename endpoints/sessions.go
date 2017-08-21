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
        
        stmtIns, _ := Db.Prepare("INSERT INTO sessions (user_id) VALUES (?)")
        defer stmtIns.Close()
        
        insertResult, _ := stmtIns.Exec(session.User_id)
        
        id, _ := insertResult.LastInsertId()
        result := Db.QueryRow("SELECT * FROM sessions WHERE id = ? LIMIT 1", id)
        
        var resultSession models.Session
        result.Scan(&resultSession.Id, &resultSession.User_id, &resultSession.Duration, &resultSession.Created_at)
        
        resultJson, _ := json.Marshal(resultSession)
        fmt.Fprintf(w, string(resultJson))
    case "PATCH":
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
        
        stmtUpdate, _ := Db.Prepare("UPDATE sessions SET duration = ? WHERE id = ?")
        defer stmtUpdate.Close()
        
        stmtUpdate.Exec(session.Duration, session.Id)
        result := Db.QueryRow("SELECT * FROM sessions WHERE id = ? LIMIT 1", session.Id)
        
        var resultSession models.Session
        result.Scan(&resultSession.Id, &resultSession.User_id, &resultSession.Duration, &resultSession.Created_at)
        
        resultJson, _ := json.Marshal(resultSession)
        fmt.Fprintf(w, string(resultJson))
    default:
        w.WriteHeader(404)
    }
}
