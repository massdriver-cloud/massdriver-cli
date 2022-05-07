package definition

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
)

func (art *Definition) Publish(c *client.MassdriverClient) error {

	bodyBytes, err := json.Marshal(*art)
	if err != nil {
		return err
	}

	req := client.NewRequest("PUT", "artifact-definitions", bytes.NewBuffer(bodyBytes))
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Println(string(respBodyBytes))
		return errors.New("received non-200 response from Massdriver: " + resp.Status)
	}

	return nil
}
