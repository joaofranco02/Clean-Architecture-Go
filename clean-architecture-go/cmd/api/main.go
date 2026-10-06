// Package main é o ponto de entrada da aplicação (Composition Root).
//
// É aqui que acontece a "fiação" (wiring) — onde todas as camadas são
// conectadas por meio de Injeção de Dependência.
//
// Este é o ÚNICO lugar que conhece todas as implementações concretas.
// Todas as outras camadas dependem apenas de interfaces.
//
// Fluxo de composição:
//   Repositório (infra) → Use Case (application) → Handler (adapter) → HTTP Server (framework)
package main

import (
	"clean-architecture-go/internal/handler"
	"clean-architecture-go/internal/infra/repository"
	"clean-architecture-go/internal/usecase"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// ┌─────────────────────────────────────────────────────────┐
	// │              COMPOSIÇÃO / DEPENDENCY INJECTION          │
	// │                                                         │
	// │  Aqui montamos o grafo de dependências da aplicação.    │
	// │  Se quisermos trocar o banco de dados, mudamos APENAS   │
	// │  esta linha — nada mais no sistema precisa mudar.       │
	// └─────────────────────────────────────────────────────────┘

	// Camada 4: Infraestrutura — Implementação concreta do repositório
	// Para usar PostgreSQL, bastaria trocar por:
	//   userRepo := repository.NewUserPostgresRepository(db)
	userRepo := repository.NewUserMemoryRepository()

	// Camada 2: Use Cases — Recebe o repositório via interface
	userUseCase := usecase.NewUserUseCase(userRepo)

	// Camada 3: Handlers — Recebe o use case
	userHandler := handler.NewUserHandler(userUseCase)

	// Camada 4: Framework — Servidor HTTP (stdlib do Go)
	mux := http.NewServeMux()
	userHandler.RegisterRoutes(mux)

	// Middleware de logging
	loggedMux := loggingMiddleware(mux)

	port := ":8080"
	fmt.Println("╔══════════════════════════════════════════════════╗")
	fmt.Println("║    🏛️  Clean Architecture API — Go              ║")
	fmt.Println("║                                                  ║")
	fmt.Println("║    Endpoints:                                    ║")
	fmt.Println("║    GET    /users       → Listar usuários         ║")
	fmt.Println("║    POST   /users       → Criar usuário           ║")
	fmt.Println("║    GET    /users/{id}  → Buscar por ID           ║")
	fmt.Println("║    PUT    /users/{id}  → Atualizar usuário       ║")
	fmt.Println("║    DELETE /users/{id}  → Remover usuário         ║")
	fmt.Println("║                                                  ║")
	fmt.Printf("║    🚀 Servidor rodando em http://localhost%s  ║\n", port)
	fmt.Println("╚══════════════════════════════════════════════════╝")

	if err := http.ListenAndServe(port, loggedMux); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}

// loggingMiddleware é um middleware simples que loga cada requisição.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("➡️  %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
