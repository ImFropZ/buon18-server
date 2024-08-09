package setting

import (
	"database/sql"
	"errors"
	"log"
	"server/models/setting"
	"server/utils"

	"github.com/nullism/bqb"
)

var (
	ErrCustomerNotFound = errors.New("customer not found")
)

type SettingCustomerService struct {
	DB *sql.DB
}

func (service *SettingCustomerService) Customers(qp *utils.QueryParams) ([]setting.SettingCustomerResponse, int, int, error) {
	bqbQuery := bqb.New(`
	SELECT
		"setting.customer".id,
		"setting.customer".fullname,
		"setting.customer".gender,
		"setting.customer".email,
		"setting.customer".phone,
		"setting.customer".additional_information
	FROM "setting.customer"`)

	qp.FilterIntoBqb(bqbQuery)
	qp.OrderByIntoBqb(bqbQuery, `"setting.customer".id ASC`)
	qp.PaginationIntoBqb(bqbQuery)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("%s", err)
		return nil, 0, 500, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		log.Printf("%s", err)
		return nil, 0, 500, utils.ErrInternalServer
	}

	customersResponse := make([]setting.SettingCustomerResponse, 0)
	for rows.Next() {
		tmpCustomer := setting.SettingCustomer{}
		err := rows.Scan(&tmpCustomer.Id, &tmpCustomer.FullName, &tmpCustomer.Gender, &tmpCustomer.Email, &tmpCustomer.Phone, &tmpCustomer.AdditionalInformation)
		if err != nil {
			log.Printf("%s", err)
			return nil, 0, 500, utils.ErrInternalServer
		}

		customersResponse = append(customersResponse, setting.SettingCustomerToResponse(tmpCustomer))
	}

	bqbQuery = bqb.New(`SELECT COUNT(*) FROM "setting.customer"`)
	qp.FilterIntoBqb(bqbQuery)

	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("%s", err)
		return nil, 0, 500, utils.ErrInternalServer
	}

	var total int
	err = service.DB.QueryRow(query, params...).Scan(&total)
	if err != nil {
		log.Printf("%s", err)
		return nil, 0, 500, utils.ErrInternalServer
	}

	return customersResponse, total, 0, nil
}

func (service *SettingCustomerService) Customer(id string) (setting.SettingCustomerResponse, int, error) {
	bqbQuery := bqb.New(`
	SELECT
		"setting.customer".id,
		"setting.customer".fullname,
		"setting.customer".gender,
		"setting.customer".email,
		"setting.customer".phone,
		"setting.customer".additional_information
	FROM "setting.customer" WHERE "setting.customer".id = ?`, id)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("%s", err)
		return setting.SettingCustomerResponse{}, 500, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		log.Printf("%s", err)
		return setting.SettingCustomerResponse{}, 500, utils.ErrInternalServer
	}

	var customer setting.SettingCustomer
	for rows.Next() {
		err := rows.Scan(&customer.Id, &customer.FullName, &customer.Gender, &customer.Email, &customer.Phone, &customer.AdditionalInformation)
		if err != nil {
			log.Printf("%s", err)
			return setting.SettingCustomerResponse{}, 500, utils.ErrInternalServer
		}
	}

	if customer.Id == 0 {
		return setting.SettingCustomerResponse{}, 404, ErrCustomerNotFound
	}

	return setting.SettingCustomerToResponse(customer), 0, nil
}
