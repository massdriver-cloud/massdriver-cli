package containers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/rs/zerolog/log"
)

type LogLine struct {
	Stream string `json:"stream"`
}

func BuildImage() error {
	log.Info().Msg("Building image")
	ctx := context.TODO()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	opts := types.ImageBuildOptions{
		Dockerfile: "./Dockerfile",
		// TODO: auto-tag with md_name_prefix and current environment
		Tags: []string{"latest"},
		// Remove: true,
	}
	tar, err := archive.TarWithOptions(".", &archive.TarOptions{})
	if err != nil {
		return err
	}

	res, errBuild := cli.ImageBuild(ctx, tar, opts)
	if errBuild != nil {
		return errBuild
	}

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		logStr := strings.TrimSuffix(scanner.Text(), "\n")
		log := []byte(logStr)

		var logLine LogLine
		err := json.Unmarshal(log, &logLine)
		if err != nil {
			return err
		}
		msg := strings.TrimSuffix(logLine.Stream, "\n")
		if msg == "" {
			continue
		}
		fmt.Println(msg)
	}

	defer res.Body.Close()

	return nil
}

func ListImages() error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.DiscardEmptyColumns)
	fmt.Fprintln(w, "REPOSITORY\tTAG\tIMAGE ID\tCREATED\tSIZE")

	// log.Info().Msg("REPOSITORY\tTAG\tIMAGE ID\tCREATED\tSIZE")
	for _, image := range images {
		if len(image.RepoDigests) == 0 || len(image.RepoTags) == 0 {
			continue
		}

		repo := strings.Split(image.RepoDigests[0], "@")[0]
		mostRecentTag := strings.Split(image.RepoTags[0], ":")[1]
		if repo == "<none>" {
			continue
		}
		fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%s\t%s", repo, mostRecentTag, image.ID[:10], "Created", humanFileSize(float64(image.VirtualSize))))

		// log.Info().Msgf("%s\t%s\t%s\t%s\t%s", repo, mostRecentTag, "ID", "Created", humanFileSize(float64(image.Size)))
	}
	w.Flush()
	return nil
}

const repoHost = "https://hub.docker.com"

func PushImage(imageURI string) error {
	ctx := context.TODO()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	reader, errPush := cli.ImagePush(ctx, imageURI, types.ImagePushOptions{
		// All           bool
		// RegistryAuth  string // RegistryAuth is the base64 encoded credentials for the registry
		// PrivilegeFunc RequestPrivilegeFunc
		// Platform      string
	})
	if errPush != nil {
		return errPush
	}
	defer reader.Close()
	return print(reader)
}

func print(rd io.Reader) error {
	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		logStr := strings.TrimSuffix(scanner.Text(), "\n")
		log := []byte(logStr)
		fmt.Println(log)
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
