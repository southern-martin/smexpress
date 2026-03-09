package http

import (
	"net/http"

	"github.com/smexpress/services/customer-service/internal/interface/http/handler"
)

type Handlers struct {
	Customer *handler.CustomerHandler
	Contact  *handler.ContactHandler
	Address  *handler.AddressHandler
	Note     *handler.NoteHandler
}

func NewRouter(h Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"customer-service"}`))
	})

	// Customers
	mux.HandleFunc("GET /customers", h.Customer.List)
	mux.HandleFunc("POST /customers", h.Customer.Create)
	mux.HandleFunc("GET /customers/{id}", h.Customer.Get)
	mux.HandleFunc("PUT /customers/{id}", h.Customer.Update)
	mux.HandleFunc("DELETE /customers/{id}", h.Customer.Delete)
	mux.HandleFunc("PATCH /customers/{id}/activate", h.Customer.Activate)
	mux.HandleFunc("PATCH /customers/{id}/deactivate", h.Customer.Deactivate)
	mux.HandleFunc("PATCH /customers/{id}/credit-hold", h.Customer.SetCreditHold)

	// Contacts
	mux.HandleFunc("GET /customers/{id}/contacts", h.Contact.List)
	mux.HandleFunc("POST /customers/{id}/contacts", h.Contact.Create)
	mux.HandleFunc("PUT /customers/{id}/contacts/{contactId}", h.Contact.Update)
	mux.HandleFunc("DELETE /customers/{id}/contacts/{contactId}", h.Contact.Delete)

	// Addresses
	mux.HandleFunc("GET /customers/{id}/addresses", h.Address.List)
	mux.HandleFunc("POST /customers/{id}/addresses", h.Address.Create)
	mux.HandleFunc("PUT /customers/{id}/addresses/{addressId}", h.Address.Update)
	mux.HandleFunc("DELETE /customers/{id}/addresses/{addressId}", h.Address.Delete)

	// Notes
	mux.HandleFunc("GET /customers/{id}/notes", h.Note.List)
	mux.HandleFunc("POST /customers/{id}/notes", h.Note.Create)
	mux.HandleFunc("DELETE /customers/{id}/notes/{noteId}", h.Note.Delete)

	return mux
}
