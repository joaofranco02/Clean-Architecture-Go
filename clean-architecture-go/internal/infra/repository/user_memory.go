// Package repository contém as implementações concretas dos repositórios
// (Camada 4 — Frameworks & Drivers / Infraestrutura).
//
// Esta camada implementa as interfaces definidas na camada de use cases.
// É aqui que vivem os detalhes de implementação: banco de dados, APIs externas, etc.
//
// PRINCÍPIO: Esta é a camada que mais muda. Trocar de banco de dados significa
// apenas criar uma nova implementação da interface UserRepository.
package repository

import (
	"clean-architecture-go/internal/entity"
	"fmt"
	"sync"
)

// UserMemoryRepository é uma implementação em memória do UserRepository.
// Em produção, seria substituída por PostgresUserRepository, MongoUserRepository, etc.
// O use case não precisa saber qual implementação está sendo usada.
type UserMemoryRepository struct {
	mu    sync.RWMutex
	users map[string]*entity.User
	nextID int
}

// NewUserMemoryRepository cria um novo repositório em memória.
func NewUserMemoryRepository() *UserMemoryRepository {
	return &UserMemoryRepository{
		users:  make(map[string]*entity.User),
		nextID: 1,
	}
}

// Save persiste um novo usuário em memória e atribui um ID auto-incrementado.
func (r *UserMemoryRepository) Save(user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = fmt.Sprintf("%d", r.nextID)
	r.nextID++

	// Armazena uma cópia para evitar mutação externa
	stored := *user
	r.users[user.ID] = &stored

	return nil
}

// FindByID busca um usuário pelo ID. Retorna nil, nil se não encontrado.
func (r *UserMemoryRepository) FindByID(id string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, nil
	}

	// Retorna uma cópia para evitar mutação externa
	copy := *user
	return &copy, nil
}

// FindAll retorna todos os usuários armazenados.
func (r *UserMemoryRepository) FindAll() ([]*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*entity.User, 0, len(r.users))
	for _, user := range r.users {
		copy := *user
		users = append(users, &copy)
	}

	return users, nil
}

// Update atualiza os dados de um usuário existente.
func (r *UserMemoryRepository) Update(user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return fmt.Errorf("usuário com ID %s não encontrado", user.ID)
	}

	stored := *user
	r.users[user.ID] = &stored

	return nil
}

// Delete remove um usuário pelo ID.
func (r *UserMemoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return fmt.Errorf("usuário com ID %s não encontrado", id)
	}

	delete(r.users, id)
	return nil
}
