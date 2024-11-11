package api

import (
	"fmt"
	"net/http"
	db "simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,oneof=USD INR CAD"`
}

func (s *Server) TransferAmount(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := s.store.TransferTx(ctx, db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if !s.validAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}

	if !s.validAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)

	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}

// func (server *Server) checkCurrency(ctx *gin.Context, accountID int64, currency string) bool {
// 	account, err := server.store.GetAccount(ctx, accountID)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// 			return false
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return false
// 	}

// 	// if account.Currency != currency {
// 	// 	err := "currency mismath"
// 	// 	ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 	// 	return false
// 	// }

// 	return true
// }