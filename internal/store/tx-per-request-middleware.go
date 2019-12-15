package store

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-pg/pg"
	"github.com/hashicorp/go-hclog"
	"github.com/leonardochaia/vendopunkto/errors"
)

// Key to use when setting the request ID.
type ctxKeyTransaction int

// RequestIDKey is the key that holds th unique request ID in a request context.
const TransactionPerRequestKey ctxKeyTransaction = 0

func GetTransactionFromContext(ctx context.Context) LazyTransaction {
	if tx := ctx.Value(TransactionPerRequestKey); tx != nil {
		return tx.(LazyTransaction)
	}
	return nil
}

func GetTransactionFromContextOrCreate(ctx context.Context, db *pg.DB) (*pg.Tx, error) {
	tx := GetTransactionFromContext(ctx)
	if tx == nil {
		tx = NewLazyTransaction(db)
	}

	return tx.GetCurrent()
}

func NewTxPerRequestMiddleware(
	parentLogger hclog.Logger,
	db *pg.DB,
) func(next http.Handler) http.Handler {
	logger := parentLogger.Named("tx-per-request")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(errors.WrapHandler(
			func(w http.ResponseWriter, r *http.Request) *errors.APIError {
				// create the transaction
				tx := NewLazyTransaction(db)

				// add it to context
				ctx := r.Context()
				ctx = context.WithValue(ctx, TransactionPerRequestKey, tx)

				var requestID string
				if reqID := r.Context().Value(middleware.RequestIDKey); reqID != nil {
					requestID = reqID.(string)
				}

				// handle panic recovering
				defer func() {
					if rcv := recover(); rcv != nil {
						var err error
						switch r := rcv.(type) {
						case error:
							err = r
						default:
							err = fmt.Errorf("%v", r)
						}

						apiError := errors.InternalServerError(err)
						errors.RenderAPIError(w, r, apiError)

						logger.Warn("Handler errored. Rolling back transaction",
							"requestId", requestID,
							"error", err.Error())

						tx.RollbackIfNeeded()
					}
				}()

				ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
				next.ServeHTTP(ww, r.WithContext(ctx))

				if ww.Status() > 400 {
					logger.Debug("Request returned error status code. Rolling back transaction",
						"requestId", requestID,
						"status", ww.Status())
					tx.RollbackIfNeeded()
				}

				tx.CommitIfNeeded()

				return nil
			}))
	}
}
