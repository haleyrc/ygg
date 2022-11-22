package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"tailscale.com/client/tailscale"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		fmt.Printf("--> %s\n", r.RemoteAddr)
		fmt.Fprintln(w, "Hello!")
		fmt.Printf("<-- %s\n", time.Since(start))
	})

	ip := fmt.Sprintf("%s:%d", iface.String(), DEFAULT_PORT)
	fmt.Printf("Listening on http://%s...\n", ip)
	if err := http.ListenAndServe(ip, nil); err != nil {
		panic(err)
	}
}
