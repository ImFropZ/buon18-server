package sales

import (
	"database/sql"
	"errors"
	"log"
	"server/database"
	"server/models"
	"server/models/accounting"
	"server/models/sales"
	"server/models/setting"
	"server/utils"
	"strings"

	"github.com/lib/pq"
	"github.com/nullism/bqb"
)

var (
	ErrOrderNotFound       = errors.New("order not found")
	ErrOrderNameExists     = errors.New("order name already exists")
	ErrPaymentTermNotFound = errors.New("payment term not found")
)

type SalesOrderService struct {
	DB *sql.DB
}

func (service *SalesOrderService) Orders(qp *utils.QueryParams) ([]sales.SalesOrderResponse, int, int, error) {
	bqbQuery := bqb.New(`WITH "limited_orders" AS (
		SELECT
			id,
			name,
			commitment_date,
			note,
			sales_quotation_id,
			accounting_payment_term_id
		FROM
			"sales.order"`)

	qp.FilterIntoBqb(bqbQuery)
	qp.PaginationIntoBqb(bqbQuery)

	bqbQuery.Space(`)
	SELECT
		"limited_orders".id,
		"limited_orders".name,
		"limited_orders".commitment_date,
		"limited_orders".note,
		"sales.quotation".id,
		"sales.quotation".name,
		"sales.quotation".creation_date,
		"sales.quotation".validity_date,
		"sales.quotation".discount,
		"sales.quotation".amount_delivery,
		"sales.quotation".status,
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
		"sales.order_item".discount,
		"accounting.payment_term".id,
		"accounting.payment_term".name,
		"accounting.payment_term".description,
		"accounting.payment_term_line".id,
		"accounting.payment_term_line".sequence,
		"accounting.payment_term_line".value_amount_percent,
		"accounting.payment_term_line".number_of_days
	FROM 
		"limited_orders"
	INNER JOIN "sales.quotation" ON "sales.quotation".id = "limited_orders".sales_quotation_id
	INNER JOIN "setting.customer" ON "setting.customer".id = "sales.quotation".setting_customer_id
	INNER JOIN "sales.order_item" ON "sales.order_item".sales_quotation_id = "sales.quotation".id
	INNER JOIN "accounting.payment_term" ON "accounting.payment_term".id = "limited_orders".accounting_payment_term_id
	INNER JOIN "accounting.payment_term_line" ON "accounting.payment_term_line".accounting_payment_term_id = "accounting.payment_term".id`)

	qp.OrderByIntoBqb(bqbQuery, `"limited_orders".id ASC, "sales.quotation".id ASC, "sales.order_item".id ASC, "accounting.payment_term_line".sequence ASC`)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("%v", err)
		return nil, 0, 500, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		log.Printf("%v", err)
		return nil, 0, 500, utils.ErrInternalServer
	}

	ordersResponse := make([]sales.SalesOrderResponse, 0)
	lastOrder := sales.SalesOrder{}
	lastCustomer := setting.SettingCustomer{}
	lastQuotation := sales.SalesQuotation{}
	orderItems := make([]sales.SalesOrderItem, 0)
	lastPaymentTerm := accounting.AccountingPaymentTerm{}
	paymentTermLines := make([]accounting.AccountingPaymentTermLine, 0)
	for rows.Next() {
		tmpOrder := sales.SalesOrder{}
		tmpCustomer := setting.SettingCustomer{}
		tmpQuotation := sales.SalesQuotation{}
		tmpOrderItem := sales.SalesOrderItem{}
		tmpPaymentTerm := accounting.AccountingPaymentTerm{}
		tmpPaymentTermLine := accounting.AccountingPaymentTermLine{}
		err := rows.Scan(
			&tmpOrder.Id,
			&tmpOrder.Name,
			&tmpOrder.CommitmentDate,
			&tmpOrder.Note,
			&tmpQuotation.Id,
			&tmpQuotation.Name,
			&tmpQuotation.CreationDate,
			&tmpQuotation.ValidityDate,
			&tmpQuotation.Discount,
			&tmpQuotation.AmountDelivery,
			&tmpQuotation.Status,
			&tmpCustomer.Id,
			&tmpCustomer.FullName,
			&tmpCustomer.Gender,
			&tmpCustomer.Email,
			&tmpCustomer.Phone,
			&tmpCustomer.AdditionalInformation,
			&tmpOrderItem.Id,
			&tmpOrderItem.Name,
			&tmpOrderItem.Description,
			&tmpOrderItem.Price,
			&tmpOrderItem.Discount,
			&tmpPaymentTerm.Id,
			&tmpPaymentTerm.Name,
			&tmpPaymentTerm.Description,
			&tmpPaymentTermLine.Id,
			&tmpPaymentTermLine.Sequence,
			&tmpPaymentTermLine.ValueAmountPercent,
			&tmpPaymentTermLine.NumberOfDays,
		)
		if err != nil {
			log.Printf("%v", err)
			return nil, 0, 500, utils.ErrInternalServer
		}

		if (lastQuotation.Id != tmpQuotation.Id && lastQuotation.Id != 0) && (lastPaymentTerm.Id != tmpPaymentTerm.Id && lastPaymentTerm.Id != 0) {
			orderItemsResponse := make([]sales.SalesOrderItemResponse, 0)
			for _, item := range orderItems {
				orderItemsResponse = append(orderItemsResponse, sales.SalesOrderItemToResponse(item))
			}
			customerResponse := setting.SettingCustomerToResponse(lastCustomer)
			quotationResponse := sales.SalesQuotationToResponse(lastQuotation, customerResponse, orderItemsResponse)
			paymentTermLinesResponse := make([]accounting.AccountingPaymentTermLineResponse, 0)
			for _, line := range paymentTermLines {
				paymentTermLinesResponse = append(paymentTermLinesResponse, accounting.AccountingPaymentTermLineToResponse(line))
			}
			paymentTermResponse := accounting.AccountingPaymentTermToResponse(lastPaymentTerm, paymentTermLinesResponse)
			ordersResponse = append(ordersResponse, sales.SalesOrderToResponse(lastOrder, quotationResponse, paymentTermResponse))

			// Reset and append new data
			lastOrder = tmpOrder
			lastQuotation = tmpQuotation
			lastCustomer = tmpCustomer
			orderItems = make([]sales.SalesOrderItem, 0)
			orderItems = append(orderItems, tmpOrderItem)
			lastPaymentTerm = tmpPaymentTerm
			paymentTermLines = make([]accounting.AccountingPaymentTermLine, 0)
			paymentTermLines = append(paymentTermLines, tmpPaymentTermLine)
			continue
		}

		if lastQuotation.Id == 0 || lastPaymentTerm.Id == 0 || lastOrder.Id == 0 {
			lastOrder = tmpOrder
			lastQuotation = tmpQuotation
			lastCustomer = tmpCustomer
			lastPaymentTerm = tmpPaymentTerm
		}

		if tmpOrderItem.Id != 0 {
			if len(orderItems) != 0 {
				if orderItems[len(orderItems)-1].Id != tmpOrderItem.Id {
					orderItems = append(orderItems, tmpOrderItem)
				}
			} else {
				orderItems = append(orderItems, tmpOrderItem)
			}
		}

		if tmpPaymentTermLine.Id != 0 {
			if len(paymentTermLines) != 0 {
				if paymentTermLines[len(paymentTermLines)-1].Id != tmpPaymentTermLine.Id {
					paymentTermLines = append(paymentTermLines, tmpPaymentTermLine)
				}
			} else {
				paymentTermLines = append(paymentTermLines, tmpPaymentTermLine)
			}
		}
	}
	if lastOrder.Id != 0 {
		orderItemsResponse := make([]sales.SalesOrderItemResponse, 0)
		for _, item := range orderItems {
			orderItemsResponse = append(orderItemsResponse, sales.SalesOrderItemToResponse(item))
		}
		customerResponse := setting.SettingCustomerToResponse(lastCustomer)
		quotationResponse := sales.SalesQuotationToResponse(lastQuotation, customerResponse, orderItemsResponse)
		paymentTermLinesResponse := make([]accounting.AccountingPaymentTermLineResponse, 0)
		for _, line := range paymentTermLines {
			paymentTermLinesResponse = append(paymentTermLinesResponse, accounting.AccountingPaymentTermLineToResponse(line))
		}
		paymentTermResponse := accounting.AccountingPaymentTermToResponse(lastPaymentTerm, paymentTermLinesResponse)
		ordersResponse = append(ordersResponse, sales.SalesOrderToResponse(lastOrder, quotationResponse, paymentTermResponse))
	}

	bqbQuery = bqb.New(`SELECT COUNT(*) FROM "sales.order"`)
	qp.FilterIntoBqb(bqbQuery)

	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("%v", err)
		return nil, 0, 500, utils.ErrInternalServer
	}

	var total int
	err = service.DB.QueryRow(query, params...).Scan(&total)
	if err != nil {
		log.Printf("%v", err)
		return nil, 0, 500, utils.ErrInternalServer
	}

	return ordersResponse, total, 200, nil
}

