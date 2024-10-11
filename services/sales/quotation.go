package sales

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/lib/pq"
	"github.com/nullism/bqb"
	"system.buon18.com/m/database"
	"system.buon18.com/m/models"
	"system.buon18.com/m/models/sales"
	"system.buon18.com/m/models/setting"
	"system.buon18.com/m/utils"
)

type SalesQuotationService struct {
	DB *sql.DB
}

func (service *SalesQuotationService) Quotations(qp *utils.QueryParams) ([]sales.SalesQuotationResponse, int, int, error) {
	bqbQuery := bqb.New(`
	WITH "limited_quotations" AS (
		SELECT
			*
		FROM
			"sales.quotation"`)

	qp.FilterIntoBqb(bqbQuery)
	qp.PaginationIntoBqb(bqbQuery)

	bqbQuery.Space(`)
	SELECT
		"limited_quotations".id,
		"limited_quotations".name,
		"limited_quotations".creation_date,
		"limited_quotations".validity_date,
		"limited_quotations".discount,
		"limited_quotations".amount_delivery,
		"limited_quotations".status,
		"setting.customer".id,
		"setting.customer".full_name,
		"setting.customer".gender,
		"setting.customer".email,
		"setting.customer".phone,
		"setting.customer".additional_information,
		"sales.order_item".id,
		"sales.order_item".name,
		"sales.order_item".description,
		"sales.order_item".price,
		"sales.order_item".discount
	FROM
		"limited_quotations"
	INNER JOIN "setting.customer" ON "limited_quotations".setting_customer_id = "setting.customer".id
	LEFT JOIN "sales.order_item" ON "limited_quotations"."id" = "sales.order_item".sales_quotation_id`)

	qp.OrderByIntoBqb(bqbQuery, `"limited_quotations".id ASC, "sales.order_item".id ASC`)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return []sales.SalesQuotationResponse{}, 0, http.StatusInternalServerError, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return []sales.SalesQuotationResponse{}, 0, http.StatusInternalServerError, utils.ErrInternalServer
	}

	quotationsResponse := make([]sales.SalesQuotationResponse, 0)
	lastQuotation := sales.SalesQuotation{}
	lastCustomer := setting.SettingCustomer{}
	orderItems := make([]sales.SalesOrderItem, 0)
	for rows.Next() {
		var tmpQuotation sales.SalesQuotation
		var tmpCustomer setting.SettingCustomer
		var tmpOrderItem sales.SalesOrderItem

		err = rows.Scan(&tmpQuotation.Id, &tmpQuotation.Name, &tmpQuotation.CreationDate, &tmpQuotation.ValidityDate, &tmpQuotation.Discount, &tmpQuotation.AmountDelivery, &tmpQuotation.Status, &tmpCustomer.Id, &tmpCustomer.FullName, &tmpCustomer.Gender, &tmpCustomer.Email, &tmpCustomer.Phone, &tmpCustomer.AdditionalInformation, &tmpOrderItem.Id, &tmpOrderItem.Name, &tmpOrderItem.Description, &tmpOrderItem.Price, &tmpOrderItem.Discount)
		if err != nil {
			slog.Error(fmt.Sprintf("%v", err))
			return []sales.SalesQuotationResponse{}, 0, http.StatusInternalServerError, utils.ErrInternalServer
		}

		if lastQuotation.Id != tmpQuotation.Id && lastQuotation.Id != nil {
			customerResponse := setting.SettingCustomerToResponse(lastCustomer)
			orderItemsResponse := make([]sales.SalesOrderItemResponse, 0)
			for _, item := range orderItems {
				orderItemsResponse = append(orderItemsResponse, sales.SalesOrderItemToResponse(item))
			}
			quotationsResponse = append(quotationsResponse, sales.SalesQuotationToResponse(lastQuotation, &customerResponse, &orderItemsResponse))

			// Reset and append new data
			lastQuotation = tmpQuotation
			lastCustomer = tmpCustomer
			orderItems = make([]sales.SalesOrderItem, 0)
			orderItems = append(orderItems, tmpOrderItem)
			continue
		}

		if lastQuotation.Id == nil {
			lastQuotation = tmpQuotation
			lastCustomer = tmpCustomer
		}

		if tmpOrderItem.Id != nil {
			orderItems = append(orderItems, tmpOrderItem)
		}
	}
	if lastQuotation.Id != nil {
		customerResponse := setting.SettingCustomerToResponse(lastCustomer)
		orderItemsResponse := make([]sales.SalesOrderItemResponse, 0)
		for _, item := range orderItems {
			orderItemsResponse = append(orderItemsResponse, sales.SalesOrderItemToResponse(item))
		}
		quotationsResponse = append(quotationsResponse, sales.SalesQuotationToResponse(lastQuotation, &customerResponse, &orderItemsResponse))
	}

	bqbQuery = bqb.New(`SELECT COUNT(*) FROM "sales.quotation"`)
	qp.FilterIntoBqb(bqbQuery)

	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return []sales.SalesQuotationResponse{}, 0, http.StatusInternalServerError, utils.ErrInternalServer
	}

	var total int
	err = service.DB.QueryRow(query, params...).Scan(&total)
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return []sales.SalesQuotationResponse{}, 0, http.StatusInternalServerError, utils.ErrInternalServer
	}

	return quotationsResponse, total, http.StatusOK, nil
}

