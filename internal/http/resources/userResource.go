package resources

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	lru "github.com/hashicorp/golang-lru"
	"net/http"
	"project/internal/models"
	"project/internal/pkg"
	"project/internal/store"
	"strconv"
)

type UserResource struct {
	store store.Store
	cache *lru.TwoQueueCache
}

func NewUserResource(store store.Store, cache *lru.TwoQueueCache) *UserResource {
	return &UserResource{
		store: store,
		cache: cache,
	}
}

func (ur *UserResource) Routes(auth func(handler http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	r.Post("/registration", ur.CreateUser)
	r.Group(func(r chi.Router) {
		r.Use(auth)
		r.Get("/", ur.AllUsers)
		r.Put("/", ur.UpdateUser)
		r.Delete("/{id}", ur.DeleteUser)
	})
	return r

}

func (ur *UserResource) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := new(models.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if _, err := ur.store.Users().ByEmail(r.Context(), user.Email); err == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error. User already exists")
		return
	}

	err := ur.store.Users().Create(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ur *UserResource) AllUsers(w http.ResponseWriter, r *http.Request) {
	if !pkg.IsUserAdmin(r.Context(), w) {
		return
	}
	users, err := ur.store.Users().All(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB error: %v", err)
		return
	}
	render.JSON(w, r, users)
}

func (ur *UserResource) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := new(models.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	if _, err := ur.store.Users().ByEmail(r.Context(), user.Email); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error. User doesn't exist")
		return
	}
	userInfo := r.Context().Value(pkg.CtxKeyUser).(*models.AuthorizedInfo)
	user.ID = userInfo.Id

	if err := ur.store.Users().Update(r.Context(), user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB error: %v", err)
		return
	}
}

func (ur *UserResource) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	if err = ur.store.Users().Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB error : %v", err)
		return
	}
}
