package satus

import (
	"fmt"
	"strings"

	"github.com/ternarybob/funktion"

	"github.com/gookit/color"
)

// GetScopedDataConnection ...
func (cfg *AppConfig) GetScopedDataConnection(name string) (DataConfig, error) {

	var connections map[string]DataConfig
	var connection DataConfig

	if funktion.IsEmpty(name) {
		fmt.Println(color.Warn.Render("Data Connection -> Name is Empty -> Returning"))
		return connection, nil
	}

	connections, err := cfg.getScopedDataConnections()

	if err != nil {
		fmt.Println(color.Warn.Render(err))
		return connection, err
	}

	return connections[name], nil

}

func (cfg *AppConfig) GetScopedDataConnectionbyType(t string) (DataConfig, error) {

	output := DataConfig{}

	if funktion.IsEmpty(t) {
		fmt.Println(color.Warn.Render("Data Connection -> Type is Empty -> Returning"))
		return output, nil
	}

	output, err := cfg.getConnectionbyType(t)
	if err != nil {
		fmt.Println(color.Warn.Render(err))
		return output, err
	}

	return output, nil

}

// GetScopedDataConnections ...
func (cfg *AppConfig) getScopedDataConnections() (map[string]DataConfig, error) {

	output := map[string]DataConfig{}
	scope := cfg.Service.Scope

	// Add Scoped Matches
	for _, d := range cfg.Connections {

		// Scope exact match
		if contains(scope, d.Scope) {
			output[d.Name] = d
		}
	}

	if len(output) <= 0 {
		return output, fmt.Errorf("No connection matched scope:%s", scope)
	}

	return output, nil
}

// GetConnectionbyType ...
func (cfg *AppConfig) getConnectionbyType(t string) (DataConfig, error) {

	var connections map[string]DataConfig

	output := DataConfig{}

	connections, err := cfg.getScopedDataConnections()
	if err != nil {
		fmt.Println(color.Warn.Render(err))
		return output, err
	}

	// Find first data connection with type t
	for _, d := range connections {

		// Type not matching
		if !strings.EqualFold(d.Type, t) {
			continue
		}

		return output, nil
	}

	return output, nil
}
