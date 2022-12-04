package http

import (
	"context"
	"github.com/go-chi/chi"
	lru "github.com/hashicorp/golang-lru"
	"log"
	"net/http"
	"project/internal/store"
	"time"
)

type Server struct {
	ctx         context.Context
	idleConnsCH chan struct{}
	store       store.Store
	cache       *lru.TwoQueueCache
	Address     string
}

func NewServer(ctx context.Context, opts ...ServerOption) *Server {
	srv := &Server{
		ctx:         ctx,
		idleConnsCH: make(chan struct{}),
	}
	for _, opts := range opts {
		opts(srv)
	}

	return srv
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()
	brandsResource := NewBrandResources(s.store, s.cache)
	r.Mount("/brands", brandsResource.Routes())
	carsResource := NewCarResource(s.store, s.cache)
	r.Mount("/cars", carsResource.Routes())
	return r
}

func (s *Server) Run() error {

	srv := &http.Server{
		Addr:         s.Address,
		Handler:      s.basicHandler(),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 30,
	}
	go s.ListenCtxForGt(srv)

	log.Println("Server running on ", s.Address)
	return srv.ListenAndServe()
}

func (s *Server) ListenCtxForGt(srv *http.Server) {
	<-s.ctx.Done() // блокируемся пока контекст приложения не отменен

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("[HTTTP] GOt err while shutting down %v", err)
		return
	}

	log.Println("[HTTP] Processed all idle connections")
	close(s.idleConnsCH)
}

func (s *Server) WaitForGracefulTermination() {
	//блок до записи или закрытия канала
	<-s.idleConnsCH
}
