package main

import "jello_backend_homework/websocket/client"

func main() {
	// newClient user2 to connect to server.
	// When need to use client test server need to change server "host" and "client".
	// Example:
	// case: docker-compose,  client.NewClient("192.168.99.100", "12345"), host 192.168.99.100 is docker service host.
	// case: minikube,  client.NewClient("192.168.99.101", "30330"), host 192.168.99.101 is minikube ip.
	client.NewClient("192.168.99.101", "30330")
}