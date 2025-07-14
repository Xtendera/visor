package validations

import (
	"github.com/santhosh-tekuri/jsonschema/v6"
	"io"
)

func ValidateSchemaFromFile(reader io.Reader, schemaFile string) error {
	c := jsonschema.NewCompiler()
	sch, err := c.Compile(schemaFile)
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
