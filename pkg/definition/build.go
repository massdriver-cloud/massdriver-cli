package definition

import (
	b64 "encoding/base64"

	"github.com/rs/zerolog/log"
)

func Build(filePath string) error {
	def, err := Parse(filePath)
	if err != nil {
		return err
	}
	log.Info().Msgf("Building definition %v", def)
	if errAdd := addProvisioners(def); errAdd != nil {
		return errAdd
	}
	return nil
}

func addProvisioners(def *Definition) error {
	// for each provisioner,
	base64Provisioner, err := buildProvisioner("")
	if err != nil {
		return err
	}
	(*def)["$md.provisioners.terraform"] = base64Provisioner
	return nil
}

// for each provisioner in
func buildProvisioner(filePath string) (string, error) {
	data := "abc123!?$*&()'-=@~"
	sEnc := b64.StdEncoding.EncodeToString([]byte(data))
	return sEnc, nil
}
