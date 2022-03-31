package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"main/activate_onos_app"
	"main/generate_file4topo_create"
)

func main() {
	//
	// Initialize a new API client to communicate with the docker daemon
	// Print err msg and exit if the operation fails
	//
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Print("The initialization of API client for docker daemon failed as the following reason: ")
		panic(err)
	}

	//
	// Pull "mininet" and "onos" images from docker hub
	// Print err msg and exit if the operation fails
	//
	image_list := []string{"registry.hub.docker.com/iwaseyusuke/mininet", "registry.hub.docker.com/onosproject/onos"}
	for _, image := range image_list {
		_, err := cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
		if err != nil {
			log.Print("Pull docker image" + image + "failed as the following reason: ")
			panic(err)
		}
	}

	//
	// Create a container named "onos" based on the given "&container.Config" and "&container.HostConfig" configuration
	// The "&container.Config" and "&container.HostConfig" are mainly used for port mapping here
	// Print err msg and exit if the operation fails
	//
	onos_resp, err := cli.ContainerCreate(context.Background(),
		&container.Config{
			Image: "registry.hub.docker.com/onosproject/onos",
			ExposedPorts: nat.PortSet{
				"8181/tcp": struct{}{},
				"8101/tcp": struct{}{},
				"5005/tcp": struct{}{},
				"830/tcp":  struct{}{}},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				"8181/tcp": []nat.PortBinding{{HostPort: "8181"}},
				"8101/tcp": []nat.PortBinding{{HostPort: "8101"}},
				"5005/tcp": []nat.PortBinding{{HostPort: "5005"}},
				"830/tcp":  []nat.PortBinding{{HostPort: "830"}}},
		}, nil, nil, "onos")
	if err != nil {
		log.Print("Creating ONOS contianer failed as the following reason: ")
		panic(err)
	}

	//
	// Start the "onos" container based on the "onos_resp.ID" created above
	// Print err msg and exit if the operation fails
	//
	if err := cli.ContainerStart(context.Background(), onos_resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Print("Starting ONOS container failed as the following reason: ")
		panic(err)
	}

	//
	// Create a container named "mininet" based on the given "&container.HostConfig" configuration
	// The "&container.HostConfig" is mainly used for volumn mapping
	// Print err msg and exit if the operation fails
	//
	mininet_resp, err := cli.ContainerCreate(context.Background(),
		&container.Config{
			Tty:   true,
			Image: "registry.hub.docker.com/iwaseyusuke/mininet",
		},
		&container.HostConfig{
			Privileged: true,
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: "/Users/yanjing/Desktop/jingyan_go_exer/mininet",
					Target: "/tmp",
				},
			},
		}, nil, nil, "mininet")
	if err != nil {
		log.Print("Creating mininet contianer failed as the following reason: ")
		panic(err)
	}

	//
	// Start the "mininet" container based on the "mininet_resp.ID" created above
	// Print err msg and exit if the operation fails
	//
	if err := cli.ContainerStart(context.Background(), mininet_resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Print("Starting mininet container failed as the following reason: ")
		panic(err)
	}

	//
	// Wait for 30s until ONOS container is ready for accepting HTTP request
	//
	time.Sleep(30 * time.Second)

	//
	// Call "func Activate_onos_app(app_name string) int" to activate necessary applications for ONOS container
	// Print err msg and exit if the operation fails
	//
	onos_apps_list := []string{"org.onosproject.fwd", "org.onosproject.openflow"}
	for _, onos_app := range onos_apps_list {
		resp_StatusCode := activate_onos_app.Activate_onos_app(onos_app)
		if resp_StatusCode != 200 {
			panic("Activate the " + onos_app + " failed because of " + string(rune(resp_StatusCode)) + " HTTP StatusCode")
		}
	}

	//
	// Obtain the IP address of ONOS container and store the value in the Onos_ip variable
	//
	var Onos_ip = ""

	onos_inspect, err := cli.ContainerInspect(context.Background(), onos_resp.ID)
	if err != nil {
		log.Print("Fialed to get the infomation of ONOS container as the following reason: ")
		panic(err)
	}

	for _, network := range onos_inspect.NetworkSettings.Networks {
		Onos_ip = network.IPAddress
	}

	generate_file4topo_create.Generate_file4topo_create(Onos_ip)

	// Execute the "/bin/bash /tmp/create_tree_topo.sh" in the mininet container

	//
	// Create a exec configuration to run an exec process which is specified in the "types.ExecConfig"
	//
	idre, err := cli.ContainerExecCreate(context.Background(), mininet_resp.ID, types.ExecConfig{
		Privileged:   true,
		Tty:          true,
		AttachStderr: true,
		AttachStdout: true,
		WorkingDir:   "/tmp",
		Cmd:          []string{"/bin/bash", "create_tree_topo.sh"},
	})
	if err != nil {
		log.Print("Failed to create the exec configuration to run the exec process as the following reason: ")
		panic(err)
	}

	//
	// Attach the connection to the exec process in the server
	//
	hijare, err := cli.ContainerExecAttach(context.Background(), idre.ID, types.ExecStartCheck{})
	if err != nil {
		log.Print("Failed to Attach the connection to the exec process in the server as the following reason: ")
		panic(err)
	}
	defer hijare.Close()

	//
	// Output of the command "/bin/bash /tmp/create_tree_topo.sh"
	//
	data, _ := ioutil.ReadAll(hijare.Reader)
	fmt.Println(string(data))
}
