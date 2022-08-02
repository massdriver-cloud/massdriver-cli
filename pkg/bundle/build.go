package bundle

import (
	"fmt"

	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
	"github.com/massdriver-cloud/massdriver-cli/pkg/provisioners/terraform"

	"github.com/rs/zerolog/log"
)

const configFile = "massdriver.yaml"

func (b *Bundle) Build(c *client.MassdriverClient, output string) error {
	// b, err := Parse(configFile, nil)
	// if err != nil {
	// 	log.Error().Err(err).Msg("an error occurred while parsing bundle")
	// 	return err
	// }

	err := b.Hydrate(configFile, c)
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while hydrating bundle")
		return err
	}

	err = b.GenerateSchemas(output)
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while generating bundle schema files")
		return err
	}

	steps := b.Steps
	if b.Steps == nil {
		steps = []Step{
			{
				Path:        "src",
				Provisioner: "terraform",
			},
		}
	}

	// TODO: merge w/ app package
	for _, step := range steps {
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
