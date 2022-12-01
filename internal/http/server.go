package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-ozzo/ozzo-validation"
	"log"
	"net/http"
	"project/internal/models"
	"project/internal/store"
	"strconv"
	"time"
)

type Server struct {
	ctx         context.Context
	idleConnsCH chan struct{}
	store       store.Store
	Address     string
}

func NewServer(ctx context.Context, address string, store store.Store) *Server {
	return &Server{
		ctx:         ctx,
		idleConnsCH: make(chan struct{}),
		store:       store,
		Address:     address,
	}
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()
	r.Post("/brands", func(w http.ResponseWriter, r *http.Request) {
		brand := new(models.Brand)
		if err := json.NewDecoder(r.Body).Decode(brand); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		if err := s.store.Brands().Create(r.Context(), brand); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB error : %v", err)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})
	r.Get("/brands", func(w http.ResponseWriter, r *http.Request) {
		brands, err := s.store.Brands().All(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}
		render.JSON(w, r, brands)
	})
	r.Get("/brands/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		brand, err := s.store.Brands().ByID(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}
		render.JSON(w, r, brand)
	})
	r.Put("/brands", func(w http.ResponseWriter, r *http.Request) {
		brand := new(models.Brand)
		if err := json.NewDecoder(r.Body).Decode(brand); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		err := validation.ValidateStruct(
			brand,
			validation.Field(&brand.ID, validation.Required),
			validation.Field(&brand.Name, validation.Required))
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Unknown err : %v", err)
			return
		}

		if err := s.store.Brands().Update(r.Context(), brand); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB error: %v", err)
			return
		}
	})
	r.Delete("/brands/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		if err := s.store.Brands().Delete(r.Context(), id); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB error %v", err)
			return
		}
	})
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
