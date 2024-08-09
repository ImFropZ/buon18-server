package sales

import (
	"database/sql"
	"errors"
	"log"
	"server/models/sales"
	"server/models/setting"
	"server/utils"

	"github.com/nullism/bqb"
)

var (
	ErrQuotationNotFound = errors.New("quotation not found")
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
			quotationsResponse = append(quotationsResponse, sales.SalesQuotationToResponse(lastQuotation, lastCustomer, orderItems))
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
		quotationsResponse = append(quotationsResponse, sales.SalesQuotationToResponse(lastQuotation, lastCustomer, orderItems))
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

	return sales.SalesQuotationToResponse(quotation, customer, orderItems), 200, nil
}
