package definition

import (
	b64 "encoding/base64"
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"
)

func Build(filePath string) (*DefinitionFile, error) {
	log.Info().Msgf("Building %v", filePath)
	def, err := Parse(filePath)
	if err != nil {
		return def, err
	}
	log.Info().Msgf("Def is %+v", def)
	log.Info().Msgf("Name is %v", def.Md.Name)
	if errAdd := addProvisioners(def); errAdd != nil {
		return def, errAdd
	}
	log.Info().Msgf("Definition built %+v", def)
	return def, nil
}

// walk provisioners, base64 various files
func addProvisioners(def *DefinitionFile) error {
	// defName := (*def)["$md"].(map[string]interface{})["name"].(string)
	defName := def.Md.Name
	terrformProvisionerPath := defName + "/provisioners/terraform"

	fileInfo, err := os.Stat(terrformProvisionerPath)
	if os.IsNotExist(err) || !fileInfo.IsDir() {
		log.Debug().Msgf("No provisioners found for %v", defName)
		return nil
	}

	// loop over all files in provisioner dir?
	base64File, err := base64File(terrformProvisionerPath + "/_provider.tf")
	if err != nil {
		return err
	}
	def.Md.Provisioners.Terraform = base64File
	return nil
}

func base64File(path string) (string, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	encoded := b64.StdEncoding.EncodeToString(bytes)
	return encoded, nil
}
