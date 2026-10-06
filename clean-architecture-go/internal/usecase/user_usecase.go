package usecase

import (
	"clean-architecture-go/internal/entity"
	"errors"
)

// Erros da camada de aplicação — diferentes dos erros de domínio.
// Erros aqui representam falhas no fluxo da aplicação, não na regra de negócio.
var (
	ErrUserNotFound = errors.New("usuário não encontrado")
)

// UserUseCase implementa todos os casos de uso relacionados a usuários.
// Ela depende APENAS de interfaces (UserRepository), nunca de implementações concretas.
type UserUseCase struct {
	repo UserRepository
}

// NewUserUseCase cria uma nova instância do use case.
// O repositório é injetado via construtor (Dependency Injection).
// Isso permite trocar a implementação sem alterar este código.
func NewUserUseCase(repo UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

// CreateUser orquestra a criação de um novo usuário.
// Fluxo:
//  1. Cria a entidade (que valida as regras de negócio)
//  2. Persiste via repositório
//  3. Retorna o DTO de saída
func (uc *UserUseCase) CreateUser(input CreateUserInput) (*UserOutput, error) {
	// 1. Criar entidade — aqui as regras de negócio são validadas
	user, err := entity.NewUser(input.Name, input.Email, input.Age)
	if err != nil {
		return nil, err
	}

	// 2. Persistir — o use case não sabe COMO é persistido, apenas que É
	if err := uc.repo.Save(user); err != nil {
		return nil, err
	}

	// 3. Converter entidade para DTO de saída
	return toOutput(user), nil
}

// GetUserByID busca um usuário pelo ID.
func (uc *UserUseCase) GetUserByID(id string) (*UserOutput, error) {
	user, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return toOutput(user), nil
}

// ListUsers retorna todos os usuários.
func (uc *UserUseCase) ListUsers() ([]*UserOutput, error) {
	users, err := uc.repo.FindAll()
	if err != nil {
		return nil, err
	}

	output := make([]*UserOutput, len(users))
	for i, user := range users {
		output[i] = toOutput(user)
	}

	return output, nil
}

// UpdateUser atualiza os dados de um usuário existente.
func (uc *UserUseCase) UpdateUser(input UpdateUserInput) (*UserOutput, error) {
	// 1. Buscar o usuário existente
	user, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// 2. Aplicar alterações (a entidade revalida as regras de negócio)
	if err := user.Update(input.Name, input.Email, input.Age); err != nil {
		return nil, err
	}

	// 3. Persistir
	if err := uc.repo.Update(user); err != nil {
		return nil, err
	}

	return toOutput(user), nil
}

// DeleteUser remove um usuário pelo ID.
func (uc *UserUseCase) DeleteUser(id string) error {
	user, err := uc.repo.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	return uc.repo.Delete(id)
}

// toOutput converte uma entidade User para o DTO de saída.
// Essa conversão garante que a entidade nunca é exposta diretamente.
func toOutput(user *entity.User) *UserOutput {
	return &UserOutput{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Age:       user.Age,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
