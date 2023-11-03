package protected

import (
	"net/http"
	custommiddleware "pos/internal/custom_middleware"

	"github.com/go-chi/chi/v5"
)

func Routes() *chi.Mux {
	r := chi.NewMux()

	r.Use(custommiddleware.ProtectedMiddleware)

	// Create a route for "Create Sale" permission using POST method
	r.Post("/create-sale", func(w http.ResponseWriter, req *http.Request) {
		_, err := w.Write([]byte("Creating Sale"))
		if err != nil {
			return
		}
	})

	// Create a route for "Edit Sale" permission using PUT method
	r.Put("/edit-sale", func(w http.ResponseWriter, req *http.Request) {
		_, err := w.Write([]byte("Editing Sale"))
		if err != nil {
			return
		}
	})

	// Create a route for "Refund Transaction" permission using POST method
	r.Post("/refund-transaction", func(w http.ResponseWriter, req *http.Request) {
		_, err := w.Write([]byte("Processing Refunds"))
		if err != nil {
			return
		}
	})

	// Create a route for "View Inventory" permission using GET method
	r.Get("/view-inventory", func(w http.ResponseWriter, req *http.Request) {
		_, err := w.Write([]byte("Viewing Inventory"))
		if err != nil {
			return
		}
	})

	// Create a route for "Manage Inventory" permission using POST, PUT, and DELETE methods
	r.Post("/manage-inventory", func(w http.ResponseWriter, req *http.Request) {
		_, err := w.Write([]byte("Adding items to Inventory"))
		if err != nil {
			return
		}
	})
	r.Put("/manage-inventory", func(w http.ResponseWriter, req *http.Request) {
		_, err := w.Write([]byte("Updating Inventory"))
		if err != nil {
			return
		}
	})
	r.Delete("/manage-inventory", func(w http.ResponseWriter, req *http.Request) {
		_, err := w.Write([]byte("Removing items from Inventory"))
		if err != nil {
			return
		}
	})

	// Create a route for "Generate Reports" permission using GET method
	r.Get("/generate-reports", func(w http.ResponseWriter, req *http.Request) {
		_, err := w.Write([]byte("Generating Reports"))
		if err != nil {
			return
		}
	})

	// Create a route for "Generate Reports" permission using GET method
	r.Get("/generate-reports/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")
		_, err := w.Write([]byte("Generating Reports: " + id))
		if err != nil {
			return
		}
	})

	// Create a route for "Customer Management" permission using POST, PUT, and DELETE methods
	r.Post("/customer-management", func(w http.ResponseWriter, req *http.Request) {
		_, err := w.Write([]byte("Adding Customers"))
		if err != nil {
			return
		}
	})
	r.Put("/customer-management", func(w http.ResponseWriter, req *http.Request) {
		_, err := w.Write([]byte("Editing Customers"))
		if err != nil {
			return
		}
	})
	r.Delete("/customer-management", func(w http.ResponseWriter, req *http.Request) {
		_, err := w.Write([]byte("Deleting Customers"))
		if err != nil {
			return
		}
	})

	// Create a route for "User Management" permission using POST, PUT, and DELETE methods
	r.Post("/user-management", func(w http.ResponseWriter, req *http.Request) {
		_, err := w.Write([]byte("Adding Users"))
		if err != nil {
			return
		}
	})
	r.Put("/user-management", func(w http.ResponseWriter, req *http.Request) {
		_, err := w.Write([]byte("Editing Users"))
		if err != nil {
			return
		}
	})
	r.Delete("/user-management", func(w http.ResponseWriter, req *http.Request) {
		_, err := w.Write([]byte("Deleting Users"))
		if err != nil {
			return
		}
	})

	// Create a route for "Access Settings" permission using GET method
	r.Get("/access-settings", func(w http.ResponseWriter, req *http.Request) {
		_, err := w.Write([]byte("Accessing Settings"))
		if err != nil {
			return
		}
	})

	return r
}
