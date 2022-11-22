package main

import (
	"context"
	"fmt"

	"tailscale.com/client/tailscale"

	"github.com/haleyrc/ygg/server"
	"github.com/haleyrc/ygg/share"
)

const DEFAULT_PORT = 8080

func main() {
	ctx := context.Background()

	iface, err := getTailscaleInterface()
	if err != nil {
		panic(err)
	}
	fmt.Println("Tailscale is up on:", iface)

	srv, err := server.New(ctx, iface, DEFAULT_PORT)
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening on:", srv.URL())
	fmt.Printf("Visit: %s/public/\n", srv.URL())

	srv.Sites().Mount("/share", share.NewSite())
	srv.Routes()

	if err := srv.Listen(); err != nil {
		panic(err)
	}
}

func getTailscaleInterface() (string, error) {
	status, err := tailscale.Status(context.Background())
	if err != nil {
		return "", nil
	}

	if len(status.TailscaleIPs) == 0 {
		return "", fmt.Errorf("no tailscale ips found")
	}

	iface := status.TailscaleIPs[0]
	return iface.String(), nil
}
