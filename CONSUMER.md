# EventDoctor Consumer

O `eventdoctor-consumer` é uma ferramenta de linha de comando que atua como um consumidor de eventos Kafka.

Roadmap de funcionalidades planejadas:

| Funcionalidade                        | Descrição                                                                          | Status |
| ------------------------------------- | ---------------------------------------------------------------------------------- | ------ |
| Base de tópicos não documentados      | Enriquecimento da base de tópicos não documentados no EventDoctor                  | -      |
| Base de eventos não documentados      | Enriquecimento da base de eventos não documentados no EventDoctor                  | -      |
| Base de consumidores não documentados | Enriquecimento da base de consumidores não documentados no EventDoctor             | -      |
| Validação de eventos Kafka            | Consome eventos de tópicos Kafka e valida contra esquemas definidos no EventDoctor | -      |


## Consumidores/Eventos/Tópicos não documentados

É considerado:

- Tópico sem `owner` definido no arquivo de configuração do EventDoctor.
- Evento publicado em um tópico que não está documentado no EventDoctor.
- Consumidor que consome de um tópico que não está documentado no EventDoctor.