# EventDoctor

Projeto de código aberto para gerenciamento e documentação de eventos em sistemas distribuídos.

É capaz de gerar documentação automática, validar conformidade de eventos e facilitar a comunicação entre produtores e consumidores de eventos.

- [EventDoctor](#eventdoctor)
  - [Plataformas suportadas](#plataformas-suportadas)
  - [API Server](#api-server)
  - [Subscriber](#subscriber)
  - [Produtores e Consumidores](#produtores-e-consumidores)
    - [Visão Geral](#visão-geral)
    - [Estrutura do Arquivo](#estrutura-do-arquivo)
    - [Especificação Detalhada dos Campos](#especificação-detalhada-dos-campos)
      - [Campos Globais](#campos-globais)
      - [Producers](#producers)
      - [Consumers](#consumers)
    - [Exemplos de Uso](#exemplos-de-uso)
      - [Produtor Simples](#produtor-simples)
      - [Consumidor Múltiplos Tópicos](#consumidor-múltiplos-tópicos)
      - [Produtor Não-Owner](#produtor-não-owner)
    - [Validações](#validações)
    - [Casos de Uso Avançados](#casos-de-uso-avançados)
      - [Versionamento de Eventos](#versionamento-de-eventos)
  - [CI/CD](#cicd)


## Plataformas suportadas

- Kafka

## API Server

O eventdoctor-api é um serviço que expõe uma API REST para interagir com o EventDoctor. Ele permite que os usuários consultem informações sobre eventos, produtores e consumidores, além de fornecer endpoints para validação de eventos e geração de documentação.


## Subscriber

O eventdoctor-subscriber é um serviço que consome mensagens de um broker (como Kafka) e valida os eventos recebidos contra os schemas definidos no arquivo de configuração `eventdoctor.yml`. Ele registra logs detalhados sobre a conformidade dos eventos, ajudando a garantir que todos os eventos estejam em conformidade com as especificações definidas.

## Produtores e Consumidores

### Visão Geral

O EventDoctor utiliza um arquivo de configuração YAML para definir produtores e consumidores de eventos. Este arquivo serve como um registro central de todos os eventos em seu sistema, permitindo rastreamento, versionamento e documentação automática.

Cada projeto que produz ou consome eventos deve ter um arquivo `eventdoctor.yml`.

### Estrutura do Arquivo

```yaml
# eventdoctor.yml
version: "1.0"
metadata:
  server_url: "https://eventdoctor.empresa.com/api"
  repository: "https://github.com/empresa/chat-microservice"

producers:
  - name: "chat-microservice"
    topic: "chat-microservice.events"
    owner: true
    title: "Chat Events"
    description: "Eventos relacionados ao sistema de chat"
    events:
      - type: "ChatCreated"
        version: "1.0.0"
        description: "Disparado quando um novo chat é criado"
        schema_url: "https://gitlab.com/nicolascorrea/eventdoctor/schemas/chat_created.json"
      - type: "ChatDeleted"
        version: "1.0.0"
        description: "Disparado quando um chat é removido"
        schema_url: "https://gitlab.com/nicolascorrea/eventdoctor/schemas/chat_deleted.json"

consumers:
  - service: "notification-microservice"
    group: "notification-group"
    description: "Serviço responsável por notificações push"
    topics: 
      - name: "chat-microservice.events"
        events: 
          - type: "ChatCreated"
          - type: "ChatDeleted"
      - name: "user-microservice.events"
        events:
          - type: "UserRegistered"
```

### Especificação Detalhada dos Campos

#### Campos Globais

| Campo                | Tipo   | Descrição                              | Obrigatório | Padrão |
| -------------------- | ------ | -------------------------------------- | ----------- | ------ |
| version              | string | Versão da especificação                | sim         | -      |
| metadata             | object | Metadados do projeto                   | não         | -      |
| metadata.repository  | string | URL do repositório                     | sim         | -      |
| metadata.server      | string | URL do servidor EventDoctor            | sim         | -      |

#### Producers

| Campo                     | Tipo    | Descrição                                    | Obrigatório           | Exemplo                    |
| ------------------------- | ------- | -------------------------------------------- | --------------------- | -------------------------- |
| name                      | string  | Nome único do serviço produtor               | sim                   | "chat-microservice"        |
| topic                     | string  | Nome do tópico                               | sim                   | "chat-microservice.events" |
| owner                     | boolean | Se o serviço é proprietário do tópico        | sim                   | true                       |
| title                     | string  | Título legível do tópico                     | sim se owner for true | "Chat Events"              |
| description               | string  | Descrição do propósito do tópico             | não                   | "Eventos do sistema..."    |
| events                    | array   | Lista de eventos produzidos                  | sim                   | -                          |
| events[].type             | string  | Tipo/nome do evento                          | sim                   | "ChatCreated"              |
| events[].version          | string  | Versão do schema do evento                   | não                   | "1.0.0"                    |
| events[].description      | string  | Descrição do evento                          | não                   | "Disparado quando..."      |
| events[].schema_url       | string  | URL do JSON Schema do evento                 | sim se owner for true | "https://..."              |

#### Consumers

| Campo                     | Tipo   | Descrição                                    | Obrigatório | Exemplo                         |
| ------------------------- | ------ | -------------------------------------------- | ----------- | ------------------------------- |
| service                   | string | Nome único do serviço consumidor             | sim         | "notification-microservice"     |
| group                     | string | Grupo de consumidores (Kafka consumer group) | sim         | "notification-group"            |
| description               | string | Descrição do serviço                         | não         | "Serviço de notificações..."    |
| topics                    | array  | Lista de tópicos consumidos                  | sim         | -                               |
| topics[].name             | string | Nome do tópico                               | sim         | "chat-microservice.events"      |
| topics[].events           | array  | Lista de eventos consumidos                  | sim         | -                               |
| topics[].events[].type    | string | Tipo do evento                               | sim         | "ChatCreated"                   |

### Exemplos de Uso

#### Produtor Simples
```yaml
producers:
  - name: "user-service"
    topic: "user-service.events"
    owner: true
    title: "User Events"
    events:
      - type: "UserRegistered"
        schema_url: "https://schemas.empresa.com/user_registered.json"
```

#### Consumidor Múltiplos Tópicos
```yaml
consumers:
  - service: "analytics-service"
    group: "analytics-group"
    topics:
      - name: "user-service.events"
        events:
          - type: "UserRegistered"
          - type: "UserDeleted"
      - name: "chat-service.events"
        events:
          - type: "ChatCreated"
```

#### Produtor Não-Owner
```yaml
producers:
  - name: "gateway-service"
    topic: "user-service.events"  # Tópico de outro serviço
    owner: false
    events:
      - type: "UserLoggedIn"  # Apenas produz, não define schema
```

### Validações

1. **Tópicos únicos por owner**: Apenas um serviço pode ser owner de um tópico
2. **Schema obrigatório para owners**: Serviços owner devem fornecer schema_url
3. **Nomes únicos**: Nomes de produtores e consumidores devem ser únicos
4. **Referências válidas**: Eventos consumidos devem existir em algum produtor
5. **Formato de versionamento**: Versões devem seguir semantic versioning (x.y.z)

### Casos de Uso Avançados

#### Versionamento de Eventos
```yaml
producers:
  - name: "order-service"
    topic: "order-service.events"
    owner: true
    title: "Order Events"
    events:
      - type: "OrderCreated"
        version: "2.0.0"
        description: "Nova versão com campo adicional 'priority'"
        schema_url: "https://schemas.empresa.com/order_created_v2.json"
      - type: "OrderCreated"
        version: "1.0.0"
        description: "Versão legada (deprecated)"
        schema_url: "https://schemas.empresa.com/order_created_v1.json"
        deprecated: true
```

## CI/CD