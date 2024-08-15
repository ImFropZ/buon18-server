package sales

import (
	"database/sql"
	"errors"
	"log"
	"server/database"
	"server/models"
	"server/models/sales"
	"server/models/setting"
	"server/utils"
	"sync"

	"github.com/lib/pq"
	"github.com/nullism/bqb"
)

var (
	ErrQuotationNotFound   = errors.New("quotation not found")
	ErrQuotationNameExists = errors.New("quotation name already exists")
	ErrCustomerNotFound    = errors.New("customer not found")
)

type SalesQuotationService struct {
	DB *sql.DB
}

func (service *SalesQuotationService) Quotations(qp *utils.QueryParams) ([]sales.SalesQuotationResponse, int, int, error) {
	bqbQuery := bqb.New(`
	WITH "limited_quotations" AS (
		SELECT
			id,
			name,
			creation_date,
			validity_date,
			discount,
			amount_delivery,
			status,
			setting_customer_id
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
		"setting.customer".fullname,
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
		log.Printf("%v", err)
		return []sales.SalesQuotationResponse{}, 0, 500, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		log.Printf("%v", err)
		return []sales.SalesQuotationResponse{}, 0, 500, utils.ErrInternalServer
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
			log.Printf("%v", err)
			return []sales.SalesQuotationResponse{}, 0, 500, utils.ErrInternalServer
		}

		if lastQuotation.Id != tmpQuotation.Id && lastQuotation.Id != 0 {
			customerResponse := setting.SettingCustomerToResponse(lastCustomer)
			orderItemsResponse := make([]sales.SalesOrderItemResponse, 0)
			for _, item := range orderItems {
				orderItemsResponse = append(orderItemsResponse, sales.SalesOrderItemToResponse(item))
			}
			quotationsResponse = append(quotationsResponse, sales.SalesQuotationToResponse(lastQuotation, customerResponse, orderItemsResponse))

			// Reset and append new data
			lastQuotation = tmpQuotation
			lastCustomer = tmpCustomer
			orderItems = make([]sales.SalesOrderItem, 0)
			orderItems = append(orderItems, tmpOrderItem)
			continue
		}

		if lastQuotation.Id == 0 {
			lastQuotation = tmpQuotation
			lastCustomer = tmpCustomer
		}

		if tmpOrderItem.Id != 0 {
			orderItems = append(orderItems, tmpOrderItem)
		}
	}
	if lastQuotation.Id != 0 {
		customerResponse := setting.SettingCustomerToResponse(lastCustomer)
		orderItemsResponse := make([]sales.SalesOrderItemResponse, 0)
		for _, item := range orderItems {
			orderItemsResponse = append(orderItemsResponse, sales.SalesOrderItemToResponse(item))
		}
		quotationsResponse = append(quotationsResponse, sales.SalesQuotationToResponse(lastQuotation, customerResponse, orderItemsResponse))
	}

	bqbQuery = bqb.New(`SELECT COUNT(*) FROM "sales.quotation"`)
	qp.FilterIntoBqb(bqbQuery)

	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("%v", err)
		return []sales.SalesQuotationResponse{}, 0, 500, utils.ErrInternalServer
	}

	var total int
	err = service.DB.QueryRow(query, params...).Scan(&total)
	if err != nil {
		log.Printf("%v", err)
		return []sales.SalesQuotationResponse{}, 0, 500, utils.ErrInternalServer
	}

	return quotationsResponse, total, 200, nil
}

func (service *SalesQuotationService) Quotation(id string) (sales.SalesQuotationResponse, int, error) {
	bqbQuery := bqb.New(`
	WITH "limited_quotations" AS (
		SELECT
			id,
			name,
			creation_date,
			validity_date,
			discount,
			amount_delivery,
			status,
			setting_customer_id
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
		"setting.customer".fullname,
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
		log.Printf("%v", err)
		return sales.SalesQuotationResponse{}, 500, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		log.Printf("%v", err)
		return sales.SalesQuotationResponse{}, 500, utils.ErrInternalServer
	}

	quotation := sales.SalesQuotation{}
	customer := setting.SettingCustomer{}
	orderItems := make([]sales.SalesOrderItem, 0)
	for rows.Next() {
		var tmpOrderItem sales.SalesOrderItem
		err = rows.Scan(&quotation.Id, &quotation.Name, &quotation.CreationDate, &quotation.ValidityDate, &quotation.Discount, &quotation.AmountDelivery, &quotation.Status, &customer.Id, &customer.FullName, &customer.Gender, &customer.Email, &customer.Phone, &customer.AdditionalInformation, &tmpOrderItem.Id, &tmpOrderItem.Name, &tmpOrderItem.Description, &tmpOrderItem.Price, &tmpOrderItem.Discount)
		if err != nil {
			log.Printf("%v", err)
			return sales.SalesQuotationResponse{}, 500, utils.ErrInternalServer
		}

		orderItems = append(orderItems, tmpOrderItem)
	}

	if quotation.Id == 0 {
		return sales.SalesQuotationResponse{}, 404, ErrQuotationNotFound
	}

	customerResponse := setting.SettingCustomerToResponse(customer)
	orderItemsResponse := make([]sales.SalesOrderItemResponse, 0)
	for _, item := range orderItems {
		orderItemsResponse = append(orderItemsResponse, sales.SalesOrderItemToResponse(item))
	}

	return sales.SalesQuotationToResponse(quotation, customerResponse, orderItemsResponse), 200, nil
}

func (service *SalesQuotationService) CreateQuotation(ctx *utils.CtxW, quotation *sales.SalesQuotationCreateRequest) (int, error) {
	commonModel := models.CommonModel{}
	commonModel.PrepareForCreate(ctx.User.Id, ctx.User.Id)

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	bqbQuery := bqb.New(`INSERT INTO "sales.quotation" 
	(name, creation_date, validity_date, discount, status, setting_customer_id, cid, ctime, mid, mtime)
	VALUES
	(?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`, quotation.Name, quotation.CreationDate, quotation.ValidityDate, quotation.Discount, quotation.Status, quotation.CustomerId, commonModel.CId, commonModel.CTime, commonModel.MId, commonModel.MTime)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("%v", err)
		tx.Rollback()
		return 500, utils.ErrInternalServer
	}

	var id int
	err = tx.QueryRow(query, params...).Scan(&id)
	if err != nil {
		switch err.(*pq.Error).Constraint {
		case database.KEY_SALES_QUOTATION_NAME:
			return 409, ErrQuotationNameExists
		case database.FK_SALES_QUOTATION_CUSTOMER_ID:
			return 400, ErrCustomerNotFound
		}

		log.Printf("%v", err)
		tx.Rollback()
		return 500, utils.ErrInternalServer
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
		log.Printf("%v", err)
		tx.Rollback()
		return 500, utils.ErrInternalServer
	}

	_, err = tx.Exec(query, params...)
	if err != nil {
		log.Printf("%v", err)
		tx.Rollback()
		return 500, utils.ErrInternalServer
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	return 201, nil
}

func (service *SalesQuotationService) UpdateQuotation(ctx *utils.CtxW, id string, quotation *sales.SalesQuotationUpdateRequest) (int, error) {
	commonModel := models.CommonModel{}
	commonModel.PrepareForCreate(ctx.User.Id, ctx.User.Id)

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	bqbQuery := bqb.New(`SELECT status FROM "sales.quotation" WHERE id = ?`, id)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("%v", err)
		tx.Rollback()
		return 500, utils.ErrInternalServer
	}

	var status string
	err = tx.QueryRow(query, params...).Scan(&status)
	if err != nil {
		tx.Rollback()

		if err == sql.ErrNoRows {
			return 404, ErrQuotationNotFound
		}

		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	if status == models.SalesQuotationStatusSalesOrder || status == models.SalesQuotationStatusSalesCancelled {
		tx.Rollback()
		return 400, errors.New("this quotation is not allowed to be updated")
	}

	bqbQuery = bqb.New(`UPDATE "sales.quotation" SET mid = ?, mtime = ?`, commonModel.MId, commonModel.MTime)
	utils.PrepareUpdateBqbQuery(bqbQuery, quotation)
	bqbQuery.Space(`WHERE id = ?`, id)

	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("%v", err)
		tx.Rollback()
		return 500, utils.ErrInternalServer
	}

	result, err := tx.Exec(query, params...)
	if err != nil {
		switch err.(*pq.Error).Constraint {
		case database.KEY_SALES_QUOTATION_NAME:
			return 409, ErrQuotationNameExists
		case database.FK_SALES_QUOTATION_CUSTOMER_ID:
			return 400, ErrCustomerNotFound
		}

		log.Printf("%v", err)
		tx.Rollback()
		return 500, utils.ErrInternalServer
	}

	if n, _ := result.RowsAffected(); n == 0 {
		tx.Rollback()
		return 404, ErrQuotationNotFound
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
				log.Printf("%v", err)
			}

			result, err = tx.Exec(query, params...)
			if err != nil {
				log.Printf("%v", err)
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
				utils.PrepareUpdateBqbQuery(bqbQuery, &item)
				bqbQuery.Space(`WHERE id = ? AND sales_quotation_id = ?`, item.Id, id)

				query, params, err := bqbQuery.ToPgsql()
				if err != nil {
					log.Printf("%v", err)
				}

				_, err = tx.Exec(query, params...)
				if err != nil {
					log.Printf("%v", err)
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
				log.Printf("%v", err)
			}

			_, err = tx.Exec(query, params...)
			if err != nil {
				log.Printf("%v", err)
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
				errorMessage = ErrQuotationNameExists
			case database.FK_SALES_QUOTATION_CUSTOMER_ID:
				errorMessage = ErrCustomerNotFound
			}
			if !hasError {
				hasError = true
				tx.Rollback()
			}
		}
	}()

	wg.Wait()
	close(errorChan)

	if hasError {
		switch errorMessage {
		case ErrQuotationNameExists:
			return 409, errorMessage
		case ErrCustomerNotFound:
			return 400, errorMessage
		}

		log.Printf("%v", errorMessage)
		return 500, errorMessage
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	return 200, nil
}
