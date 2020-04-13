package models

import (
	"encoding/json"
	"fmt"
)

type ValidationError struct {
	Messages map[string][]string
}

func (e *ValidationError) Error() string {
	buf, _ := json.Marshal(e.Messages)
	return fmt.Sprintf("Validation error: %s", string(buf))
}

type Base struct {
	identified bool
}

func (a *Base) Identify() {
	a.identified = true
}

func (a *Base) Identified() bool {
	return a.identified
}
