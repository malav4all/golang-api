package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/malav4all/golang-api/model"
	"github.com/malav4all/golang-api/repository/alert"
)

type Alert struct {
	Repo *alert.RedisRepo
}

func (a *Alert) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UserId    uuid.UUID         `json:"user_id"`
		AlertData []model.AlertType `json:"alert_data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	now := time.Now().UTC()
	alert := model.Alert{
		AlertID:   rand.Uint64(),
		UserId:    body.UserId,
		AlertData: body.AlertData,
		CreatedAt: &now,
	}
	err := a.Repo.InsertAlert(r.Context(), alert)
	if err != nil {
		fmt.Println("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(alert)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(res)

}

type FindAllPage struct {
	Size   uint64
	Offset uint64
}

type FindResult struct {
	Alerts []model.Alert
	Cursor uint64
}

func (a *Alert) ListAlert(w http.ResponseWriter, r *http.Request) {
	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}

	const decimal = 10
	const bitSize = 64
	cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	const size = 50
	res, err := a.Repo.FindAll(r.Context(), alert.FindAllPage{
		Offset: cursor,
		Size:   size,
	})
	if err != nil {
		fmt.Println("failed to find all:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Items []model.Alert `json:"items"`
		Next  uint64        `json:"next"`
	}
	response.Items = res.Alerts
	response.Next = res.Cursor

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func (a *Alert) GetAlertId(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	fmt.Println("ID:", id)
	const base = 10
	const bitSize = 64
	alertID, err := strconv.ParseUint(id, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Print(alertID)

	o, err := a.Repo.FindByID(r.Context(), alertID)
	if errors.Is(err, alert.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(o); err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *Alert) UpdateById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update Alert")
}

func (a *Alert) DeleteById(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = a.Repo.DeleteByID(r.Context(), orderID)
	if errors.Is(err, alert.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
