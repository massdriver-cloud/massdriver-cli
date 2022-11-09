package cmd

import (
	"github.com/massdriver-cloud/massdriver-cli/pkg/bundle"
	"github.com/massdriver-cloud/massdriver-cli/pkg/common"
	"github.com/massdriver-cloud/massdriver-cli/pkg/containers"

	"github.com/spf13/cobra"
)

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Container development tools",
	Long:  ``,
}

var imageBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "",

	RunE: runImageBuild,
}

var imageListCmd = &cobra.Command{
	Use:   "list",
	Short: "",

	RunE: runImageList,
}

var imagePushCmd = &cobra.Command{
	Use:   "push",
	Short: "",

	RunE: runImagePush,
}

var packageImgCmd = &cobra.Command{
	Use:   "package",
	Short: "",

	RunE: runPackageCmd,
}

func init() {
	rootCmd.AddCommand(imageCmd)

	imageCmd.AddCommand(imageBuildCmd)
	imageCmd.AddCommand(imageListCmd)
	imageCmd.AddCommand(imagePushCmd)
	imageCmd.AddCommand(packageImgCmd)
}

func runImageBuild(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	_, errClient := initClient(cmd)
	if errClient != nil {
		return errClient
	}

	errBuild := containers.BuildImage(containers.BuildOptions{})
	if errBuild != nil {
		return errBuild
	}

	return nil
}

func runImageList(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	_, errClient := initClient(cmd)
	if errClient != nil {
		return errClient
	}

	errBuild := containers.ListImages()
	if errBuild != nil {
		return errBuild
	}

	return nil
}

func runImagePush(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	_, errClient := initClient(cmd)
	if errClient != nil {
		return errClient
	}

	errBuild := containers.PushImage("latest")
	if errBuild != nil {
		return errBuild
	}
	return nil
}

func runPackageCmd(cmd *cobra.Command, args []string) error {
	setupLogging(cmd)

	_, errClient := initClient(cmd)
	if errClient != nil {
		return errClient
	}
	b, err := bundle.Parse(common.MassdriverYamlFilename, nil)
	if err != nil {
		return err
	}

	errBuild := containers.Package(b)
	if errBuild != nil {
		return errBuild
	}

	return nil
}
