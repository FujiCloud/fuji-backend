package endpoints

import (
    "encoding/json"
    "fmt"
    "net/http"
    
    "../models"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {
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
        var user models.User
        decoder.Decode(&user)
        defer r.Body.Close()
        
        stmtIns, _ := Db.Prepare("INSERT INTO users (os, device, locale, voiceover, bold_text, reduce_motion, reduce_transparency) VALUES (?, ?, ?, ?, ?, ?, ?)")
        defer stmtIns.Close()
        
        stmtIns.Exec(user.Os, user.Device, user.Locale, user.Voiceover, user.Bold_text, user.Reduce_motion, user.Reduce_transparency)
        
        fmt.Fprintf(w, "Yay!")
    default:
        w.WriteHeader(404)
    }
}
