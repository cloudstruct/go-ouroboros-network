package main

import (
	"flag"
	"fmt"
	ouroboros "github.com/cloudstruct/go-ouroboros-network"
	"net"
	"os"
)

type serverFlags struct {
	flagset *flag.FlagSet
	//txFile  string
}

func newServerFlags() *serverFlags {
	f := &serverFlags{
		flagset: flag.NewFlagSet("server", flag.ExitOnError),
	}
	//f.flagset.StringVar(&f.txFile, "tx-file", "", "path to the transaction file to submit")
	return f
}

func createListenerSocket(f *globalFlags) (net.Listener, error) {
	var err error
	var listen net.Listener
	if f.socket != "" {
		if err := os.RemoveAll(f.socket); err != nil {
			return nil, fmt.Errorf("failed to remove existing socket: %s", err)
		}
		listen, err = net.Listen("unix", f.socket)
		if err != nil {
			return nil, fmt.Errorf("failed to open listening socket: %s", err)
		}
	} else if f.address != "" {
		listen, err = net.Listen("tcp", f.address)
		if err != nil {
			return nil, fmt.Errorf("failed to open listening socket: %s", err)
		}
	}
	return listen, nil
}

func testServer(f *globalFlags) {
	serverFlags := newServerFlags()
	err := serverFlags.flagset.Parse(f.flagset.Args()[1:])
	if err != nil {
		fmt.Printf("failed to parse subcommand args: %s\n", err)
		os.Exit(1)
	}

	listen, err := createListenerSocket(f)
	if err != nil {
		fmt.Printf("ERROR: failed to create listener: %s\n", err)
		os.Exit(1)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("ERROR: failed to accept connection: %s\n", err)
			continue
		}
		errorChan := make(chan error)
		oOpts := &ouroboros.OuroborosOptions{
			Conn:                  conn,
			NetworkMagic:          uint32(f.networkMagic),
			ErrorChan:             errorChan,
			UseNodeToNodeProtocol: f.ntnProto,
			Server:                true,
		}
		go func() {
			for {
				err := <-errorChan
				fmt.Printf("ERROR: %s\n", err)
			}
		}()
		_, err = ouroboros.New(oOpts)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}
		fmt.Printf("handshake completed...disconnecting\n")
		conn.Close()
	}
}
