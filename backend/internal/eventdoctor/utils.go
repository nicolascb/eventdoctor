package eventdoctor

import "fmt"

// validateOwnerEventsSchema checks that all events from owner producers have a schema_url defined.
func validateOwnerEventsSchema(producers []Producer) error {
	for _, producer := range producers {
		if producer.Owner {
			for _, event := range producer.Events {
				if event.SchemaURL == "" {
					return fmt.Errorf("producer is owner of topic %q but event %q has no schema_url defined",
						producer.Topic, event.Name)
				}
			}
		}
	}
	return nil
}

// validateUniqueEnvironments checks that each environment name appears only once.
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
