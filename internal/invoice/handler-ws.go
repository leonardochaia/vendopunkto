package invoice

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/leonardochaia/vendopunkto/errors"
	"github.com/leonardochaia/vendopunkto/internal/store"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (handler *Handler) invoiceWebSocket(w http.ResponseWriter, r *http.Request) error {
	const op errors.Op = "api.invoice.websocket"

	invoiceID := chi.URLParam(r, "id")
	if invoiceID == "" {
		return errors.E(op, errors.Parameters, errors.Str("No ID was provided"))
	}

	invoice, err := handler.manager.GetInvoice(r.Context(), invoiceID)
	if err != nil {
		return errors.E(op, err)
	}

	if tx := store.GetTransactionFromContext(r.Context()); tx != nil {
		// no need to keep the db tx open since we won't use it.
		err := tx.RollbackIfNeeded()
		if err != nil {
			return errors.E(op, err)
		}
	}

	// upgrade this connection to a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return errors.E(op, errors.Internal, err)
	}
	defer conn.Close()

	// send the current invoice via ws
	err = handler.writeInvoiceToWs(*invoice, conn)
	if err != nil {
		return err
	}

	// used to signal when a read error ocurred, assumed ws closed.
	closedChan := make(chan struct{})

	// keep reading from the socket until an error happens
	go func(closedChan chan struct{}) {
		defer close(closedChan)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					handler.logger.Error("WS read errored", "error", err)
				}
				closedChan <- struct{}{}
				return
			}
		}
	}(closedChan)

	// this channel will get a fresh invoice when it's updated
	invoiceChan := handler.topic.Register(invoiceID)
	defer handler.topic.Unregister(invoiceID)

	// keep listening on the invoiceChan and send back any invoice changes we get
	for {
		select {
		case invoice, ok := <-invoiceChan:
			if !ok {
				return nil
			}
			err = handler.writeInvoiceToWs(invoice, conn)
			if err != nil {
				handler.logger.Error("Error while responding WS", "error", err)
				return nil
			}
		case <-closedChan:
			handler.logger.Debug("Closing websocket", "id", invoiceID)
			return nil
		}
	}
}

func (handler *Handler) writeInvoiceToWs(invoice Invoice, conn *websocket.Conn) error {
	const op errors.Op = "api.invoice.writeInvoiceToWs"
	dto, err := convertInvoiceToDto(invoice, handler.pluginMgr)
	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	err = conn.WriteJSON(dto)
	if err != nil {
		return errors.E(op, errors.Internal, err)
	}

	return nil
}
