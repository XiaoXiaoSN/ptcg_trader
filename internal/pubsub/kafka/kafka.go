package kafka

import "encoding/json"

// Config for kafka client
type Config struct {
	Debug          bool     `yaml:"debug"`
	BrokerList     []string `yaml:"address"`
	ClientID       string   `yaml:"client_id"`
	GroupName      string   `yaml:"group_name"`
	OffsetsInitial int64    `yaml:"offsets_initial"`
}

// NewInjection ...
func (c *Config) NewInjection() *Config {
	return c
}

// MsgData ...
type MsgData struct {
	TraceID   string          `json:"trace_id"`
	Data      json.RawMessage `json:"data,omitempty"`
	ConsumeID string          `json:"consume_id"`
}

// ProducerMessage ...
type ProducerMessage struct {
	Topic string
	Data  []byte
}

// ProducerError ...
type ProducerError struct {
	Msg ProducerMessage
	Err error
}
