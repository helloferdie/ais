package response

import (
	"ais/lib/libresponse"
	"ais/pkg/article"
)

// Article -
func Article(at *article.Model, format map[string]interface{}) map[string]interface{} {
	m := libresponse.MapOutput(at, true, format)
	return m
}
