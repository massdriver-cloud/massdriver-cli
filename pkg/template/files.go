package template

import (
	"encoding/json"
	"os"

	"github.com/osteele/liquid"
)

func WriteToFile(filePath string, template []byte, data *Data) error {
	engine := liquid.NewEngine()

	var bindings map[string]interface{}
	inrec, _ := json.Marshal(data)
	json.Unmarshal(inrec, &bindings)

	out, renderErr := engine.ParseAndRender(template, bindings)

	if renderErr != nil {
		return renderErr
	}

	return os.WriteFile(filePath, out, 0600)
}
