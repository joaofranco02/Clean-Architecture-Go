package usecase

import "time"

// CreateUserInput contém os dados necessários para criar um usuário.
// DTOs protegem as entidades de serem expostas diretamente —
// o mundo externo nunca manipula entidades, apenas DTOs.
type CreateUserInput struct {
	Name  string
	Email string
	Age   int
}

// UpdateUserInput contém os dados para atualizar um usuário existente.
type UpdateUserInput struct {
	ID    string
	Name  string
	Email string
	Age   int
}

// UserOutput é o DTO de saída — a representação do usuário
// que o use case retorna para as camadas externas.
// Note que é diferente da entidade: podemos omitir campos internos
// ou formatar dados de maneira específica para a saída.
type UserOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
