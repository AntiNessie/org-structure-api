package handler

import (
	"encoding/json"
	"net/http"
	"org-structure-api/internal/service"
	"strconv"

	"github.com/gorilla/mux"
)

type DepartmentHandler struct {
	service service.DepartmentService
}

func NewDepartmentHandler(service service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: service}
}

type createDepartmentRequest struct {
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id,omitempty"`
}

func (h *DepartmentHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var req createDepartmentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dept, err := h.service.Create(req.Name, req.ParentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dept)
}

func (h *DepartmentHandler) GetDepartment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid department ID", http.StatusBadRequest)
		return
	}

	includeEmployees := r.URL.Query().Get("include_employees") != "false"

	dept, err := h.service.GetByID(uint(id), includeEmployees)
	if err != nil {
		switch err {
		case service.ErrDepartmentNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dept)
}
