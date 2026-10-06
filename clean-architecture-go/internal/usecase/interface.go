// Package usecase contém os casos de uso da aplicação (Camada 2 — Application Business Rules).
// Esta camada orquestra o fluxo de dados entre as entidades e o mundo externo.
//
// PRINCÍPIO CHAVE: Esta camada define INTERFACES (ports/portas) que as camadas
// externas devem implementar. Isso é a Inversão de Dependência (DIP) —
// o use case dita o contrato, não a infraestrutura.
package usecase

import "clean-architecture-go/internal/entity"

// UserRepository define o contrato para persistência de usuários.
// Esta interface é definida AQUI (camada 2), mas IMPLEMENTADA
// na camada de infraestrutura (camada 4).
//
// Isso garante que:
// - O use case não depende de nenhum banco específico
// - Podemos trocar a implementação (memória, PostgreSQL, MongoDB) sem alterar o use case
// - Testes podem usar mocks facilmente
type UserRepository interface {
	// Save persiste um novo usuário e atribui um ID.
	Save(user *entity.User) error

	// FindByID busca um usuário pelo seu identificador único.
	// Retorna nil, nil se o usuário não for encontrado.
	FindByID(id string) (*entity.User, error)

	// FindAll retorna todos os usuários cadastrados.
	FindAll() ([]*entity.User, error)

	// Update atualiza os dados de um usuário existente.
	Update(user *entity.User) error

	// Delete remove um usuário pelo ID.
	Delete(id string) error
}
