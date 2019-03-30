package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

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
	ClientCert = "cert/client.crt"
	ClientKey  = "cert/client.key"
	ServerName = "server"
)

func run() error {

	// Load certs from the d
	cert, err := tls.LoadX509KeyPair(ClientCert, ClientKey)
	if err != nil {
		return fmt.Errorf("Could not load client key pair : %v", err)
	}

	// Create certpool from the CA
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(AkbarCA)
	if err != nil {
		return fmt.Errorf("Could not read Cert CA : %v", err)
	}

	// Append the certs from the CA
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return fmt.Errorf("Failed to append CA cert : %v", err)
	}

	// Create transport creds based on TLS.
	creds := credentials.NewTLS(&tls.Config{
		ServerName:   ServerName,
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	})

	// Creates a client connection to the given target
	gConn, err := grpc.Dial(Addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return fmt.Errorf("Could not dail : %s, Error : %v", Addr, err)
	}

	// Create new client
	client := proto.NewMaxClient(gConn)
	s, err := client.Num(context.Background())
	if err != nil {
		return fmt.Errorf("Could not find Max : %v", err)
	}

	var max int32
	ctx := s.Context()
	done := make(chan bool)
	rand.Seed(time.Now().Unix())

	// Send random increasing numbers.
	go func() {
		for i := 1; i < 10; i++ {
			n := int32(rand.Intn(i))
			if err := s.Send(&proto.Request{Num: int32(n)}); err != nil {
				log.Fatalf("Could not send number %d, Error : %v", n, err)
			}
			log.Printf("Sent : %d", n)
			time.Sleep(time.Millisecond * 100)
		}
		if err := s.CloseSend(); err != nil {
			log.Printf("error while closing the send : %v", err)
		}
	}()

	// Receive max number from stream
	go func() {
		for {
			r, err := s.Recv()
			if err == io.EOF {
				close(done)
				return
			}
			if err != nil {
				log.Fatalf("Could not receive max number : %v", err)
			}
			log.Printf("Max Received : %d", r.Result)
			max = r.Result
		}
	}()

	// Close the done channel, If context is done.
	go func() {
		<-ctx.Done()
		if err := ctx.Err(); err != nil {
			log.Printf("Done error : %v", err)
		}
		close(done)
	}()

	// wait for finish
	<-done

	// Print last max number received.
	log.Printf("Finished, Max Received : %d", max)
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	}
}
