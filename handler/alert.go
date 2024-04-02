package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/malav4all/golang-api/model"
	"github.com/malav4all/golang-api/repository/alert"
)

type Alert struct {
	Repo *alert.RedisRepo
}

func (a *Alert) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UserId    uuid.UUID         `json:"customer_id"`
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

func (a *Alert) ListAlert(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List Alert")
}

func (a *Alert) getAlertId(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Alert ID")
}

func (a *Alert) UpdateById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update Alert")
}

func (a *Alert) DeleteById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete Alert")
}
