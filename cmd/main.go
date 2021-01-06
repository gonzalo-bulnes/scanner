package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gonzalo-bulnes/scanner/cmd/domain"
)

func main() {
	d := domain.New("apt.freedom.press") // Proof of concept with a single domain

	status := d.Check(context.Background())
	if status.Err() != nil {
		os.Exit(2)
	}

	metadata := status.Value().(domain.Metadata)
	if metadata.CertificateValidUntil.Before(time.Now().Add(168 * time.Hour)) {
		fmt.Printf("The TLS certificate for %s expires soon! It is valid until %v.\n", metadata.Name, metadata.CertificateValidUntil)
		os.Exit(1)
	}
	fmt.Printf("The TLS certificate for %s is valid until %v.\n", metadata.Name, metadata.CertificateValidUntil)
}
