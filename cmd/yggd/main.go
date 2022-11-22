package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gorilla/mux"
	"tailscale.com/client/tailscale"

	"github.com/haleyrc/ygg/cloudflare/r2"
	"github.com/haleyrc/ygg/server"
	"github.com/haleyrc/ygg/share"
)

const DEFAULT_PORT = 8080

func main() {
	ctx := context.Background()

	iface := mustGetTailscaleInterface()
	fmt.Println("Tailscale is up on:", iface)

	router := mux.NewRouter().StrictSlash(true)
	public := router.PathPrefix("/public").Subrouter()

	r2Client := r2.Must(r2.New(ctx, r2.Credentials{
		AccountID:       os.Getenv("CLOUDFLARE_ACCOUNT_ID"),
		AccessKeyID:     os.Getenv("CLOUDFLARE_ACCESS_KEY_ID"),
		AccessKeySecret: os.Getenv("CLOUDFLARE_ACCESS_KEY_SECRET"),
	}))

	shareController, err := share.NewController(r2Client)
	if err != nil {
		panic(err)
	}
	public.HandleFunc("/share", shareController.Index).Methods("GET")
	public.HandleFunc("/share", shareController.Create).Methods("POST")

	srv, err := server.New(ctx, iface, DEFAULT_PORT, router)
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening on:", srv.Addr())

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func mustGetTailscaleInterface() string {
	status, err := tailscale.Status(context.Background())
	if err != nil {
		panic(err)
	}

	if len(status.TailscaleIPs) == 0 {
		panic("no tailscale ips found")
	}

	iface := status.TailscaleIPs[0]
	return iface.String()
}
