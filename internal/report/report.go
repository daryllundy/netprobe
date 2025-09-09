package report

import (
	"encoding/json"

	"github.com/daryllundy/netprobe/internal/logger"
)

type Reporter struct {
	logger logger.Logger
}

func New(log logger.Logger) *Reporter {
	return &Reporter{logger: log}
}

func (r *Reporter) GenerateJSON(data interface{}) ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}
