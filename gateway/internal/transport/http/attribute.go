package http

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
)

func JsonAttribute(key string, value any) attribute.KeyValue {
	return attribute.Key(key).String(fmt.Sprintf("%v", value))
}
