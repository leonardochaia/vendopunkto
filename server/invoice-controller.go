package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/leonardochaia/vendopunkto/store"
)

func (s *Server) CreateInvoice() http.HandlerFunc {

	type creationParams struct {
		Amount       uint   `json:"amount"`
		Denomination string `json:"denomination"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var params = new(creationParams)
		if err := render.DecodeJSON(r.Body, &params); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		invoice, err := s.invoice.CreateInvoice(r.Context(), params.Amount, params.Denomination)
		if err != nil {
			render.Render(w, r, s.ErrInternalLog(err))
			return
		}
		render.JSON(w, r, invoice)
	}

}

func (s *Server) GetInvoice() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		thingID := chi.URLParam(r, "id")
		if thingID == "" {
			render.Render(w, r, ErrInvalidRequest(fmt.Errorf("Invalid ID")))
			return
		}

		b, err := s.invoice.GetInvoice(r.Context(), thingID)
		
		if err == store.ErrNotFound {
			render.Render(w, r, ErrNotFound)
			return
		} else if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		render.JSON(w, r, b)
	}

}

// // ThingSave creates or updates a thing
// func (s *Server) ThingSave() http.HandlerFunc {

// 	type idResponse struct {
// 		ID string `json:"id"`
// 	}

// 	return func(w http.ResponseWriter, r *http.Request) {

// 		var b = new(gorestapi.Thing)
// 		if err := render.DecodeJSON(r.Body, &b); err != nil {
// 			render.Render(w, r, ErrInvalidRequest(err))
// 			return
// 		}
// 		thingID, err := s.thingStore.ThingSave(r.Context(), b)
// 		if err != nil {
// 			render.Render(w, r, ErrInvalidRequest(err))
// 			return
// 		}
// 		render.JSON(w, r, &idResponse{ID: thingID})
// 	}

// }

// // ThingDelete deletes a thing
// func (s *Server) ThingDelete() http.HandlerFunc {

// 	return func(w http.ResponseWriter, r *http.Request) {

// 		// Get the thingID
// 		thingID := chi.URLParam(r, "id")
// 		if thingID == "" {
// 			render.Render(w, r, ErrInvalidRequest(fmt.Errorf("Invalid ID")))
// 			return
// 		}
// 		err := s.thingStore.ThingDeleteByID(r.Context(), thingID)
// 		if err != nil {
// 			render.Render(w, r, ErrInvalidRequest(err))
// 			return
// 		}

// 		render.NoContent(w, r)
// 	}

// }
