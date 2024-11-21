package dto

import (
	"encoding/json"
)

type ServerResponse struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
	Ts   int64           `json:"ts"`
}