func (service *SalesQuotationService) Quotation(id string) (sales.SalesQuotationResponse, int, error) {
	bqbQuery := bqb.New(`
	WITH "limited_quotations" AS (
		SELECT
			*
		FROM
			"sales.quotation"
		WHERE id = ?)
	SELECT
		"limited_quotations".id,
		"limited_quotations".name,
		"limited_quotations".creation_date,
		"limited_quotations".validity_date,
		"limited_quotations".discount,
		"limited_quotations".amount_delivery,
		"limited_quotations".status,
		"setting.customer".id,
		"setting.customer".full_name,
		"setting.customer".gender,
		"setting.customer".email,
		"setting.customer".phone,
		"setting.customer".additional_information,
		"sales.order_item".id,
		"sales.order_item".name,
		"sales.order_item".description,
		"sales.order_item".price,
		"sales.order_item".discount
	FROM
		"limited_quotations"
	INNER JOIN "setting.customer" ON "limited_quotations".setting_customer_id = "setting.customer".id
	LEFT JOIN "sales.order_item" ON "limited_quotations"."id" = "sales.order_item".sales_quotation_id
	ORDER BY "limited_quotations".id ASC, "sales.order_item".id ASC`, id)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return sales.SalesQuotationResponse{}, http.StatusInternalServerError, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return sales.SalesQuotationResponse{}, http.StatusInternalServerError, utils.ErrInternalServer
	}

	quotation := sales.SalesQuotation{}
	customer := setting.SettingCustomer{}
	orderItems := make([]sales.SalesOrderItem, 0)
	for rows.Next() {
		var tmpOrderItem sales.SalesOrderItem
		err = rows.Scan(&quotation.Id, &quotation.Name, &quotation.CreationDate, &quotation.ValidityDate, &quotation.Discount, &quotation.AmountDelivery, &quotation.Status, &customer.Id, &customer.FullName, &customer.Gender, &customer.Email, &customer.Phone, &customer.AdditionalInformation, &tmpOrderItem.Id, &tmpOrderItem.Name, &tmpOrderItem.Description, &tmpOrderItem.Price, &tmpOrderItem.Discount)
		if err != nil {
			slog.Error(fmt.Sprintf("%v", err))
			return sales.SalesQuotationResponse{}, http.StatusInternalServerError, utils.ErrInternalServer
		}

		orderItems = append(orderItems, tmpOrderItem)
	}

	if quotation.Id == nil {
		return sales.SalesQuotationResponse{}, http.StatusNotFound, utils.ErrQuotationNotFound
	}

	customerResponse := setting.SettingCustomerToResponse(customer)
	orderItemsResponse := make([]sales.SalesOrderItemResponse, 0)
	for _, item := range orderItems {
		orderItemsResponse = append(orderItemsResponse, sales.SalesOrderItemToResponse(item))
	}

	return sales.SalesQuotationToResponse(quotation, &customerResponse, &orderItemsResponse), http.StatusOK, nil
}

