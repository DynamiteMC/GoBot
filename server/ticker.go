package server

import (
	"time"
)

type Ticker struct {
	srv *Server

	TickingFrequency int
	tick             uint

	ticker *time.Ticker
}

func (srv *Server) createTicker() {
	srv.ticker = Ticker{
		srv:              srv,
		TickingFrequency: srv.cfg.TPS,
		ticker:           time.NewTicker(time.Second / time.Duration(srv.cfg.TPS)),
	}
}

func (t Ticker) Start() {
	go func() {
		for range t.ticker.C {
			t.tick++
		}
	}()
}
