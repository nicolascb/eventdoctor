# TODO

## Projeto

- [x] Criar pacote de logger usando slog
- [ ] Criar `Makefile` para comandos comuns (build, test, lint, etc.)
- [ ] Configurar linter (golangci-lint)
- [ ] Definir política de versionamento (semver)

## CLI

- [x] `eventdoctor-cli config validate`
- [x] `eventdoctor-cli config apply`
- [ ] `eventdoctor-cli get --topics --services`
- [ ] Criar testes unitários

## Subscriber

- [ ] Verificar tópicos e consumers de X em X tempo
- [ ] Criar testes unitários

## API

- [x] Substituir o Gin pelo net/http
- [x] Substituir o zap.Logger pelo slog
- [ ] Expor métricas prometheus
- [ ] Adicionar health check endpoints (`/health`, `/ready`)
- [ ] Adicionar autenticação/autorização (API keys)
- [ ] Implementar rate limiting
- [ ] Documentação OpenAPI/Swagger
- [ ] Criar testes unitários
- [ ] Criar testes de integração
- [x] Criar drivers para outros bancos de dados (PostgreSQL, MySQL, etc.)

## Observabilidade

- [ ] Adicionar tracing distribuído (OpenTelemetry)
- [ ] Centralizar logs (formato JSON estruturado)

## Infra

- [ ] Dockerizar a API
- [ ] Dockerizar a CLI
- [ ] Dockerizar o subscriber
- [ ] Criar `docker-compose.yaml` para desenvolvimento local
- [ ] Criar `.goreleaser.yaml` para automatizar releases
- [ ] Configurar GitHub Actions (CI/CD)
- [ ] Scan de vulnerabilidades nas dependências
- [ ] Scan de vulnerabilidades nas imagens Docker
- [ ] Disponibilizar CLI via home brew

## Documentação

- [ ] Criar repositório eventdoctor-docs
- [ ] Adicionar CONTRIBUTING.md
- [ ] Adicionar CHANGELOG.md
- [ ] Documentar variáveis de ambiente
- [ ] Criar exemplos de uso mais completos 