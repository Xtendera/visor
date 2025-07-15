package validations

import (
	"github.com/santhosh-tekuri/jsonschema/v6"
	"io"
)

func ValidateSchemaFromPath(reader io.Reader, schemaPath string) error {
	c := jsonschema.NewCompiler()
	sch, err := c.Compile(schemaPath)
	if err != nil {
		return err
	}
	inst, err := jsonschema.UnmarshalJSON(reader)
	if err != nil {
		return err
	}
	err = sch.Validate(inst)
	if err != nil {
		return err
	}
	return nil
}
