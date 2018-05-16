package utils

import (
	"encoding/json"
	"io"

	"github.com/qiangxue/fasthttp-routing"
)

func DecodeJson(i interface{}, r io.Reader) (error) {
	dec := json.NewDecoder(r)
	return dec.Decode(i)
}

func EncodeJson(i interface{}, wr io.Writer) (error) {
	enc := json.NewEncoder(wr)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(i)
}

func String(c *routing.Context, k string) (string) {
	return string(c.FormValue(k))
}

func Bool(c *routing.Context, k string) (bool) {
	return String(c, k) == "true"
}

func StringOp(c *routing.Context, k string) (*string) {
	s := String(c, k)
	if s == "" {
		return nil
	}
	return &s
}

func BoolOp(c *routing.Context, k string) (*bool) {
	s := String(c, k)
	if s == "" {
		return nil
	}
	b := s == "true"
	return &b
}

func Or(exp bool, a, b string) string {
	if exp {
		return a
	}
	return b
}
