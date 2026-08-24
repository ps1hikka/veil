package server

import (
	"log"
	"net"

	"gateway/internal/auth"
	"gateway/internal/config"
	"gateway/internal/proxy"
	"gateway/internal/session"
	"gateway/internal/vless"
)

func ListenAndServe(cfg config.Config) error {
	authenticator, err := auth.New(cfg)
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", cfg.Listen)
	if err != nil {
		return err
	}
	log.Println("gateway listening on", cfg.Listen)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept:", err)
			continue
		}
		go handle(conn, cfg, authenticator)
	}
}

func handle(conn net.Conn, cfg config.Config, a *auth.Authenticator) {
	defer conn.Close()

	header, err := vless.ParseHeader(conn)
	if err != nil {
		log.Println("vless:", err)
		return
	}

	client, ok := a.Allow(header.UUID)
	if !ok {
		log.Println("auth failed:", vless.FormatUUID(header.UUID))
		return
	}

	log.Printf("accepted uuid=%s cmd=%d dest=%s:%d",
		vless.FormatUUID(header.UUID),
		header.Command,
		header.Address,
		header.Port,
	)

	sess := session.New(conn, header, client)
	if err := proxy.Handle(sess); err != nil {
		log.Println("proxy:", err)
	}
}
