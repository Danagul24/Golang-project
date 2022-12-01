package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
	lru "github.com/hashicorp/golang-lru"
	"net/http"
	"project/internal/models"
	"project/internal/store"
	"strconv"
)

type BrandResource struct {
	store store.Store
	cache *lru.TwoQueueCache
}

func NewBrandResources(store store.Store, cache *lru.TwoQueueCache) *BrandResource {
	return &BrandResource{
		store: store,
		cache: cache,
	}
}

func (br *BrandResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", br.CreateBrand)
	r.Get("/", br.AllBrands)
	r.Get("/{id}", br.ByID)
	r.Put("/", br.UpdateBrand)
	r.Delete("/{id}", br.DeleteBrand)
	return r
}

func (br *BrandResource) CreateBrand(w http.ResponseWriter, r *http.Request) {
	brand := new(models.Brand)
	if err := json.NewDecoder(r.Body).Decode(brand); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := br.store.Brands().Create(r.Context(), brand); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB error : %v", err)
		return
	}

	br.cache.Purge() // чистка кэша после создания бренда

	w.WriteHeader(http.StatusCreated)
}

func (br *BrandResource) AllBrands(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.BrandFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		brandsFromCache, ok := br.cache.Get(searchQuery)
		if ok {
			render.JSON(w, r, brandsFromCache)
			return
		}
		filter.Query = &searchQuery
	}

	brands, err := br.store.Brands().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
	if searchQuery != "" {
		br.cache.Add(searchQuery, brands)
	}
	render.JSON(w, r, brands)
}
func (br *BrandResource) ByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	brandFromCache, ok := br.cache.Get(id)
	if ok {
		render.JSON(w, r, brandFromCache)
		return
	}

	brand, err := br.store.Brands().ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	br.cache.Add(id, brand)
	render.JSON(w, r, brand)
}

func (br *BrandResource) UpdateBrand(w http.ResponseWriter, r *http.Request) {
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

	if err := br.store.Brands().Update(r.Context(), brand); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB error: %v", err)
		return
	}

	br.cache.Remove(brand.ID)
}

func (br *BrandResource) DeleteBrand(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	if err := br.store.Brands().Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB error %v", err)
		return
	}

	br.cache.Remove(id)
}
