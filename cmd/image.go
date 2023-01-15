package cmd

import (
	"os"
	"strings"

	"github.com/massdriver-cloud/massdriver-cli/pkg/api2"
	"github.com/massdriver-cloud/massdriver-cli/pkg/image"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type PushFlag struct {
	Flag      string
	Attribute *string
}

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Container image integration Massdriver",
	Long:  ``,
}

var imagePushCmd = &cobra.Command{
	Use:   "push <namespace>/<image-name> <cloud-region> <massdriver-artifactId>",
	Short: "Push an image to ECR, ACR or GAR",
	RunE:  runImagePush,
}

var DockerBuildContext string
var DockerfileName string
var Tag string
var ArtifactId string
var ImageName string
var Region string

func init() {
	rootCmd.AddCommand(imageCmd)
	imageCmd.AddCommand(imagePushCmd)
	imagePushCmd.Flags().StringVarP(&DockerBuildContext, "build-context", "b", ".", "Path to the directory to build the image from")
	imagePushCmd.Flags().StringVarP(&DockerfileName, "dockerfile", "f", "Dockerfile", "Name of the dockerfile to build from if you have named it anything other than Dockerfile")
	imagePushCmd.Flags().StringVarP(&Tag, "image-tag", "t", "latest", "Unique identifier for this version of the image")
	imagePushCmd.Flags().StringVarP(&ImageName, "image-name", "n", "", "Name of the image to push in name spaced form I.E. acme-corp/mail-service")
	imagePushCmd.MarkFlagRequired("image-name")
	imagePushCmd.Flags().StringVarP(&ArtifactId, "artifact", "a", "", "Massdriver ID of the artifact used to create the repository and generate repository credentials.")
	imagePushCmd.MarkFlagRequired("artifact")
	imagePushCmd.Flags().StringVarP(&Region, "region", "r", "", "Cloud region to push the image to")
	imagePushCmd.MarkFlagRequired("region")
}

func runImagePush(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	orgID := os.Getenv("MASSDRIVER_ORG_ID")
	if orgID == "" {
		log.Fatal().Msg("MASSDRIVER_ORG_ID must be set")
	}

	APIKey := os.Getenv("MASSDRIVER_API_KEY")
	if orgID == "" {
		log.Fatal().Msg("MASSDRIVER_API_KEY must be set")
	}

	pushInput := image.PushImageInput{
		OrganizationId: orgID,
	}

	err := validatePushInputAndAddFlags(&pushInput, cmd)

	if err != nil {
		return err
	}

	gqlclient := api2.NewClient(APIKey)
	imageClient, err := image.NewImageClient()

	if err != nil {
		return err
	}

	return image.Push(gqlclient, pushInput, imageClient)
}

func validatePushInputAndAddFlags(input *image.PushImageInput, cmd *cobra.Command) error {
	flagsToSet := []PushFlag{
		{Flag: "dockerfile", Attribute: &input.Dockerfile},
		{Flag: "build-context", Attribute: &input.DockerBuildContext},
		{Flag: "image-tag", Attribute: &input.Tag},
		{Flag: "image-name", Attribute: &input.ImageName},
		{Flag: "artifact", Attribute: &input.ArtifactId},
		{Flag: "region", Attribute: &input.Location},
	}

	for _, flag := range flagsToSet {
		value, err := cmd.Flags().GetString(flag.Flag)

		if err != nil {
			return err
		}

		*flag.Attribute = value
	}

	if invalidImageName(input.ImageName) {
		log.Fatal().Msg("massdriver enforces the practice of namespacing images. please enter an image name in the format of namespace/image-name")
	}

	return nil
}

func invalidImageName(imageName string) bool {
	return !(len(strings.Split(imageName, "/")) == 2)
}
