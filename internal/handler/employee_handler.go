package handler

import (
	"encoding/json"
	"net/http"
	"org-structure-api/internal/service"
	"strconv"

	"github.com/gorilla/mux"
)

type EmployeeHandler struct {
	service service.EmployeeService
}

func NewEmployeeHandler(service service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

type createEmployeeRequest struct {
	FullName string `json:"full_name"`
	Position string `json:"position"`
}

func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	departmentID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid department ID", http.StatusBadRequest)
		return
	}

	var req createEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	emp, err := h.service.Create(uint(departmentID), req.FullName, req.Position, nil)
	if err != nil {
		switch err {
		case service.ErrDepartmentNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		case service.ErrEmployeeNameRequired, service.ErrPositionRequired:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(emp)
}
