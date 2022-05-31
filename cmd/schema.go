package cmd

import (
	"os"

	"github.com/massdriver-cloud/massdriver-cli/pkg/jsonschema"

	"github.com/spf13/cobra"
)

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Manage JSON Schemas",
	Long:  ``,
}

var schemaValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validates an input JSON object against a JSON schema",
	Long:  ``,
	RunE:  runSchemaValidate,
}

var schemaDereferenceCmd = &cobra.Command{
	Use:   "dereference",
	Short: "Dereferences a schema, resolving all local $ref's",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE:  runSchemaDereference,
}

func init() {
	rootCmd.AddCommand(schemaCmd)
	schemaCmd.AddCommand(schemaValidateCmd)
	schemaCmd.AddCommand(schemaDereferenceCmd)

	schemaValidateCmd.Flags().StringP("document", "d", "document.json", "Path to JSON document")
	schemaValidateCmd.Flags().StringP("schema", "s", "./schema.json", "Path to JSON Schema")

	schemaDereferenceCmd.Flags().StringP("out", "o", "", "File to output derefenced schema to (default is stdout)")
}

func runSchemaValidate(cmd *cobra.Command, args []string) error {
	schema, _ := cmd.Flags().GetString("schema")
	document, _ := cmd.Flags().GetString("document")
	_, err := jsonschema.Validate(schema, document)
	return err
}

func runSchemaDereference(cmd *cobra.Command, args []string) error {
	schema := args[0]
	out, _ := cmd.Flags().GetString("out")

	var outFile *os.File

	if out == "" || out == "-" {
		outFile = os.Stdout
	} else {
		var err error
		outFile, err = os.OpenFile(out, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer outFile.Close()
	}

	return jsonschema.WriteDereferencedSchema(schema, outFile, nil)
}
