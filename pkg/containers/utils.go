package containers

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
)

type LogLine struct {
	Stream string `json:"stream"`
}

type LogLinePush struct {
	Status   string `json:"status"`
	Progress string `json:"progress"`
}

// TODO: move hackathon notes
// AWS
// https://docs.aws.amazon.com/AmazonECR/latest/userguide/registry_auth.html
// (Get-ECRLoginCommand).Password | docker login --username AWS --password-stdin aws_account_id.dkr.ecr.region.amazonaws.com
// TOKEN=$(aws ecr get-authorization-token --output text --query 'authorizationData[].authorizationToken')
// curl -i -H "Authorization: Basic $TOKEN" https://aws_account_id.dkr.ecr.region.amazonaws.com/v2/amazonlinux/tags/list
// Azure
// https://learn.microsoft.com/en-us/azure/container-registry/container-registry-authentication?tabs=azure-cli
// GCP
// Pushing an Image
// make install.macos && \
//   REGISTRY_AUTH_TOKEN=$(gcloud auth print-access-token --impersonate-service-account massdriver-sa@md-wbeebe-0808-example-apps.iam.gserviceaccount.com) && \
//   mass image push
func getAuthConfig(authToken string, image string) (string, error) {
	authConfig := types.AuthConfig{
		Username: registryUsernameFromImage(image),
		Password: authToken,
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return "", fmt.Errorf("error when encoding authConfig. err: %s", err)
	}
	return base64.URLEncoding.EncodeToString(encodedJSON), nil
}

func registryUsernameFromImage(image string) string {
	if strings.Contains(image, "dkr.ecr") {
		return "AWS"
	}
	if strings.Contains(image, "azurecr.io") {
		return "00000000-0000-0000-0000-000000000000"
	}
	// GCP
	return "oauth2accesstoken"
}

func print(rd io.Reader) error {
	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		logStr := strings.TrimSuffix(scanner.Text(), "\n")
		log := []byte(logStr)

		var logLine LogLinePush
		err := json.Unmarshal(log, &logLine)
		if err != nil {
			return err
		}
		msg := strings.TrimSuffix(logLine.Status, "\n")
		if msg == "" {
			continue
		}
		fmt.Println(msg)
	}

	// errLine := &ErrorLine{}
	// json.Unmarshal([]byte(lastLine), errLine)
	// if errLine.Error != "" {
	// 	// return errors.New(errLine.Error)
	// 	return nil
	// }

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func humanFileSize(size float64) string {
	var suffixes = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

	base := math.Log(size) / math.Log(1024)
	getSize := Round(math.Pow(1024, base-math.Floor(base)), .5, 2)
	getSuffix := suffixes[int(math.Floor(base))]
	return strconv.FormatFloat(getSize, 'f', -1, 64) + "" + string(getSuffix)
}

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
