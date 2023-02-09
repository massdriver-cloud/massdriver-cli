package image

import "github.com/massdriver-cloud/massdriver-cli/pkg/api2"

func Build(input PushImageInput, imageClient Client) error {
	containerRepository := &api2.ContainerRepository{
		RepositoryUri: "hank/was",
	}
	res, err := imageClient.BuildImage(input, containerRepository)

	if err != nil {
		return err
	}

	err = handleResponseBuffer(res.Body)

	if err != nil {
		return err
	}

	return nil
}