func (service *SalesQuotationService) CreateQuotation(ctx *utils.CtxValue, quotation *sales.SalesQuotationCreateRequest) (int, error) {
	commonModel := models.CommonModel{}
	commonModel.PrepareForCreate(*ctx.User.Id, *ctx.User.Id)

	tx, err := service.DB.Begin()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}
	defer tx.Rollback()

	bqbQuery := bqb.New(`INSERT INTO "sales.quotation"
	(name, creation_date, validity_date, discount, status, setting_customer_id, cid, ctime, mid, mtime)
	VALUES
	(?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`, quotation.Name, quotation.CreationDate, quotation.ValidityDate, quotation.Discount, quotation.Status, quotation.CustomerId, commonModel.CId, commonModel.CTime, commonModel.MId, commonModel.MTime)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	var id int
	err = tx.QueryRow(query, params...).Scan(&id)
	if err != nil {
		switch err.(*pq.Error).Constraint {
		case database.KEY_SALES_QUOTATION_NAME:
			return http.StatusConflict, utils.ErrQuotationNameExists
		case database.FK_SETTING_CUSTOMER_ID:
			return http.StatusBadRequest, utils.ErrCustomerNotFound
		}

		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	bqbQuery = bqb.New(`INSERT INTO "sales.order_item" (name, description, price, discount, sales_quotation_id, cid, ctime, mid, mtime) VALUES`)

	for index, item := range quotation.SalesOrderItems {
		bqbQuery.Space(`(?, ?, ?, ?, ?, ?, ?, ?, ?)`, item.Name, item.Description, item.Price, item.Discount, id, commonModel.CId, commonModel.CTime, commonModel.MId, commonModel.MTime)
		if index != len(quotation.SalesOrderItems)-1 {
			bqbQuery.Space(",")
		}
	}

	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if _, err = tx.Exec(query, params...); err != nil {
		switch err.(*pq.Error).Constraint {
		case database.FK_SETTING_CUSTOMER_ID:
			return http.StatusBadRequest, utils.ErrCustomerNotFound
		}

		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if err = tx.Commit(); err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	return http.StatusCreated, nil
}

func (service *SalesQuotationService) UpdateQuotation(ctx *utils.CtxValue, id string, quotation *sales.SalesQuotationUpdateRequest) (int, error) {
	commonModel := models.CommonModel{}
	commonModel.PrepareForCreate(*ctx.User.Id, *ctx.User.Id)

	tx, err := service.DB.Begin()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}
	defer tx.Rollback()

	bqbQuery := bqb.New(`SELECT status FROM "sales.quotation" WHERE id = ?`, id)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	var status string
	err = tx.QueryRow(query, params...).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, utils.ErrQuotationNotFound
		}

		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if status == models.SalesQuotationStatusSalesOrder || status == models.SalesQuotationStatusSalesCancelled {
		return http.StatusBadRequest, errors.New("this quotation is not allowed to be updated")
	}

	bqbQuery = bqb.New(`UPDATE "sales.quotation" SET mid = ?, mtime = ?`, commonModel.MId, commonModel.MTime)
	if quotation.Name != nil {
		bqbQuery.Space(`name = ?`, *quotation.Name)
	}
	if quotation.CreationDate != nil {
		bqbQuery.Space(`creation_date = ?`, *quotation.CreationDate)
	}
	if quotation.ValidityDate != nil {
		bqbQuery.Space(`validity_date = ?`, *quotation.ValidityDate)
	}
	if quotation.Discount != nil {
		bqbQuery.Space(`discount = ?`, *quotation.Discount)
	}
	if quotation.Status != nil {
		bqbQuery.Space(`status = ?`, *quotation.Status)
	}
	if quotation.CustomerId != nil {
		bqbQuery.Space(`setting_customer_id = ?`, *quotation.CustomerId)
	}
	bqbQuery.Space(`WHERE id = ?`, id)

	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	result, err := tx.Exec(query, params...)
	if err != nil {
		switch err.(*pq.Error).Constraint {
		case database.KEY_SALES_QUOTATION_NAME:
			return http.StatusConflict, utils.ErrQuotationNameExists
		case database.FK_SETTING_CUSTOMER_ID:
			return http.StatusBadRequest, utils.ErrCustomerNotFound
		}

		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if n, _ := result.RowsAffected(); n == 0 {
		return http.StatusNotFound, utils.ErrQuotationNotFound
	}

	errorChan := make(chan error)
	var wg sync.WaitGroup

	if quotation.AddSalesOrderItems != nil {
		wg.Add(1)
		go func() {
			bqbQuery := bqb.New(`INSERT INTO "sales.order_item" (name, description, price, discount, sales_quotation_id, cid, ctime, mid, mtime) VALUES`)
			for index, item := range *quotation.AddSalesOrderItems {
				bqbQuery.Space(`(?, ?, ?, ?, ?, ?, ?, ?, ?)`, item.Name, item.Description, item.Price, item.Discount, id, commonModel.CId, commonModel.CTime, commonModel.MId, commonModel.MTime)
				if index != len(*quotation.AddSalesOrderItems)-1 {
					bqbQuery.Space(",")
				}
			}

			query, params, err := bqbQuery.ToPgsql()
			if err != nil {
				slog.Error(fmt.Sprintf("%v", err))
			}

			result, err = tx.Exec(query, params...)
			if err != nil {
				slog.Error(fmt.Sprintf("%v", err))
				errorChan <- err
			}

			wg.Done()
		}()
	}

	if quotation.UpdateSalesOrderItems != nil {
		for _, item := range *quotation.UpdateSalesOrderItems {
			wg.Add(1)
			go func() {
				bqbQuery := bqb.New(`UPDATE "sales.order_item" SET mid = ?, mtime = ?`, commonModel.MId, commonModel.MTime)
				if item.Name != nil {
					bqbQuery.Space(`name = ?`, *item.Name)
				}
				if item.Description != nil {
					bqbQuery.Space(`description = ?`, *item.Description)
				}
				if item.Price != nil {
					bqbQuery.Space(`price = ?`, *item.Price)
				}
				if item.Discount != nil {
					bqbQuery.Space(`discount = ?`, *item.Discount)
				}
				bqbQuery.Space(`WHERE id = ? AND sales_quotation_id = ?`, item.Id, id)

				query, params, err := bqbQuery.ToPgsql()
				if err != nil {
					slog.Error(fmt.Sprintf("%v", err))
				}

				_, err = tx.Exec(query, params...)
				if err != nil {
					slog.Error(fmt.Sprintf("%v", err))
					errorChan <- err
				}

				wg.Done()
			}()
		}
	}

	if quotation.DeleteSalesOrderItemIds != nil {
		wg.Add(1)
		go func() {
			bqbQuery := bqb.New(`DELETE FROM "sales.order_item" WHERE id IN (`)
			for index, id := range *quotation.DeleteSalesOrderItemIds {
				bqbQuery.Space(`?`, id)
				if index != len(*quotation.DeleteSalesOrderItemIds)-1 {
					bqbQuery.Space(",")
				}
			}
			bqbQuery.Space(`) AND sales_quotation_id = ?`, id)

			query, params, err := bqbQuery.ToPgsql()
			if err != nil {
				slog.Error(fmt.Sprintf("%v", err))
			}

			_, err = tx.Exec(query, params...)
			if err != nil {
				slog.Error(fmt.Sprintf("%v", err))
				errorChan <- err
			}

			wg.Done()
		}()
	}

	hasError := false
	var errorMessage error
	go func() {
		for err := range errorChan {
			switch err.(*pq.Error).Constraint {
			case database.KEY_SALES_QUOTATION_NAME:
				errorMessage = utils.ErrQuotationNameExists
			case database.FK_SETTING_CUSTOMER_ID:
				errorMessage = utils.ErrCustomerNotFound
			}
			if !hasError {
				hasError = true
			}
		}
	}()

	wg.Wait()
	close(errorChan)

	if hasError {
		switch errorMessage {
		case utils.ErrQuotationNameExists:
			return http.StatusConflict, errorMessage
		case utils.ErrCustomerNotFound:
			return http.StatusBadRequest, errorMessage
		}

		slog.Error(fmt.Sprintf("%v", errorMessage))
		return http.StatusInternalServerError, errorMessage
	}

	err = tx.Commit()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	return http.StatusOK, nil
}

func (service *SalesQuotationService) DeleteQuotation(id string) (int, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}
	defer tx.Rollback()

	bqbQuery := bqb.New(`SELECT status FROM "sales.quotation" WHERE id = ?`, id)
	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	var status string
	err = tx.QueryRow(query, params...).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, utils.ErrQuotationNotFound
		}

		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if status == models.SalesQuotationStatusSalesOrder || status == models.SalesQuotationStatusSalesCancelled {
		return http.StatusForbidden, utils.ErrUnableToDeleteQuotation
	}

	bqbQuery = bqb.New(`DELETE FROM "sales.order_item" WHERE sales_quotation_id = ?`, id)
	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if _, err := tx.Exec(query, params...); err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	bqbQuery = bqb.New(`DELETE FROM "sales.quotation" WHERE id = ?`, id)
	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	result, err := tx.Exec(query, params...)
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if n, _ := result.RowsAffected(); n == 0 {
		return http.StatusNotFound, utils.ErrQuotationNotFound
	}

	err = tx.Commit()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	return http.StatusOK, nil
}
