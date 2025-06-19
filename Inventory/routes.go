package main

import (
	retrival "github.com/APouzi/inventory-management/Inventory-Retrieval"
	"github.com/go-chi/chi/v5"
)

func routerDigest(router *chi.Mux){
	router.Get("/inventory-locations",retrival.InventoryLocation)
}


