# EventDoctor CLI

- Validate/Send configuration to the API
- Fetch information from the API


Build documentation at https://github.com/imfing/hextra.

## Commands

| Command                           | Description                                        |
| --------------------------------- | -------------------------------------------------- |
| `eventdoctor-cli config validate` | Validates the `eventdoctor.yml` configuration file |
| `eventdoctor-cli config apply`    | Applies the configurations defined in the file     |
| `eventdoctor-cli get`             | Retrieves information from the EventDoctor API     |

**Global flags:**

| Flag                      | Description                                               | Required |
| ------------------------- | --------------------------------------------------------- | -------- |
| `-h, --help`              | Displays help for the specified command                   | No       |
| `-v, --version`           | Displays the EventDoctor CLI version                      | No       |
| `-e, --env <environment>` | Specifies the environment (development, production, etc.) | Yes      |


### Command: `config`

| Command    | Options             | Description                                                                                                                                           |
| ---------- | ------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- |
| `validate` | `-f, --file <path>` | Validates the `eventdoctor.yml` configuration file. If no file is specified, the default is `./eventdoctor.yml`.                                      |
| `apply`    | `-f, --file <path>` | Applies the configurations defined in the `eventdoctor.yml` file to the EventDoctor API. If no file is specified, the default is `./eventdoctor.yml`. |


### Command: `get`

The `eventdoctor-cli get` command allows retrieving information from the EventDoctor API. It can be used to list events, topics, or other entities managed by the API.

At least one filter must be provided for the search to be performed.

The EventDoctor API URL must be previously configured in the `~/.eventdoctor` configuration file.

```yaml
  servers:
    - environment: "development"
      url: "http://localhost:8080"
    - environment: "production"
      url: "https://eventdoctor.empresa.com"
```


| Flag          | Description        | Required | Example                                                        |
| ------------- | ------------------ | -------- | -------------------------------------------------------------- |
| `--topics `   | Filter by topics   | No       | `eventdoctor-cli get --topics=user.events,payments.events `    |
| `--services ` | Filter by services | No       | `eventdoctor-cli get --services=user-service,payment-service ` |
