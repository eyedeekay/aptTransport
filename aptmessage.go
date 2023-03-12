package apttransport

import "strings"

// AptMessage is a message from the apt transport.
type AptMessage struct {
	Status     string
	StatusCode int
	Header     Header
	Exit       int
}

// Header is a simple map of headers for an Apt message
type Header map[string][]string

// Add a header to the message.
func (h Header) Add(key, value string) {
	h[key] = append(h[key], value)
}

// Get a header from the message.
func (h Header) Get(key string) string {
	if value, ok := h[key]; ok {
		if len(value) > 0 {
			return value[0]
		}
	}
	return ""
}

// String returns the apt message as a string
func (m *AptMessage) String() string {
	s := []string{m.Status}
	for k, values := range m.Header {
		for _, v := range values {
			s = append(s, k+": "+v)
		}
	}
	return strings.Join(s, "\n") + "\n\n"
}
