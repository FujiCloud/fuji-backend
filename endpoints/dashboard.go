package endpoints

import (
    "net/http"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "dashboard.html")
}
