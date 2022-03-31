# 1. Create namespaces
ip netns add h1
ip netns add h2
ip netns add h3
ip netns add h4
ip netns add h5
ip netns add h6

# 2. Create openvswitch
ovs-vsctl add-br s1
ovs-vsctl set bridge s1 protocols=OpenFlow14
ovs-vsctl add-br s2
ovs-vsctl set bridge s2 protocols=OpenFlow14
ovs-vsctl add-br s3
ovs-vsctl set bridge s3 protocols=OpenFlow14
ovs-vsctl add-br s4
ovs-vsctl set bridge s4 protocols=OpenFlow14

# 3.1 Create vethernet links among switches and hosts
ip link add h1-eth1 type veth peer name s2-eth1
ip link add h2-eth1 type veth peer name s3-eth1
ip link add h3-eth1 type veth peer name s3-eth2
ip link add h4-eth1 type veth peer name s4-eth1
ip link add h5-eth1 type veth peer name s1-eth4
ip link add h6-eth1 type veth peer name s1-eth5

# 3.2 Create vethernet links among switches and switches
ip link add s1-eth1 type veth peer name s2-eth2
ip link add s1-eth2 type veth peer name s3-eth3
ip link add s1-eth3 type veth peer name s4-eth2

# 4. Move host ports into namespaces
ip link set h1-eth1 netns h1
ip link set h2-eth1 netns h2
ip link set h3-eth1 netns h3
ip link set h4-eth1 netns h4
ip link set h5-eth1 netns h5
ip link set h6-eth1 netns h6

# 5. Connect switch ports to OVS
ovs-vsctl add-port s1 s1-eth1
ovs-vsctl add-port s1 s1-eth2
ovs-vsctl add-port s1 s1-eth3
ovs-vsctl add-port s1 s1-eth4
ovs-vsctl add-port s1 s1-eth5
ovs-vsctl add-port s2 s2-eth1
ovs-vsctl add-port s2 s2-eth2
ovs-vsctl add-port s3 s3-eth1
ovs-vsctl add-port s3 s3-eth2
ovs-vsctl add-port s3 s3-eth3
ovs-vsctl add-port s4 s4-eth1
ovs-vsctl add-port s4 s4-eth2

# 6. Connect controller to switch
ovs-vsctl set-controller s1 tcp:onos_ip:6653
ovs-vsctl set-controller s2 tcp:onos_ip:6653
ovs-vsctl set-controller s3 tcp:onos_ip:6653
ovs-vsctl set-controller s4 tcp:onos_ip:6653

# 7.1 Setup ip for hosts
ip netns exec h1 ifconfig h1-eth1 192.168.1.10 hw ether 0E:76:73:DE:95:0B
ip netns exec h1 ifconfig lo up
ip netns exec h2 ifconfig h2-eth1 192.168.1.20 hw ether 46:F1:1F:61:13:BE
ip netns exec h2 ifconfig lo up
ip netns exec h3 ifconfig h3-eth1 192.168.1.30 hw ether C2:CF:6D:1F:04:68
ip netns exec h3 ifconfig lo up
ip netns exec h4 ifconfig h4-eth1 192.168.1.40 hw ether B6:9D:D9:84:A7:64
ip netns exec h4 ifconfig lo up
ip netns exec h5 ifconfig h5-eth1 192.168.1.50 hw ether 42:00:FB:6C:9D:E3
ip netns exec h5 ifconfig lo up
ip netns exec h6 ifconfig h6-eth1 192.168.1.60 hw ether 1A:16:D9:D6:FA:6E
ip netns exec h6 ifconfig lo up

# 7.2 Setup ip for switches
ifconfig s1-eth1 up
ifconfig s1-eth2 up
ifconfig s1-eth3 up
ifconfig s1-eth4 up
ifconfig s1-eth5 up
ifconfig s2-eth1 up
ifconfig s2-eth2 up
ifconfig s3-eth1 up
ifconfig s3-eth2 up
ifconfig s3-eth3 up
ifconfig s4-eth1 up
ifconfig s4-eth2 up

# 8 Test connections between hosts
ip netns exec h1 ping -c1 192.168.1.10
ip netns exec h1 ping -c1 192.168.1.20
ip netns exec h1 ping -c1 192.168.1.30
ip netns exec h1 ping -c1 192.168.1.40
ip netns exec h1 ping -c1 192.168.1.50
ip netns exec h1 ping -c1 192.168.1.60
