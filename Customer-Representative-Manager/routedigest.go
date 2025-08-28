package main

import (
	orders "github.com/apouzi/customer-representative-manager/Orders"
	"github.com/go-chi/chi"
)

func RouteDigest(digest *chi.Mux) *chi.Mux {

	digest.Post("/pdf",orders.OrderHandler)
	digest.Post("/create-summary-order",orders.CreateOrderSummaryRecord)

	return digest
}