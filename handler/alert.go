package handler

import (
	"fmt"
	"net/http"
)

type Alert struct{}

func (a *Alert) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create Alert")
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
