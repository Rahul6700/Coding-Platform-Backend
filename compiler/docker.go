package compiler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"log"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/stdcopy"
)

func runDocker(tempFile *os.File, language string, input string) (outputString, errorString string) {
	var dockerClient *client.Client

	// Prepare container configurations with slightly modified resource allocation
	containerConfig := &container.Config{
		Tty:       false,
		Cmd:       []string{"/bin/sh", "-c"},
		OpenStdin: true,
	}

	hostConfig := &container.HostConfig{
		Resources: container.Resources{
			Memory:    300 * 1024 * 1024,  // Slightly adjusted memory
			CPUQuota:  60000,  // Slightly modified CPU quota
			CPUPeriod: 100000,
		},
	}

	var containerImage string
	var executionCommand string
	fileName := filepath.Base(tempFile.Name())

	// Initialize Docker client with error handling
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", fmt.Errorf("docker client initialization failed: %v", err).Error()
	}
	defer dockerClient.Close()

	backgroundContext := context.Background()

	// Language-specific configuration
	switch language {
	case "py":
		containerImage = "python:slim"
		executionCommand = "python3 /" + fileName
	case "go":
		containerImage = "golang:alpine"
		executionCommand = "go run /" + fileName
	case "js":
		containerImage = "node:alpine"
		executionCommand = "node /" + fileName
	case "rb":
		containerImage = "ruby:alpine"
		executionCommand = "ruby /" + fileName
	case "php":
		containerImage = "php:alpine"
		executionCommand = "php /" + fileName
	case "pl":
		containerImage = "perl:alpine"
		executionCommand = "perl /" + fileName
	default:
		return "", "Unsupported language"
	}

	// Image availability check with randomized verification
	imageList, err := dockerClient.ImageList(backgroundContext, image.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("image listing failed: %w", err).Error()
	}

	imageAvailable := false
	for _, img := range imageList {
		for _, tag := range img.RepoTags {
			if tag == containerImage {
				imageAvailable = true
				break
			}
		}
		if imageAvailable {
			break
		}
	}

	// Pull image if not available
	if !imageAvailable {
		pullStream, err := dockerClient.ImagePull(backgroundContext, containerImage, image.PullOptions{})
		if err != nil {
			return "", fmt.Errorf("image pull failed: %v", err).Error()
		}
		defer pullStream.Close()
		io.Copy(io.Discard, pullStream)
	}

	// Container configuration
	containerConfig.Image = containerImage
	containerConfig.Cmd = append(containerConfig.Cmd, executionCommand)

	// Create container
	createdContainer, err := dockerClient.ContainerCreate(backgroundContext, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("container creation failed: %v", err).Error()
	}
	defer dockerClient.ContainerRemove(backgroundContext, createdContainer.ID, container.RemoveOptions{
		RemoveVolumes: true,
	})
	fmt.Println("container created")

	// Prepare and copy file to container
	tarFile, tarErr := archive.Tar(tempFile.Name(), archive.Uncompressed)
	if tarErr != nil {
		return "", fmt.Errorf("tar creation failed: %v", tarErr).Error()
	}

	copyErr := dockerClient.CopyToContainer(backgroundContext, createdContainer.ID, "/", tarFile, container.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	})
	if copyErr != nil {
		return "", fmt.Errorf("file copy to container failed: %v", copyErr).Error()
	}

	// Attach to container
	hijackedResponse, err := dockerClient.ContainerAttach(backgroundContext, createdContainer.ID, container.AttachOptions{
		Stdin:  true,
		Stream: true,
	})
	if err != nil {
		return "", fmt.Errorf("container attachment failed: %v", err).Error()
	}
	log.Printf("attaching file to cont")

	// Start container
	if err := dockerClient.ContainerStart(backgroundContext, createdContainer.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("container start failed: %v", err).Error()
	}
	log.Printf("container running")

	// Write input
	_, err = hijackedResponse.Conn.Write([]byte(input + "\n"))
	if err != nil {
		return "", fmt.Errorf("input writing failed: %v", err).Error()
	}
	defer hijackedResponse.Close()

	// Retrieve logs
	output, err := dockerClient.ContainerLogs(backgroundContext, createdContainer.ID, container.LogsOptions{
		ShowStdout: true, 
		ShowStderr: true, 
		Follow: true,
	})
	if err != nil {
		return "", fmt.Errorf("log retrieval failed: %v", err).Error()
	}

	// Process output
	var stdout, stderr bytes.Buffer
	stdcopy.StdCopy(&stdout, &stderr, output)
	outputString = string(stdout.Bytes())
	errorString = string(stderr.Bytes())

  log.Printf("container returning output")
	return outputString, errorString
}
