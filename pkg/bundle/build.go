package bundle

import (
	"fmt"

	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
	"github.com/massdriver-cloud/massdriver-cli/pkg/provisioners/terraform"

	"github.com/rs/zerolog/log"
)

const configFile = "massdriver.yaml"

func (bundle *Bundle) Build(c *client.MassdriverClient, output string) error {
	// b, err := Parse(configFile, nil)
	// if err != nil {
	// 	log.Error().Err(err).Msg("an error occurred while parsing bundle")
	// 	return err
	// }

	err := bundle.Hydrate(configFile, c)
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while hydrating bundle")
		return err
	}

	err = bundle.GenerateSchemas(output)
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while generating bundle schema files")
		return err
	}

	for _, step := range bundle.Steps {
		switch step.Provisioner {
		case "terraform":
			err = terraform.GenerateFiles(output, step.Path)
			if err != nil {
				log.Error().Err(err).Str("provisioner", step.Provisioner).Msg("an error occurred while generating provisioner files")
				return err
			}
		default:
			log.Error().Msg("unknown provisioner: " + step.Provisioner)
			return fmt.Errorf("unknown provisioner: %v", step.Provisioner)
		}
	}

	return nil
}