func (service *SalesOrderService) Order(id string) (sales.SalesOrderResponse, int, error) {
	bqbQuery := bqb.New(`WITH "limited_orders" AS (
		SELECT
			id,
			name,
			commitment_date,
			note,
			sales_quotation_id,
			accounting_payment_term_id
		FROM
			"sales.order"
		WHERE
			id = ?
	)
	SELECT
		"limited_orders".id,
		"limited_orders".name,
		"limited_orders".commitment_date,
		"limited_orders".note,
		"sales.quotation".id,
		"sales.quotation".name,
		"sales.quotation".creation_date,
		"sales.quotation".validity_date,
		"sales.quotation".discount,
		"sales.quotation".amount_delivery,
		"sales.quotation".status,
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
		"sales.order_item".discount,
		"accounting.payment_term".id,
		"accounting.payment_term".name,
		"accounting.payment_term".description,
		"accounting.payment_term_line".id,
		"accounting.payment_term_line".sequence,
		"accounting.payment_term_line".value_amount_percent,
		"accounting.payment_term_line".number_of_days
	FROM 
		"limited_orders"
	INNER JOIN "sales.quotation" ON "sales.quotation".id = "limited_orders".sales_quotation_id
	INNER JOIN "setting.customer" ON "setting.customer".id = "sales.quotation".setting_customer_id
	INNER JOIN "sales.order_item" ON "sales.order_item".sales_quotation_id = "sales.quotation".id
	INNER JOIN "accounting.payment_term" ON "accounting.payment_term".id = "limited_orders".accounting_payment_term_id
	INNER JOIN "accounting.payment_term_line" ON "accounting.payment_term_line".accounting_payment_term_id = "accounting.payment_term".id
	ORDER BY "limited_orders".id ASC, "sales.quotation".id ASC, "sales.order_item".id ASC, "accounting.payment_term_line".sequence ASC`, id)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("%v", err)
		return sales.SalesOrderResponse{}, 500, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		log.Printf("%v", err)
		return sales.SalesOrderResponse{}, 500, utils.ErrInternalServer
	}

	order := sales.SalesOrder{}
	customer := setting.SettingCustomer{}
	quotation := sales.SalesQuotation{}
	orderItems := make([]sales.SalesOrderItem, 0)
	paymentTerm := accounting.AccountingPaymentTerm{}
	paymentTermLines := make([]accounting.AccountingPaymentTermLine, 0)
	for rows.Next() {
		tmpOrderItem := sales.SalesOrderItem{}
		tmpPaymentTermLine := accounting.AccountingPaymentTermLine{}
		err := rows.Scan(
			&order.Id,
			&order.Name,
			&order.CommitmentDate,
			&order.Note,
			&quotation.Id,
			&quotation.Name,
			&quotation.CreationDate,
			&quotation.ValidityDate,
			&quotation.Discount,
			&quotation.AmountDelivery,
			&quotation.Status,
			&customer.Id,
			&customer.FullName,
			&customer.Gender,
			&customer.Email,
			&customer.Phone,
			&customer.AdditionalInformation,
			&tmpOrderItem.Id,
			&tmpOrderItem.Name,
			&tmpOrderItem.Description,
			&tmpOrderItem.Price,
			&tmpOrderItem.Discount,
			&paymentTerm.Id,
			&paymentTerm.Name,
			&paymentTerm.Description,
			&tmpPaymentTermLine.Id,
			&tmpPaymentTermLine.Sequence,
			&tmpPaymentTermLine.ValueAmountPercent,
			&tmpPaymentTermLine.NumberOfDays,
		)
		if err != nil {
			log.Printf("%v", err)
			return sales.SalesOrderResponse{}, 500, utils.ErrInternalServer
		}

		if tmpOrderItem.Id != 0 {
			if len(orderItems) != 0 {
				if orderItems[len(orderItems)-1].Id != tmpOrderItem.Id {
					orderItems = append(orderItems, tmpOrderItem)
				}
			} else {
				orderItems = append(orderItems, tmpOrderItem)
			}
		}

		if tmpPaymentTermLine.Id != 0 {
			if len(paymentTermLines) != 0 {
				if paymentTermLines[len(paymentTermLines)-1].Id != tmpPaymentTermLine.Id {
					paymentTermLines = append(paymentTermLines, tmpPaymentTermLine)
				}
			} else {
				paymentTermLines = append(paymentTermLines, tmpPaymentTermLine)
			}
		}
	}
	if order.Id == 0 {
		return sales.SalesOrderResponse{}, 404, ErrOrderNotFound
	}

	orderItemsResponse := make([]sales.SalesOrderItemResponse, 0)
	for _, item := range orderItems {
		orderItemsResponse = append(orderItemsResponse, sales.SalesOrderItemToResponse(item))
	}
	customerResponse := setting.SettingCustomerToResponse(customer)
	quotationResponse := sales.SalesQuotationToResponse(quotation, customerResponse, orderItemsResponse)
	paymentTermLinesResponse := make([]accounting.AccountingPaymentTermLineResponse, 0)
	for _, line := range paymentTermLines {
		paymentTermLinesResponse = append(paymentTermLinesResponse, accounting.AccountingPaymentTermLineToResponse(line))
	}
	paymentTermResponse := accounting.AccountingPaymentTermToResponse(paymentTerm, paymentTermLinesResponse)

	return sales.SalesOrderToResponse(order, quotationResponse, paymentTermResponse), 200, nil
}

func (service *SalesOrderService) CreateOrder(ctx *utils.CtxW, order *sales.SalesOrderCreateRequest) (int, error) {
	commonModel := models.CommonModel{}
	commonModel.PrepareForCreate(ctx.User.Id, ctx.User.Id)

	bqbQuery := bqb.New(`CALL create_sales_order(?, ?, ?, ?, ?, ?, ?, ?, ?)`, order.Name, order.CommitmentDate, order.Note, order.QuotationId, order.PaymentTermId, commonModel.CId, commonModel.CTime, commonModel.MId, commonModel.MTime)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	_, err = service.DB.Exec(query, params...)
	if err != nil {
		if message := err.(*pq.Error).Message; strings.HasPrefix(message, "custom_error:") {
			return 400, errors.New(strings.TrimPrefix(message, "custom_error:"))
		}

		switch err.(*pq.Error).Constraint {
		case database.KEY_SALES_ORDER_NAME:
			return 409, ErrOrderNameExists
		case database.FK_ACCOUNTING_PAYMENT_TERM_ID:
			return 400, ErrPaymentTermNotFound
		}

		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	return 201, nil
}
