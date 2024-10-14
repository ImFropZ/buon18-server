package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"system.buon18.com/m/models"
	"system.buon18.com/m/models/sales"
	"system.buon18.com/m/services"
	"system.buon18.com/m/utils"
)

type SalesHandler struct {
	DB            *sql.DB
	ServiceFacade *services.ServiceFacade
}

func (handler *SalesHandler) Quotations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	qp := utils.NewQueryParams().
		PrepareFilters(sales.SalesQuotation{}, r, `"sales.quotation"`).
		PrepareSorts(sales.SalesQuotation{}, r, `"limited_quotations"`).
		PrepareLimitAndOffset(r)

	quotations, total, statusCode, err := handler.ServiceFacade.SalesQuotationService.Quotations(qp)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.Header().Set("X-Total-Count", fmt.Sprintf("%d", total))
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "", map[string]interface{}{
		"quotations": quotations,
	}))
}

func (handler *SalesHandler) Quotation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	quotation, statusCode, err := handler.ServiceFacade.SalesQuotationService.Quotation(id)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "", map[string]interface{}{
		"quotation": quotation,
	}))
}

func (handler *SalesHandler) CreateQuotation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(utils.CtxKey{}).(*utils.CtxValue)
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[sales.SalesQuotationCreateRequest](r, w, false)
	if !ok {
		return
	}

	statusCode, err := handler.ServiceFacade.SalesQuotationService.CreateQuotation(ctx, &req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "created", nil))
}

func (handler *SalesHandler) UpdateQuotation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(utils.CtxKey{}).(*utils.CtxValue)
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[sales.SalesQuotationUpdateRequest](r, w, true)
	if !ok {
		return
	}

	id := mux.Vars(r)["id"]
	statusCode, err := handler.ServiceFacade.SalesQuotationService.UpdateQuotation(ctx, id, &req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "updated", nil))
}

func (handler *SalesHandler) DeleteQuotations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[models.CommonDelete](r, w, false)
	if !ok {
		return
	}

	statusCode, err := handler.ServiceFacade.SalesQuotationService.DeleteQuotations(&req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "deleted", nil))
}

func (handler *SalesHandler) Orders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	qp := utils.NewQueryParams().
		PrepareFilters(sales.SalesOrder{}, r, `"sales.order"`).
		PrepareSorts(sales.SalesOrder{}, r, `"limited_orders"`).
		PrepareLimitAndOffset(r)

	orders, total, statusCode, err := handler.ServiceFacade.SalesOrderService.Orders(qp)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.Header().Set("X-Total-Count", fmt.Sprintf("%d", total))
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "", map[string]interface{}{
		"orders": orders,
	}))
}

func (handler *SalesHandler) Order(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	order, statusCode, err := handler.ServiceFacade.SalesOrderService.Order(id)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "", map[string]interface{}{
		"order": order,
	}))
}

func (handler *SalesHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(utils.CtxKey{}).(*utils.CtxValue)
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[sales.SalesOrderCreateRequest](r, w, false)
	if !ok {
		return
	}

	statusCode, err := handler.ServiceFacade.SalesOrderService.CreateOrder(ctx, &req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "created", nil))
}

func (handler *SalesHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(utils.CtxKey{}).(*utils.CtxValue)
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[sales.SalesOrderUpdateRequest](r, w, true)
	if !ok {
		return
	}

	id := mux.Vars(r)["id"]
	statusCode, err := handler.ServiceFacade.SalesOrderService.UpdateOrder(ctx, id, &req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "updated", nil))
}
