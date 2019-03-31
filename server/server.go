package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"secure-grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Cert file details
const (
	Addr       = ":50052"
	AkbarCA    = "cert/akbar.com.crt"
	ServerCert = "cert/server.crt"
	ServerKey  = "cert/server.key"
	ServerName = "server"
)

// MaxServer to define the find max number.
type MaxServer struct{}

// Num to find the max number among the stream of  numbers received.
func (m *MaxServer) Num(srv proto.Max_NumServer) error {
	ctx := srv.Context()
	var max int32

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// receive from the client stream.
		r, err := srv.Recv()
		if err == io.EOF {
			log.Print("exiting...")
			return nil
		}
		if err != nil {
			log.Printf("received error : %v", err)
			continue
		}

		// check max number with prev max.
		if r.Num <= max {
			continue
		}

		// update max number.
		max = r.Num

		// Send to client stream
		if err := srv.Send(&proto.Response{Result: max}); err != nil {
			log.Printf("send error : %v", err)
		}
		log.Printf("Sent Max number : %d", max)
	}
}

// run starts the secure TLS grpc server
func run() error {

	// Load certs from the disk.
	cert, err := tls.LoadX509KeyPair(ServerCert, ServerKey)
	if err != nil {
		return fmt.Errorf("could not server key pairs: %s", err)
	}

	// Create certpool from the CA
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(AkbarCA)
	if err != nil {
		return fmt.Errorf("could not read CA cert: %s", err)
	}

	// Append the certs from the CA
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return fmt.Errorf("Failed to append the CA certs: %s", err)
	}

	// Create the TLS config for gRPC server.
	creds := credentials.NewTLS(
		&tls.Config{
			ClientAuth:   tls.RequireAnyClientCert,
			Certificates: []tls.Certificate{cert},
			ClientCAs:    certPool,
		})

	// Create the new gRPC server with TLS config
	gSrv := grpc.NewServer(grpc.Creds(creds))
	proto.RegisterMaxServer(gSrv, &MaxServer{})

	// Open channel to listen on the addr
	list, err := net.Listen("tcp", Addr)
	if err != nil {
		return fmt.Errorf("Could not listen on port: %s, Error %v", Addr, err)
	}

	// Serve accepts incoming connections on the listener
	if err := gSrv.Serve(list); err != nil {
		return fmt.Errorf("Failed to serve grpc, Error:%v", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(-1)
	}
}
