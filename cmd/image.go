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
	Use:   "push <namespace>/<image-name>",
	Short: "Push an image to ECR, ACR or GAR",
	RunE:  runImagePush,
	Args:  cobra.ExactArgs(1),
}

var dockerBuildContext string
var dockerfileName string
var targetPlatform string
var buildOnly bool
var tag string
var artifactId string
var region string

func init() {
	rootCmd.AddCommand(imageCmd)
	imageCmd.AddCommand(imagePushCmd)
	imagePushCmd.Flags().StringVarP(&dockerBuildContext, "build-context", "b", ".", "Path to the directory to build the image from")
	imagePushCmd.Flags().StringVarP(&dockerfileName, "dockerfile", "f", "Dockerfile", "Name of the dockerfile to build from if you have named it anything other than Dockerfile")
	imagePushCmd.Flags().StringVarP(&tag, "image-tag", "t", "latest", "Unique identifier for this version of the image")
	imagePushCmd.Flags().StringVarP(&artifactId, "artifact", "a", "", "Massdriver ID of the artifact used to create the repository and generate repository credentials.")
	imagePushCmd.MarkFlagRequired("artifact")
	imagePushCmd.Flags().StringVarP(&region, "region", "r", "", "Cloud region to push the image to")
	imagePushCmd.MarkFlagRequired("region")
	imagePushCmd.Flags().StringVarP(&targetPlatform, "platform", "p", "linux/amd64", "")
	imagePushCmd.Flags().BoolVarP(&buildOnly, "build", "o", false, "")
}

func runImagePush(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	orgID := os.Getenv("MASSDRIVER_ORG_ID")
	if orgID == "" {
		log.Fatal().Msg("MASSDRIVER_ORG_ID must be set")
	}

	APIKey := os.Getenv("MASSDRIVER_API_KEY")
	if APIKey == "" {
		log.Fatal().Msg("MASSDRIVER_API_KEY must be set")
	}

	pushInput := image.PushImageInput{
		OrganizationId: orgID,
		ImageName:      args[0],
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
	if buildOnly {
		return image.Build(pushInput, imageClient)
	}

	return image.Push(gqlclient, pushInput, imageClient)
}

func validatePushInputAndAddFlags(input *image.PushImageInput, cmd *cobra.Command) error {
	flagsToSet := []PushFlag{
		{Flag: "dockerfile", Attribute: &input.Dockerfile},
		{Flag: "build-context", Attribute: &input.DockerBuildContext},
		{Flag: "platform", Attribute: &input.TargetPlatform},
		{Flag: "image-tag", Attribute: &input.Tag},
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
