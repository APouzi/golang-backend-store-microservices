package main

import (
	retrival "github.com/APouzi/inventory-management/Inventory-Retrieval"
	"github.com/go-chi/chi/v5"
)

func routerDigest(router *chi.Mux){
	router.Get("/inventory-locations",retrival.InventoryLocationAll)
	router.Get("/inventory-locations/{id}",retrival.InventoryLocation)
	router.Get("/inventory-locations/",retrival.HandleSearch)
	router.Get("/inventory-locations-products",retrival.AllProductLocations)
	router.Get("/inventory-locations-products/{location}",retrival.ProductsInSingularLocations)
	router.Get("/inventory-locations-transfers",retrival.AllTransfers)
	router.Get("/inventory-locations-transfers/{id}",retrival.TransfersByID)
	router.Get("/inventory-locations-transfers/",retrival.HandleSearch)
}


