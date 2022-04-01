# DigitalTwin_SDN_with_Golang
Re-write the previous digital twin with SDN project with Golang

Golang env setup
- go mod init main
- go mod tidy
- go run main.go

What does the code implement?
- Pull "mininet" and "onos" docker image from dockerhub
- Create "mininet" and "onos" containers and start them with configuring port mapping and volumn mapping for containers
- Activate two applications for "onos" container through RESTAPI
- Obtain the IP address of onos container 
- Replace the "onos_ip" string in the template file with the IP address of onos container
- Execute the template bash file in mininet container

Output of the code
- Check "docker image ls"
- Check "docker container ps -a"
- Check "http://127.0.0.1:8181/onos/ui/login.html" in your browser with "onos:rocks" as "username/password", you can see something cool
