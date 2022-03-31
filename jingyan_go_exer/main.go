package main

import (
	"bytes"
	"context"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	container "github.com/docker/docker/api/types/container"
)

const OnosImageName = "registry.hub.docker.com/onosproject/onos"
const OnosContainerName = "onos"

var onosAppsList = []string{"org.onosproject.fwd", "org.onosproject.openflow"}

const MininetImageName = "registry.hub.docker.com/iwaseyusuke/mininet"
const MininetContainerName = "mininet"
const HostVolMap = "/Users/yanjing/Desktop/jingyan_go_exer/mininet"
const MininetVolMap = "/tmp"

const TemplateTopologyFile = "/Users/yanjing/Desktop/jingyan_go_exer/topo_templates/create_tree_topo_template.sh"
const NewTopologyFile = "/Users/yanjing/Desktop/jingyan_go_exer/mininet/create_tree_topo.sh"

var MininetExecCmd = []string{"/bin/bash", "create_tree_topo.sh"}

func NewClientWithOpts() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("Initializing API client for docker daemon failed as: %v", err)
		os.Exit(1)
	}
	return cli
}

func ImagePull(cli *client.Client, imageName string) {
	_, err := cli.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
	if err != nil {
		log.Printf("Pulling docker image %s failed as: %v", imageName, err)
		os.Exit(1)
	}
}

func ContainerStart(cli *client.Client, config *container.Config, hostConfig *container.HostConfig, containerName string) string {
	resp, err := cli.ContainerCreate(context.Background(), config, hostConfig, nil, nil, containerName)
	if err != nil {
		log.Printf("Creating %s container failed as: %v", containerName, err)
		os.Exit(1)
	}

	if err := cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Printf("Starting %s container failed as %v: ", resp.ID, err)
		os.Exit(1)
	}

	return resp.ID
}

func ActivateOnosApp(appName string) {
	url := "http://127.0.0.1:8181/onos/v1/applications/" + appName + "/active"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Printf("HTTP Request establishment failed as: %v", err)
		os.Exit(1)
	}
	req.SetBasicAuth("onos", "rocks")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Sending HTTP Request failed as %v: ", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	httpStatusCode := resp.StatusCode
	if httpStatusCode != 200 {
		log.Printf("Activate the %s failed because of %d HTTP StatusCode", appName, httpStatusCode)
		os.Exit(1)
	}
}

func ContainerInspect(cli *client.Client, containerID string) string {
	var OnosIP = ""
	cli.ContainerInspect(context.Background(), containerID)

	onosInspect, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		log.Printf("Fialed to get the infomation of ONOS container as: %v", err)
		os.Exit(1)
	}

	for _, network := range onosInspect.NetworkSettings.Networks {
		OnosIP = network.IPAddress
	}

	return OnosIP
}

func GenerateTopologyCreationFile(OnosIP string) {
	input, err := ioutil.ReadFile(TemplateTopologyFile)
	if err != nil {
		log.Printf("Reading the template file failed as: %v", err)
		os.Exit(1)
	}
	output := bytes.Replace(input, []byte("onos_ip"), []byte(OnosIP), -1)

	if err = ioutil.WriteFile(NewTopologyFile, output, 0666); err != nil {
		log.Printf("Writing to the new file failed as: %v", err)
		os.Exit(1)
	}
}

func ContainerExecCmd(cli *client.Client, containerID string, mininetExecConfig types.ExecConfig) string {
	headers, err := cli.ContainerExecCreate(context.Background(), containerID, mininetExecConfig)
	if err != nil {
		log.Printf("Failed to create the exec configuration to run the exec process as: %v", err)
		os.Exit(1)
	}

	resp, err := cli.ContainerExecAttach(context.Background(), headers.ID, types.ExecStartCheck{})
	if err != nil {
		log.Printf("Failed to Attach the connection to the exec process in the server as: %v", err)
		os.Exit(1)
	}
	defer resp.Close()

	data, _ := ioutil.ReadAll(resp.Reader)
	return string(data)
}

func main() {
	cli := NewClientWithOpts()

	imageList := []string{MininetImageName, OnosImageName}
	for _, imageName := range imageList {
		ImagePull(cli, imageName)
	}

	onosConfig, onosHostConfig := &container.Config{
		Image: OnosImageName,
		ExposedPorts: nat.PortSet{
			"8181/tcp": struct{}{},
			"8101/tcp": struct{}{},
			"5005/tcp": struct{}{},
			"830/tcp":  struct{}{}},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			"8181/tcp": []nat.PortBinding{{HostPort: "8181"}},
			"8101/tcp": []nat.PortBinding{{HostPort: "8101"}},
			"5005/tcp": []nat.PortBinding{{HostPort: "5005"}},
			"830/tcp":  []nat.PortBinding{{HostPort: "830"}}},
	}

	var onosID = ContainerStart(cli, onosConfig, onosHostConfig, OnosContainerName)

	mininetConfig, mininetHostConfig := &container.Config{
		Tty:   true,
		Image: MininetImageName,
	}, &container.HostConfig{
		Privileged: true,
		Mounts: []mount.Mount{{
			Type:   mount.TypeBind,
			Source: HostVolMap,
			Target: MininetVolMap,
		}},
	}

	var mininetID = ContainerStart(cli, mininetConfig, mininetHostConfig, MininetContainerName)

	time.Sleep(30 * time.Second)

	for _, appName := range onosAppsList {
		ActivateOnosApp(appName)
	}

	var onosIP = ContainerInspect(cli, onosID)
	GenerateTopologyCreationFile(onosIP)

	mininetExecConfig := types.ExecConfig{
		Privileged:   true,
		Tty:          true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          MininetExecCmd,
	}

	ContainerExecCmd(cli, mininetID, mininetExecConfig)
}
