package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/nicolascb/eventdoctor/internal/api/response"
	"github.com/nicolascb/eventdoctor/internal/eventdoctor"
)

// TopicViewer is the interface for fetching topic and service views from the API.
type TopicViewer interface {
	GetTopicView(serverURL, topicName string) (*response.TopicView, error)
	GetServiceView(serverURL, serviceName string) (*response.ServiceView, error)
}

// GetCommand handles read queries against the EventDoctor API.
type GetCommand struct {
	client TopicViewer
}

func NewGetCommand(client TopicViewer) *GetCommand {
	return &GetCommand{client: client}
}

// ResolveServerURL returns the server URL to use for API calls.
// If directURL is non-empty it is used as-is; otherwise the URL is resolved
// from the eventdoctor.yml at filePath for the given environment.
func (c *GetCommand) ResolveServerURL(filePath, env, directURL string) (string, error) {
	if directURL != "" {
		return directURL, nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open config file %q: %w", filePath, err)
	}
	defer file.Close()

	return resolveURLFromReader(file, env)
}

func resolveURLFromReader(r io.Reader, env string) (string, error) {
	spec, err := eventdoctor.LoadSpecFromReader(r)
	if err != nil {
		return "", fmt.Errorf("failed to parse config file: %w", err)
	}

	return spec.GetServerURL(env)
}

// Topic fetches and prints the producer/consumer table for a topic.
func (c *GetCommand) Topic(serverURL, topicName string) error {
	view, err := c.client.GetTopicView(serverURL, topicName)
	if err != nil {
		return fmt.Errorf("failed to get topic %q: %w", topicName, err)
	}

	printTopicView(view)
	return nil
}

// Service fetches and prints the produce/consume table for a service.
func (c *GetCommand) Service(serverURL, serviceName string) error {
	view, err := c.client.GetServiceView(serverURL, serviceName)
	if err != nil {
		return fmt.Errorf("failed to get service %q: %w", serviceName, err)
	}

	printServiceView(view)
	return nil
}

var tableStyle = table.Style{
	Name: "EventDoctor",
	Box:  table.StyleBoxLight,
	Color: table.ColorOptions{
		Header: text.Colors{text.Bold},
	},
	Format: table.FormatOptions{
		Header: text.FormatUpper,
	},
	Options: table.OptionsDefault,
}

func newTable() table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(tableStyle)
	return t
}

func printTopicView(view *response.TopicView) {
	fmt.Printf("Topic: %s\n", view.Topic)
	if view.OwnerService != nil {
		fmt.Printf("Owner: %s\n", *view.OwnerService)
	}
	fmt.Println()

	fmt.Println("PRODUCERS")
	if len(view.Producers) == 0 {
		fmt.Println("  (none)")
	} else {
		t := newTable()
		t.AppendHeader(table.Row{"Service", "Repository", "Event", "Writes", "Owner"})
		for _, p := range view.Producers {
			t.AppendRow(table.Row{p.Service, p.Repository, p.Event, p.Writes, p.Owner})
		}
		t.Render()
	}

	fmt.Println()
	fmt.Println("CONSUMERS")

	t := newTable()
	t.AppendHeader(table.Row{"Service", "Repository", "Event", "Group", "Version"})
	for _, c := range view.Consumers {
		version := "-"
		if c.Version != nil {
			version = *c.Version
		}
		t.AppendRow(table.Row{c.Service, c.Repository, c.Event, c.Group, version})
	}

	if len(view.Consumers) == 0 {
		t.AppendRow(table.Row{"-", "-", "-", "-", "-"})
	}

	t.Render()
}

func printServiceView(view *response.ServiceView) {
	fmt.Printf("Service: %s\n\n", view.Service)

	fmt.Println("PRODUCES")

	t := newTable()
	t.AppendHeader(table.Row{"Topic", "Event", "Writes", "Owner"})
	for _, p := range view.Produces {
		t.AppendRow(table.Row{p.Topic, p.Event, p.Writes, p.Owner})
	}

	if len(view.Produces) == 0 {
		t.AppendRow(table.Row{"-", "-", "-", "-"})
	}

	t.Render()

	fmt.Println()
	fmt.Println("CONSUMES")

	tConsumers := newTable()
	tConsumers.AppendHeader(table.Row{"Topic", "Event", "Group", "Version"})
	for _, c := range view.Consumes {
		version := "-"
		if c.Version != nil {
			version = *c.Version
		}
		tConsumers.AppendRow(table.Row{c.Topic, c.Event, c.Group, version})
	}

	if len(view.Consumes) == 0 {
		tConsumers.AppendRow(table.Row{"-", "-", "-", "-"})
	}

	tConsumers.Render()
}
