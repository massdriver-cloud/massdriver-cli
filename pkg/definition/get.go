package definition

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
)

func GetDefinition(c *client.MassdriverClient, definitionType string) (map[string]interface{}, error) {
	var definition map[string]interface{}

	endpoint := path.Join("artifact-definitions", definitionType)

	req := client.NewRequest("GET", endpoint, nil)

	ctx := context.TODO()
	resp, err := c.Do(&ctx, req)

	if err != nil {
		return definition, err
	}
	defer resp.Body.Close()

	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return definition, err
	}

	if resp.StatusCode != http.StatusOk {
		fmt.Println(string(respBodyBytes))
		return definition, errors.New("received non-200 response from Massdriver: " + resp.Status + " " + definitionType)
	}

	err = json.Unmarshal(respBodyBytes, &definition)
	if err != nil {
		return definition, err
	}

	return definition, nil
}
