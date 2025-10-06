# EventDoctor - Plano e Backlog Técnico

## Sumário
- [EventDoctor - Plano e Backlog Técnico](#eventdoctor---plano-e-backlog-técnico)
  - [Sumário](#sumário)
  - [1. Visão Geral](#1-visão-geral)
  - [2. Módulos](#2-módulos)
    - [2.1 Módulo Config](#21-módulo-config)
    - [2.2 Módulo Docs](#22-módulo-docs)
    - [2.3 Módulo Subscriber](#23-módulo-subscriber)
  - [3. Modelo de Dados](#3-modelo-de-dados)
    - [3.1 Tabelas Principais](#31-tabelas-principais)
    - [3.2 Tabelas de Observabilidade](#32-tabelas-de-observabilidade)
    - [3.3 Histórico / Revisões](#33-histórico--revisões)
    - [3.4 Índices Sugeridos](#34-índices-sugeridos)
    - [3.5 Soft Delete / Hard Delete](#35-soft-delete--hard-delete)
  - [4. Regras de Validação (Formais)](#4-regras-de-validação-formais)
    - [4.1 Spec Root](#41-spec-root)
    - [4.2 Producers](#42-producers)
    - [4.3 Events](#43-events)
    - [4.4 Consumers](#44-consumers)
    - [4.5 Schema URL](#45-schema-url)
    - [4.6 Warnings Possíveis](#46-warnings-possíveis)
  - [5. API (Fase 1)](#5-api-fase-1)
  - [6. Módulo Config - Contrato Interno](#6-módulo-config---contrato-interno)
  - [7. Módulo Docs - Geração](#7-módulo-docs---geração)
  - [8. Módulo Subscriber - Fluxo](#8-módulo-subscriber---fluxo)
  - [9. Edge Cases e Decisões](#9-edge-cases-e-decisões)
  - [10. Roadmap (Milestones)](#10-roadmap-milestones)
  - [11. Módulo CLI](#11-módulo-cli)
    - [Visão Geral](#visão-geral)
    - [Comandos Principais](#comandos-principais)
    - [Parâmetros/Flags](#parâmetrosflags)
    - [Exemplo de Uso](#exemplo-de-uso)
    - [Considerações](#considerações)
  - [12. Módulo de Notificações Slack](#12-módulo-de-notificações-slack)
    - [Visão Geral](#visão-geral-1)
    - [Gatilhos de Notificação](#gatilhos-de-notificação)
    - [Pontos de Integração](#pontos-de-integração)
    - [Integração com Slack](#integração-com-slack)
    - [Formato das Mensagens](#formato-das-mensagens)
    - [Configuração](#configuração)
    - [Exemplo de Payload](#exemplo-de-payload)
    - [Roadmap de Integração](#roadmap-de-integração)
  - [13. Próximos Passos Concretos](#13-próximos-passos-concretos)

---

## 1. Visão Geral
Ferramenta para governança de eventos (Kafka): catálogo de producers, events, consumers; validação dinâmica; documentação e cobertura; detecção de lacunas.

---

## 2. Módulos

### 2.1 Módulo Config
Entrada: YAML (eventdoctor.yml)  
Saídas: validação estruturada, persistência normalizada, criação de revisão (hash).  
Responsabilidades:
- Parse + validação
- Aplicar delta (upsert / soft delete)
- Registrar revisão (hash deterministic via JSON ordenado)

### 2.2 Módulo Docs
Geração de:
- HTML (base /docs) + páginas de perfil
- Cobertura, gaps (eventos sem consumidor, consumidores de eventos inexistentes, tópicos sem owner, versões deprecated)
- Templates reutilizáveis

### 2.3 Módulo Subscriber
Responsável por:
- Descoberta de tópicos Kafka
- Registro de tópicos/eventos “missing”
- Consumo + validação via JSON Schema
- Armazenar logs e lacunas

---

## 3. Modelo de Dados

### 3.1 Tabelas Principais
| Tabela | Campos (principais) | Observações |
|--------|---------------------|-------------|
| producers | id, name (UNIQUE), topic, owner (BOOL), title, description, created_at, updated_at, removed_at NULL | Unicidade name |
| events | id, topic, type, version NULL, description, schema_url NULL, deprecated BOOL DEFAULT false, removed_at NULL | UNIQUE(topic,type,COALESCE(version,'')) |
| consumers | id, service (UNIQUE), "group", description, created_at, updated_at, removed_at NULL | service único |
| consumer_topics | id, consumer_id FK, topic | UNIQUE(consumer_id, topic) |
| consumer_topic_events | id, consumer_topic_id FK, event_type, event_version NULL | UNIQUE(consumer_topic_id, event_type, COALESCE(event_version,'')) |

### 3.2 Tabelas de Observabilidade
| Tabela | Função |
|--------|--------|
| missing_topics | Tópicos vistos no cluster mas ausentes em config |
| missing_events | Eventos (topic + type) observados e não configurados |
| validation_log | Log de validação (status, erro, sample, etc.) |
| schema_cache | Cache de schemas (etag, sucesso, conteúdo) |

### 3.3 Histórico / Revisões
| Tabela | Campos |
|--------|--------|
| config_revisions | id, hash, applied_at, raw_yaml, producers_count, events_count, consumers_count |

### 3.4 Índices Sugeridos
- events(topic, type)
- validation_log(topic, event_type, created_at)
- missing_events(topic, event_type)

### 3.5 Soft Delete / Hard Delete
- Inicial: hard delete para producers/consumers; soft para events (permitindo histórico de versões).
- Alternativa futura: uniformizar com removed_at.

---

## 4. Regras de Validação (Formais)

### 4.1 Spec Root
- version: obrigatório (string)
- metadata: opcional
  - Se tiver server_url ou repository → ambos obrigatórios
  - Padronizar uso de metadata.server_url (eliminar metadata.server)

### 4.2 Producers
- name: obrigatório, único
- topic: obrigatório
- owner: obrigatório (bool)
- owner=true → title obrigatório; cada evento precisa de schema_url
- owner=false → schema_url ignorado (warning), title opcional
- Máx. 1 producer owner por topic

### 4.3 Events
- type: obrigatório
- version: opcional (regex semver ^\d+\.\d+\.\d+$ se presente)
- Duplicidade (topic,type,version) não pode divergir em schema_url
- Múltiplas versões permitidas
- deprecated=true → version obrigatório
- version ausente → armazenar NULL

### 4.4 Consumers
- service: obrigatório, único
- group: obrigatório
- topics[]: obrigatório
- Cada tópico referenciado deve existir em producers (ao menos 1)
- events[].type: obrigatório
- events[].version (se fornecida) deve existir; se não fornecida → aceita qualquer versão existente
- Referência a evento inexistente → erro

### 4.5 Schema URL
- Obrigatório quando exigido por regra acima
- Deve ser http(s) válido (sintaxe)
- Fetch não obrigatório na fase de config (feito no Subscriber)

### 4.6 Warnings Possíveis
- Producer non-owner com schema_url → ignorado
- Consumer aponta tópico sem owner (política: atualmente NÃO permitir; pode virar erro)
- Event sem versão (incentivar versionamento)

---

## 5. API (Fase 1)

Validação & Config:
- POST /config/validate → { valid, errors[], warnings[], summary }
- POST /config/apply → aplica se válido → retorna revision { id, hash }
- GET /config/revisions
- GET /config/revisions/{id}

Consulta Domínio:
- GET /producers
- GET /producers/{name}
- GET /consumers
- GET /consumers/{service}
- GET /topics
- GET /topics/{topic}/events
- GET /events?topic=&type=&version=
- GET /coverage → { topicsWithoutOwner, eventsWithoutConsumers, consumedUnknownEvents, missingTopics, missingEvents }

Observabilidade:
- GET /validations?status=&topic=&event_type=&limit=
- (Futuro) GET /schemas/cache

Docs:
- GET /docs
- GET /docs/producers/{name}
- GET /docs/consumers/{service}

---

## 6. Módulo Config - Contrato Interno
Funções:
- Parse(rawYAML string) (ConfigSpec, error)
- Validate(spec ConfigSpec) (ValidationResult)  
  - ValidationResult { Errors[], Warnings[], Summary{ producers, events, consumers } }
- Apply(spec ConfigSpec) (Revision, error)

Fluxo Apply:
1. Normaliza → hash (sha256 JSON ordenado)
2. Inicia transação
3. Upsert producers/events/consumers (reativa se soft-deletado)
4. Marca removed_at (ou hard delete) para ausentes
5. Insere config_revisions
6. Commit

---

## 7. Módulo Docs - Geração
Páginas:
- Index: tabela de tópicos (owner, #events, #consumers)
- Producer profile: eventos (deprecated, sem schema, versões)
- Consumer profile: tópicos + eventos consumidos (highlight removidos/deprecated)
- Coverage badges:
  - % eventos com consumidor
  - % tópicos com owner
  - # missing (topics/events)
Integrações:
- missing_topics / missing_events → seção “Observações”

---

## 8. Módulo Subscriber - Fluxo

Config:
- bootstrap_servers
- group_id_base (ex: eventdoctor-validator)
- poll_timeout
- schema_fetch_timeout
- max_concurrency

Processo:
1. Listar tópicos Kafka (Admin)
2. Se tópico não existe em producers → upsert em missing_topics
3. Criar consumer group único assinando todos (ou subset + missing)
4. Consumo:
   - Parse JSON
   - Extrair type (campo 'type' configurável) e version (opcional)
   - Resolver evento esperado (match exato se version; senão: existe qualquer?)
   - Inexistente → missing_events + validation_log(status='unknown_event')
   - Com schema_url → buscar schema (cache)
     - Falha fetch → status='schema_fetch_error'
     - Validação OK / mismatch → status correspondente
   - Armazenar validation_log
Campos futuros: payload_size, processing_ms

---

## 9. Edge Cases e Decisões
- Consumer pode referenciar tópico sem owner? (atual: NÃO)
- Estratégia de delete: iniciar hard (producers/consumers) + soft (events)
- Event sem versão: permitido (warning)
- metadata.server_url vs metadata.server: padronizar server_url
- Seleção de versão quando consumer não especifica:
  - Simples: aceitar qualquer existente (não precisa escolher maior)
  - Validação dinâmica no Subscriber usa conjunto de versões

---

## 10. Roadmap (Milestones)
M1 - Banco & Migração  
M2 - Módulo Config (structs, parser, validação, apply + testes)  
M3 - API Básica (config + listagens + coverage)  
M4 - Docs (templates índice, producers, consumers)  
M5 - Subscriber Infra (Kafka + missing topics)  
M6 - Validação de Mensagens (schema cache + JSON Schema + logs)  
M7 - Observabilidade (filtros /validations + métricas Prometheus)  
M8 - Hardening (retry/backoff schema fetch, circuit breaker, paginação)

---

## 11. Módulo CLI

### Visão Geral
A CLI do EventDoctor será responsável por facilitar a integração dos clientes (usuários e pipelines CI/CD) com a API do EventDoctor, permitindo o envio e validação de arquivos de configuração YAML de forma simples e automatizável.

### Comandos Principais
- `validate`           : Valida o arquivo localmente (opcional, se houver lógica local)
- `validate-remote`    : Envia o arquivo para a API `/config/validate` e exibe o resultado
- `apply`              : Envia o arquivo para a API `/config/apply` para aplicar a configuração
- `revisions`          : Lista revisões aplicadas (`GET /config/revisions`)

### Parâmetros/Flags
- `--file`, `-f`       : Caminho do arquivo YAML (default: `eventdoctor.yaml`)
- `--api-url`          : URL da API (pode ler de env ou config)
- `--token`            : Token de autenticação (opcional)
- `--json`             : Saída em formato JSON (para uso em scripts/pipelines)

### Exemplo de Uso
```bash
eventdoctor apply -f ./eventdoctor.yaml --api-url http://localhost:8080
```

### Considerações
- A CLI deve ser fácil de instalar (binário Go, sem dependências extras)
- Saída amigável para humanos e scripts
- Ideal para uso em pipelines CI/CD
- Estrutura sugerida: `/cmd/cli/`
- Frameworks Go sugeridos: `cobra` ou `urfave/cli`

---

## 12. Módulo de Notificações Slack

### Visão Geral
Notificações automáticas no Slack para eventos relevantes do catálogo: criação, remoção, alteração de schema e eventos não documentados. Facilita visibilidade e resposta rápida a mudanças.

### Gatilhos de Notificação
- Criação de tópico, evento, consumer ou producer
- Remoção de tópico, evento, consumer ou producer
- Alteração de schema de evento
- Detecção de evento não documentado (missing_events)

### Pontos de Integração
- Hooks/disparadores após operações de apply/config (CRUD)
- Após detecção de missing_events pelo Subscriber

### Integração com Slack
- Utilizar Slack Webhook (configurável via env, arquivo ou banco)
- Mensagens customizadas por tipo de evento
- Permitir ativar/desativar notificações por tipo

### Formato das Mensagens
- Criação: "Novo [tipo] criado: [nome] ([detalhes])"
- Remoção: "[Tipo] removido: [nome] em [data]"
- Alteração de schema: "Schema alterado para evento [topic/type/version]"
- Evento não documentado: "Evento não documentado detectado: [topic/type]"
- Incluir links para docs quando possível

### Configuração
- Endpoint do webhook configurável
- Flags para ativar/desativar tipos de notificação
- Possível extensão futura: múltiplos canais

### Exemplo de Payload
```json
{
  "text": "Novo evento criado: order.created (topic: orders)"
}
```

### Roadmap de Integração
- Adicionar hooks após apply/config e detecção de missing_events
- Implementar função utilitária para envio ao Slack
- Adicionar configuração de webhook e flags
- Testes de integração

---

## 13. Próximos Passos Concretos
1. Ajustar migration.sql com novas tabelas
2. Definir structs Go (internal/config/spec.go, internal/store/models.go)
3. Implementar parser + validator
4. Criar endpoints /config/validate e /config/apply
5. Gerar primeira página de docs expandida

---