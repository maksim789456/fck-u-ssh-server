package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
)

func main() {
	config := &ssh.ServerConfig{
		PasswordCallback: func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			return nil, nil
		},
	}

	privateBytes, err := os.ReadFile("id_rsa")
	if err != nil {
		fmt.Printf("Failed to load private key\n")
		return
	}
	fmt.Printf("Loaded private key\n")

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		fmt.Printf("Failed to parse private key\n")
		return
	}
	fmt.Printf("Parsed private key\n")

	config.AddHostKey(private)

	in, out := make(chan string), make(chan string)
	for {
		listener, err := net.Listen("tcp", "0.0.0.0:22")
		if err != nil {
			fmt.Printf("Failed to accept incoming connection: %s\n", err)
			continue
		}

		nConn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept incoming connection\n")
			continue
		}

		serviceConnection(nConn, config, in, out)
		listener.Close()
	}
}

func serviceConnection(nConn net.Conn, config *ssh.ServerConfig, in, out chan string) {
	sshConn, chans, reqs, err := ssh.NewServerConn(nConn, config)
	if err != nil {
		return
	}

	fmt.Printf("Trying to connect: %s@%s (%s)\n", sshConn.User(), sshConn.RemoteAddr(), sshConn.ClientVersion())

	go func() {
		for req := range reqs {
			fmt.Printf("Discarding request type (%s), with playload:\n%s\n", req.Type, req.Payload)
			if req.WantReply {
				req.Reply(false, nil)
			}
		}
	}()

	for newChannel := range chans {
		if t := newChannel.ChannelType(); t != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "Only session type!")
			continue
		}

		channel, _, err := newChannel.Accept()
		if err != nil {
			fmt.Printf("Could not accept channel: %s\n", err)
			return
		}

		go func() {
			var bytes = []byte("^_^\r\n") // Change to you own message
			channel.Write(bytes)
			channel.Close()
		}()
	}
}
