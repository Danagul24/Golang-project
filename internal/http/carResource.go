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

type CarResource struct {
	store store.Store
	cache *lru.TwoQueueCache
}

func NewCarResource(store store.Store, cache *lru.TwoQueueCache) *CarResource {
	return &CarResource{
		store: store,
		cache: cache,
	}
}

func (cr *CarResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", cr.CreateCar)
	r.Get("/", cr.AllCars)
	r.Get("/{id}", cr.ByID)
	r.Put("/", cr.UpdateCar)
	r.Delete("/{id}", cr.DeleteCar)
	r.Get("/{city}", cr.FilterCarsByCity)
	r.Get("/sort_by={sortType}", cr.SortCars)
	return r
}

func (cr *CarResource) CreateCar(w http.ResponseWriter, r *http.Request) {
	car := new(models.Car)
	if err := json.NewDecoder(r.Body).Decode(car); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := cr.store.Cars().Create(r.Context(), car); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB error : %v", err)
		return
	}

	cr.cache.Purge() // чистка кэша после создания бренда

	w.WriteHeader(http.StatusCreated)
}

func (cr *CarResource) AllCars(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.CarFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		carsFromCache, ok := cr.cache.Get(searchQuery)
		if ok {
			render.JSON(w, r, carsFromCache)
			return
		}
		filter.Query = &searchQuery
	}

	cars, err := cr.store.Cars().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
	if searchQuery != "" {
		cr.cache.Add(searchQuery, cars)
	}
	render.JSON(w, r, cars)
}
func (cr *CarResource) ByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	carFromCache, ok := cr.cache.Get(id)
	if ok {
		render.JSON(w, r, carFromCache)
		return
	}

	car, err := cr.store.Cars().ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	cr.cache.Add(id, car)
	render.JSON(w, r, car)
}

func (cr *CarResource) UpdateCar(w http.ResponseWriter, r *http.Request) {
	car := new(models.Car)
	if err := json.NewDecoder(r.Body).Decode(car); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	err := validation.ValidateStruct(
		car,
		validation.Field(&car.ID, validation.Required),
	)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err : %v", err)
		return
	}

	if err := cr.store.Cars().Update(r.Context(), car); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB error: %v", err)
		return
	}

	cr.cache.Remove(car.ID)
}

func (cr *CarResource) DeleteCar(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	if err := cr.store.Cars().Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB error %v", err)
		return
	}

	cr.cache.Remove(id)
}

func (cr *CarResource) SortCars(w http.ResponseWriter, r *http.Request) {
	sortType := chi.URLParam(r, "sortType")

	sortedCarsFromCache, ok := cr.cache.Get(sortType)
	if ok {
		render.JSON(w, r, sortedCarsFromCache)
		return
	}

	sortedCars, err := cr.store.Cars().Sort(r.Context(), sortType)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB error: %v", err)
		return
	}

	cr.cache.Add(sortType, sortedCars)
	render.JSON(w, r, sortedCars)
}

func (cr *CarResource) FilterCarsByCity(w http.ResponseWriter, r *http.Request) {
	filter := chi.URLParam(r, "city")

	filteredCarsFromCache, ok := cr.cache.Get(filter)
	if ok {
		render.JSON(w, r, filteredCarsFromCache)
		return
	}

	filteredCars, err := cr.store.Cars().FilterByCity(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB error: %v", err)
		return
	}

	cr.cache.Add(filter, filteredCars)
	render.JSON(w, r, filteredCars)
}
