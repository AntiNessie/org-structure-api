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

type updateDepartmentRequest struct {
	Name     *string `json:"name,omitempty"`
	ParentID *uint   `json:"parent_id,omitempty"`
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

func (h *DepartmentHandler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid department ID", http.StatusBadRequest)
		return
	}

	var req updateDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dept, err := h.service.Update(uint(id), req.Name, req.ParentID)
	if err != nil {
		switch err {
		case service.ErrDepartmentNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		case service.ErrNameRequired:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dept)
}

func (h *DepartmentHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid department ID", http.StatusBadRequest)
		return
	}

	mode := r.URL.Query().Get("mode")
	if mode == "" {
		http.Error(w, "mode parameter is required (cascade or reassign)", http.StatusBadRequest)
		return
	}

	var reassignToID *uint
	if mode == "reassign" {
		reassignStr := r.URL.Query().Get("reassign_to_department_id")
		if reassignStr == "" {
			http.Error(w, "reassign_to_department_id is required for reassign mode", http.StatusBadRequest)
			return
		}
		reassignID, err := strconv.ParseUint(reassignStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid reassign_to_department_id", http.StatusBadRequest)
			return
		}
		reassignIDUint := uint(reassignID)
		reassignToID = &reassignIDUint
	}

	err = h.service.Delete(uint(id), mode, reassignToID)
	if err != nil {
		switch err {
		case service.ErrDepartmentNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
