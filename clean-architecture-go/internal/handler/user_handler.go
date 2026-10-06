// Package handler contém os handlers HTTP (Camada 3 — Interface Adapters).
//
// Handlers são adaptadores que:
// 1. Recebem requisições HTTP (dados do mundo externo)
// 2. Convertem para DTOs que os use cases entendem
// 3. Chamam o use case apropriado
// 4. Convertem a resposta do use case para HTTP (JSON)
//
// Esta camada NÃO contém lógica de negócio — apenas tradução de protocolo.
package handler

import (
	"clean-architecture-go/internal/usecase"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// UserHandler gerencia as rotas HTTP relacionadas a usuários.
type UserHandler struct {
	useCase *usecase.UserUseCase
}

// NewUserHandler cria um novo handler injetando o use case.
func NewUserHandler(uc *usecase.UserUseCase) *UserHandler {
	return &UserHandler{useCase: uc}
}

// apiResponse é a estrutura padrão de resposta da API.
type apiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// RegisterRoutes registra todas as rotas do handler no mux.
func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users", h.handleUsers)
	mux.HandleFunc("/users/", h.handleUserByID)
}

// handleUsers roteia entre GET (listar) e POST (criar).
func (h *UserHandler) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listUsers(w, r)
	case http.MethodPost:
		h.createUser(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "método não permitido")
	}
}

// handleUserByID roteia operações em um usuário específico.
func (h *UserHandler) handleUserByID(w http.ResponseWriter, r *http.Request) {
	// Extrai o ID da URL: /users/{id}
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	if id == "" {
		respondError(w, http.StatusBadRequest, "ID é obrigatório")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getUser(w, r, id)
	case http.MethodPut:
		h.updateUser(w, r, id)
	case http.MethodDelete:
		h.deleteUser(w, r, id)
	default:
		respondError(w, http.StatusMethodNotAllowed, "método não permitido")
	}
}

// createUser — POST /users
func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateUserInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	output, err := h.useCase.CreateUser(input)
	if err != nil {
		// Erros de validação do domínio retornam 422 (Unprocessable Entity)
		respondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, apiResponse{
		Success: true,
		Data:    output,
	})
}

// getUser — GET /users/{id}
func (h *UserHandler) getUser(w http.ResponseWriter, _ *http.Request, id string) {
	output, err := h.useCase.GetUserByID(id)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    output,
	})
}

// listUsers — GET /users
func (h *UserHandler) listUsers(w http.ResponseWriter, _ *http.Request) {
	output, err := h.useCase.ListUsers()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    output,
	})
}

// updateUser — PUT /users/{id}
func (h *UserHandler) updateUser(w http.ResponseWriter, r *http.Request, id string) {
	var input usecase.UpdateUserInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	input.ID = id

	output, err := h.useCase.UpdateUser(input)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    output,
	})
}

// deleteUser — DELETE /users/{id}
func (h *UserHandler) deleteUser(w http.ResponseWriter, _ *http.Request, id string) {
	if err := h.useCase.DeleteUser(id); err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    "usuário removido com sucesso",
	})
}

// respondJSON envia uma resposta JSON com status code.
func respondJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

// respondError envia uma resposta de erro padronizada.
func respondError(w http.ResponseWriter, statusCode int, message string) {
	respondJSON(w, statusCode, apiResponse{
		Success: false,
		Error:   message,
	})
}
