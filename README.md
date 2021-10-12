# middleboxer

middleboxer is a command line tool for testing middleboxes like firewalls. It
uses cooperating client instances that send packets through a middlebox and
checks which packets pass through the middlebox and which do not.

## Overview

A test setup using middleboxer is depicted in the figure below:

```
+----------+     +-----------+     +----------+
| Client 1 |---->| Middlebox |---->| Client 2 |
+----------+     +-----------+     +----------+
     |_______________    _______________|
                     |  |
                  +--------+
                  | Server |
                  +--------+
```

The `middlebox` is the device under test, e.g., a firewall. Two middleboxer
clients, `client 1` and `client 2`, are connected to a middleboxer `server`.
Both clients are placed in such a way that they communicate with each other
through the middlebox.

The server manages the properties of packets, e.g., IP addresses, layer 4
protocol and ports, that should be sent through the middlebox by the clients.
It instructs one client to send those packets and the other client to receive
them.

The clients report back to the server for each packet if the packet passed
through the middlebox or if they received error messages like ICMP errors or
TCP resets.

The server collects all results and prints them to the console or writes them
to a file.

## Usage

You can run the `middleboxer` command with the following command line
arguments:

```
Usage of middleboxer:
  -address string
        set address to connect to (client mode) or listen on (server mode)
  -diffs
        show packet diffs in results
  -id uint
        set id of the client (default 1)
  -out string
        set output file
  -ports string
        set port range to be tested (default "1:65535")
  -prot uint
        set layer 4 protocol to 6 (tcp) or 17 (udp) (default 6)
  -rdev string
        set device of the receiving client
  -rdip string
        set destination IP of the receiving client
  -rdmac string
        set destination MAC of the receiving client
  -rid uint
        set id of the receiving client (default 2)
  -rsip string
        set source IP of the receiving client
  -rsmac string
        set source MAC of the receiving client
  -sdev string
        set device of the sending client
  -sdip string
        set destination IP of the sending client
  -sdmac string
        set destination MAC of the sending client
  -server
        run as server (default: run as client)
  -sid uint
        set id of the sending client (default 1)
  -ssip string
        set source IP of the sending client
  -ssmac string
        set source MAC of the sending client
  -ssport uint
        set source port of the sending client
```

## Examples

### Server

Running a middleboxer server listening on all addresses and port 3333:

```console
$ middleboxer -server -address :3333 \
	-sdev veth2 \
	-ssmac 0a:bc:de:f0:00:12 \
	-sdmac 0a:bc:de:f0:00:22 \
	-ssip 192.168.1.1 \
	-sdip 192.168.1.2 \
	-ssport 4242 \
	-rdev veth4 \
	-rsip 192.168.1.1 \
	-rdip 192.168.1.2 \
	-prot 17 \
	-ports 1024:1032
```

The command line above configures the following properties of packets sent by
the sending client:
* Outgoing network device `veth2`
* Source MAC address `0a:bc:de:f0:00:12`
* Destination MAC address `0a:bc:de:f0:00:22`
* Source IP address `192.168.1.1`
* Destination IP address `192.168.1.2`
* Source Port `4242`

The command line above configures the following properties of packets expected
by the receiving client:
* Incoming network device `veth4`
* Source IP address `192.168.1.1`
* Destination IP address `192.168.1.2`

Additionally, the command line above specifies that UDP should be used (`-prot
17`) and that all ports from 1024 to 1032 should be tested (`-ports
1024:1032`).

### Clients

Running a client with ID 1 and connecting to the server listening on
`192.168.1.3:3333`:

```console
$ sudo middleboxer -id 1 -address 192.168.1.3:3333
```

Running a client with ID 2 and connecting to the server listening on
`192.168.1.3:3333`:

```console
$ sudo middleboxer -id 2 -address 192.168.1.3:3333
```
