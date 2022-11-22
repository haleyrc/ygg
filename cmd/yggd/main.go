package main

import (
	"context"
	"fmt"
	"net/http"

	"tailscale.com/client/tailscale"

	"github.com/haleyrc/ygg/server"
	"github.com/haleyrc/ygg/share"
)

const DEFAULT_PORT = 8080

func main() {
	ctx := context.Background()

	status, err := tailscale.Status(ctx)
	if err != nil {
		panic(err)
	}

	if len(status.TailscaleIPs) == 0 {
		panic("No Tailscale IPs found!")
	}

	iface := status.TailscaleIPs[0]
	fmt.Printf("Tailscale is up on: %s\n", iface)

	srv, err := server.New(ctx, iface.String(), DEFAULT_PORT)
	if err != nil {
		panic(err)
	}

	srv.Sites().Mount("/share", share.NewSite())

	srv.Routes()

	ip := fmt.Sprintf("%s:%d", iface.String(), DEFAULT_PORT)
	fmt.Printf("Listening on http://%s...\n", ip)
	if err := http.ListenAndServe(ip, srv); err != nil {
		panic(err)
	}
}
