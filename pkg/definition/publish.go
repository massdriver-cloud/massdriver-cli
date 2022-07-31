package definition

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/massdriver-cloud/massdriver-cli/pkg/client"
	"github.com/rs/zerolog/log"
)

func (art *Definition) Publish(c *client.MassdriverClient) error {
	bodyBytes, err := json.Marshal(*art)
	if err != nil {
		return err
	}

	req := client.NewRequest("PUT", "artifact-definitions", bytes.NewBuffer(bodyBytes))
	ctx := context.TODO()
	resp, err := c.Do(&ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		respBodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			return err2
		}
		log.Debug().Msg(string(respBodyBytes))
		return errors.New("received non-200 response from Massdriver: " + resp.Status)
	}

	return nil
}
