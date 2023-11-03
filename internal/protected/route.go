package protected

import (
	"net/http"
	custommiddleware "pos/internal/custom_middleware"

	"github.com/go-chi/chi/v5"
)

func Routes() *chi.Mux {
	r := chi.NewMux()
	permissions := []struct {
		name        string
		route       string
		method      string
		description string
		slug        string
	}{
		{"Create Sale", "/create-sale", "POST", "Adding Sale", "create-sale"},
		{"Edit Sale", "/edit-sale", "PUT", "Editing Sale", "edit-sale"},
		{"Refund Transaction", "/refund-transaction", "POST", "Processing Refunds", "refund-transaction"},
		{"View Inventory", "/view-inventory", "GET", "Viewing Inventory", "view-inventory"},
		{"Manage Inventory", "/manage-inventory", "POST", "Adding items to Inventory", "manage-inventory"},
		{"Manage Inventory", "/manage-inventory", "PUT", "Updating Inventory", "manage-inventory"},
		{"Manage Inventory", "/manage-inventory", "DELETE", "Removing items from Inventory", "manage-inventory"},
		{"Generate Reports", "/generate-reports", "GET", "Generating Reports", "generate-reports"},
		{"Generate Reports with ID", "/generate-reports/{id}", "GET", "Generating Reports: {id}", "generate-reports/{id}"},
		{"Customer Management", "/customer-management", "POST", "Adding Customers", "customer-management"},
		{"Customer Management", "/customer-management", "PUT", "Editing Customers", "customer-management"},
		{"Customer Management", "/customer-management", "DELETE", "Deleting Customers", "customer-management"},
		{"User Management", "/user-management", "POST", "Adding Users", "user-management"},
		{"User Management", "/user-management", "PUT", "Editing Users", "user-management"},
		{"User Management", "/user-management", "DELETE", "Deleting Users", "user-management"},
		{"Access Settings", "/access-settings", "GET", "Accessing Settings", "access-settings"},
	}

	for _, p := range permissions {
		r.Group(func(r chi.Router) {
			r.Use(custommiddleware.ProtectedMiddleware(p.slug))
			switch p.method {
			case "GET":
				r.Get(p.route, func(w http.ResponseWriter, req *http.Request) {
					_, err := w.Write([]byte(p.description))
					if err != nil {
						return
					}
				})
			case "POST":
				r.Post(p.route, func(w http.ResponseWriter, req *http.Request) {
					_, err := w.Write([]byte(p.description))
					if err != nil {
						return
					}
				})
			case "PUT":
				r.Put(p.route, func(w http.ResponseWriter, req *http.Request) {
					_, err := w.Write([]byte(p.description))
					if err != nil {
						return
					}
				})
			case "DELETE":
				r.Delete(p.route, func(w http.ResponseWriter, req *http.Request) {
					_, err := w.Write([]byte(p.description))
					if err != nil {
						return
					}
				})
			}
		})
	}

	return r
}
