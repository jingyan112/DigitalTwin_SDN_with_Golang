package generate_file4topo_create

import (
	"bytes"
	"io/ioutil"
	"log"
)

//
// Function name: Generate_file4topo_create
//
// Function description: Read the template file of creating network topology,
//                       and replace the "onos_ip" in template file with the ONOS container IP address
//                       then write all the info to a new file
//
// Paramters: Onos_ip (string), the IP address of ONOS container
//
// Returns:
//
func Generate_file4topo_create(Onos_ip string) {
	input, err := ioutil.ReadFile("/Users/yanjing/Desktop/jingyan_go_exer/topo_templates/create_tree_topo_template.sh")
	if err != nil {
		log.Print("Reading the template file failed as the following reason: ")
		panic(err)
	}
	output := bytes.Replace(input, []byte("onos_ip"), []byte(Onos_ip), -1)

	if err = ioutil.WriteFile("/Users/yanjing/Desktop/jingyan_go_exer/mininet/create_tree_topo.sh", output, 0666); err != nil {
		log.Print("Writing to the new file failed as the following reason: ")
		panic(err)
	}
}
