package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func extractUUID(r *http.Request, prefix string) string {
	return strings.TrimPrefix(r.URL.Path, prefix)
}

func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/services/", serviceActionHandler)
	mux.HandleFunc("/monitoring", serviceMonitoringHandler)

	return mux
}

func serviceActionHandler(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/services/"), "/")

    if len(parts) == 1 && parts[0] == "create" && r.Method == http.MethodPost {
        var service Service
        if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
            jsonResponse(w, Result{Success: false, Message: "body invalide"})
            return
        }
        jsonResponse(w, CreateService(service))
        return
    }

    if len(parts) < 2 {
        http.Error(w, "route invalide", http.StatusBadRequest)
        return
    }

    uuid := parts[0]
    action := parts[1]

	switch r.Method {
	case http.MethodPost:
		switch action {
		case "start":
			jsonResponse(w, StartService(uuid))
		case "stop":
			jsonResponse(w, StopService(uuid))
		case "restart":
			jsonResponse(w, RestartService(uuid))
		default:
			http.Error(w, "action inconnue", http.StatusNotFound)
		}

	case http.MethodDelete:
		if action == "delete" {
			jsonResponse(w, DeleteService(uuid))
		} else {
			http.Error(w, "action inconnue", http.StatusNotFound)
		}

	default:
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

func serviceMonitoringHandler(w http.ResponseWriter, r *http.Request) {
    jsonResponse(w, GetMonitoring())
}