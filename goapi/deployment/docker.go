package deployment

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"text/template"
	"time"

	"github.com/CyberLabs-Infosec/isolet/goapi/config"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

func GetClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func BuildImage(instance_name string, flag string, level int, password string) error {
	dockerClient, err := GetClient()
	if err != nil {
		log.Println(err)
		return err
	}

	imageIDs, isThere, err := ImageExists(instance_name)
	if err != nil {
		return err
	}
	if isThere {
		err = DeleteImage(imageIDs)
		if err != nil {
			return err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	DockerFile, err := GetParsedDockerfile(level, flag, password, config.DEFAULT_USERNAME)
	if err != nil {
		log.Println(err)
		return err
	}

	opts := types.ImageBuildOptions{
		Context:    DockerFile,
		Tags:       []string{instance_name},
		Dockerfile: "Dockerfile",
		Remove:     true,
	}

	res, err := dockerClient.ImageBuild(ctx, DockerFile, opts)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	err = CatchError(res.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetParsedDockerfile(level int, flag string, password string, username string) (*bytes.Reader, error) {
	var content bytes.Buffer
	temp, err := template.ParseFiles(fmt.Sprintf("challenges/level%d/Dockerfile", level))
	if err != nil {
		return nil, err
	}

	err = temp.Execute(&content, struct {
		Username string
		Password string
		Flag     string
		Wargame  string
	}{
		Username: username,
		Password: password,
		Flag:     flag,
		Wargame:  config.WARGAME_NAME,
	})
	if err != nil {
		return nil, err
	}

	filename := "Dockerfile"
	buf := new(bytes.Buffer)
	tarWriter := tar.NewWriter(buf)
	defer tarWriter.Close()

	tarHeader := &tar.Header{
		Name: filename,
		Size: int64(len(content.Bytes())),
	}

	err = tarWriter.WriteHeader(tarHeader)
	if err != nil {
		return nil, err
	}

	_, err = tarWriter.Write(content.Bytes())
	if err != nil {
		return nil, err
	}

	dockerFileTarReader := bytes.NewReader(buf.Bytes())
	return dockerFileTarReader, nil
}

func DeleteImage(imageIDs []string) error {
	cli, err := GetClient()
	if err != nil {
		return err
	}

	for _, id := range imageIDs {
		_, err := cli.ImageRemove(context.Background(), id, types.ImageRemoveOptions{
			Force:         true,
			PruneChildren: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Taken from https://github.com/Atish03/podwiz
func CatchError(rd io.Reader) error {
	var lastLine string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lastLine = scanner.Text()
	}

	errLine := &ErrorLine{}
	json.Unmarshal([]byte(lastLine), errLine)
	if errLine.Error != "" {
		return errors.New(errLine.Error)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func ImageExists(image string) ([]string, bool, error) {
	var imageIDs []string
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return imageIDs, false, err
	}

	filters := filters.NewArgs()
	filters.Add("reference", image+":latest")

	opts := types.ImageListOptions{
		Filters: filters,
	}

	imgs, err := cli.ImageList(context.Background(), opts)
	if err != nil {
		return imageIDs, false, err
	}

	for _, object := range imgs {
		imageIDs = append(imageIDs, object.ID)
	}

	if len(imageIDs) > 0 {
		return imageIDs, true, nil
	}

	return imageIDs, false, nil
}
