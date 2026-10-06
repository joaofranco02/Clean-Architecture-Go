// Package entity contém as entidades de domínio (Camada 1 — Domain).
// Esta é a camada mais interna da Clean Architecture.
// Entidades encapsulam regras de negócio fundamentais que existem
// independentemente de qualquer aplicação, framework ou tecnologia.
package entity

import (
	"errors"
	"strings"
	"time"
)

// Erros de domínio — representam violações das regras de negócio.
var (
	ErrInvalidName  = errors.New("nome é obrigatório e deve ter pelo menos 3 caracteres")
	ErrInvalidEmail = errors.New("email é obrigatório e deve ser um endereço válido")
	ErrInvalidAge   = errors.New("idade deve ser entre 1 e 150")
)

// User é a entidade central do domínio.
// Ela contém apenas dados e regras de negócio fundamentais.
// NÃO conhece banco de dados, HTTP, JSON ou qualquer detalhe externo.
type User struct {
	ID        string
	Name      string
	Email     string
	Age       int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser cria uma nova entidade User validada.
// Esta é a única forma correta de criar um User — garante que
// as regras de negócio sejam respeitadas desde a criação.
func NewUser(name, email string, age int) (*User, error) {
	user := &User{
		Name:      strings.TrimSpace(name),
		Email:     strings.TrimSpace(strings.ToLower(email)),
		Age:       age,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

// Validate verifica se todas as regras de negócio da entidade são válidas.
// Pode ser chamado a qualquer momento para revalidar o estado.
func (u *User) Validate() error {
	if len(u.Name) < 3 {
		return ErrInvalidName
	}

	if !isValidEmail(u.Email) {
		return ErrInvalidEmail
	}

	if u.Age < 1 || u.Age > 150 {
		return ErrInvalidAge
	}

	return nil
}

// Update aplica alterações à entidade, revalidando as regras de negócio.
func (u *User) Update(name, email string, age int) error {
	u.Name = strings.TrimSpace(name)
	u.Email = strings.TrimSpace(strings.ToLower(email))
	u.Age = age
	u.UpdatedAt = time.Now()

	return u.Validate()
}

// isValidEmail é uma regra de negócio: validação básica de formato de email.
func isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	// Deve ter algo antes e depois do @, e o domínio deve ter um ponto
	return len(parts[0]) > 0 && strings.Contains(parts[1], ".")
}
