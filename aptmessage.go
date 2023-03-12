package apttransport

import "strings"

type Message struct {
	Status     string
	StatusCode int
	Header     Header
	Exit       int
}

type AptMessage struct {
	Code    int
	Headers map[string]string
}

type Header map[string][]string

func (h Header) Add(key, value string) {
	h[key] = append(h[key], value)
}

func (h Header) Get(key string) string {
	if value, ok := h[key]; ok {
		if len(value) > 0 {
			return value[0]
		}
	}
	return ""
}

func (m *Message) String() string {
	s := []string{m.Status}
	for k, values := range m.Header {
		for _, v := range values {
			s = append(s, k+": "+v)
		}
	}
	return strings.Join(s, "\n") + "\n\n"
}
