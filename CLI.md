# EventDoctor CLI

- Validar/Enviar configuração para a API
- Buscar informações da API


Construir a documentação em https://github.com/imfing/hextra.

## Comandos

| Comando                           | Descrição                                          |
| --------------------------------- | -------------------------------------------------- |
| `eventdoctor-cli config validate` | Valida o arquivo de configuração `eventdoctor.yml` |
| `eventdoctor-cli config apply`    | Aplica as configurações definidas no arquivo       |
| `eventdoctor-cli get`             | Recupera informações da API do EventDoctor         |

**Flags globais:**

| Flag                      | Descrição                                               | Obrigatório |
| ------------------------- | ------------------------------------------------------- | ----------- |
| `-h, --help`              | Exibe a ajuda para o comando especificado               | Não         |
| `-v, --version`           | Exibe a versão do EventDoctor CLI                       | Não         |
| `-e, --env <environment>` | Especifica o ambiente (desenvolvimento, produção, etc.) | Sim         |


### Comando: `config`

| Comando    | Opções              | Descrição                                                                                                                                                     |
| ---------- | ------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `validate` | `-f, --file <path>` | Valida o arquivo de configuração `eventdoctor.yml`. Se nenhum arquivo for especificado, o padrão é `./eventdoctor.yml`.                                       |
| `apply`    | `-f, --file <path>` | Aplica as configurações definidas no arquivo `eventdoctor.yml` para a API do EventDoctor. Se nenhum arquivo for especificado, o padrão é `./eventdoctor.yml`. |


### Comando: `get`

O comando `eventdoctor-cli get` permite recuperar informações da API do EventDoctor. Ele pode ser usado para listar eventos, tópicos ou outras entidades gerenciadas pela API.

Ao menos um filtro deve ser fornecido para que a busca seja efetivada.

A URL da API do EventDoctor deve ser configurada previamente no arquivo de configuração `~/.eventdoctor`.

```yaml
  servers:
    - environment: "development"
      url: "http://localhost:8080"
    - environment: "production"
      url: "https://eventdoctor.empresa.com"
```


| Flag          | Descrição           | Obrigatório | Exemplo                                                        |
| ------------- | ------------------- | ----------- | -------------------------------------------------------------- |
| `--topics `   | Filtro por tópicos  | Não         | `eventdoctor-cli get --topics=user.events,payments.events `    |
| `--services ` | Filtro por serviços | Não         | `eventdoctor-cli get --services=user-service,payment-service ` |
