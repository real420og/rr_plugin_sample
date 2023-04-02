package plugin

import (
	"net/http"
)

type ResponseDecorator struct {
	http.ResponseWriter

	originalBody []byte
	code         int
	size         int
}

func NewResponseDecorator(w http.ResponseWriter) *ResponseDecorator {
	return &ResponseDecorator{
		ResponseWriter: w,
	}
}

func (r *ResponseDecorator) Write(body []byte) (int, error) {
	var err error

	r.originalBody = body

	body = rot13(body)

	r.size, err = r.ResponseWriter.Write(body)

	return r.size, err
}

func (r *ResponseDecorator) Code() int {
	return r.code
}

func (r *ResponseDecorator) Size() int {
	return r.size
}

func (r *ResponseDecorator) OriginalBody() []byte {
	return r.originalBody
}

func rot13(input []byte) []byte {
	result := make([]byte, len(input))
	for i, b := range input {
		switch {
		case b >= 'A' && b <= 'Z':
			result[i] = 'A' + (b-'A'+13)%26
		case b >= 'a' && b <= 'z':
			result[i] = 'a' + (b-'a'+13)%26
		default:
			result[i] = b
		}
	}
	return result
}
