package eventdoctor

import "fmt"

// validateOwnerEventsSchema verifica se todos os eventos de produtores owner têm schema_url definido
func validateOwnerEventsSchema(producers []Producer) error {
	for _, producer := range producers {
		if producer.Owner {
			for _, event := range producer.Events {
				if event.SchemaURL == "" {
					return fmt.Errorf("o produtor é owner do tópico '%s', mas o evento '%s' não possui schema_url definido",
						producer.Topic, event.Name)
				}
			}
		}
	}
	return nil
}

// validateUniqueEnvironments verifica se cada ambiente é único
func validateUniqueEnvironments(servers []Server) error {
	environments := make(map[string]bool)

	for _, server := range servers {
		if environments[server.Environment] {
			return fmt.Errorf("env duplicate: '%s'", server.Environment)
		}
		environments[server.Environment] = true
	}

	return nil
}
