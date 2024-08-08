package setting

import (
	"database/sql"
	"errors"
	"log"
	"server/models"
	"server/utils"

	"github.com/nullism/bqb"
)

var (
	ErrCustomerNotFound = errors.New("customer not found")
)

type SettingCustomerService struct {
	DB *sql.DB
}

func (service *SettingCustomerService) Customers(qp *utils.QueryParams) ([]models.SettingCustomerResponse, int, int, error) {
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

	customersResponse := make([]models.SettingCustomerResponse, 0)
	for rows.Next() {
		tmpCustomer := models.SettingCustomer{}
		err := rows.Scan(&tmpCustomer.Id, &tmpCustomer.FullName, &tmpCustomer.Gender, &tmpCustomer.Email, &tmpCustomer.Phone, &tmpCustomer.AdditionalInformation)
		if err != nil {
			log.Printf("%s", err)
			return nil, 0, 500, utils.ErrInternalServer
		}

		customersResponse = append(customersResponse, models.SettingCustomerToResponse(tmpCustomer))
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

func (service *SettingCustomerService) Customer(id string) (models.SettingCustomerResponse, int, error) {
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
		return models.SettingCustomerResponse{}, 500, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		log.Printf("%s", err)
		return models.SettingCustomerResponse{}, 500, utils.ErrInternalServer
	}

	customerResponse := models.SettingCustomerResponse{}
	for rows.Next() {
		tmpCustomer := models.SettingCustomer{}
		err := rows.Scan(&tmpCustomer.Id, &tmpCustomer.FullName, &tmpCustomer.Gender, &tmpCustomer.Email, &tmpCustomer.Phone, &tmpCustomer.AdditionalInformation)
		if err != nil {
			log.Printf("%s", err)
			return models.SettingCustomerResponse{}, 500, utils.ErrInternalServer
		}

		customerResponse = models.SettingCustomerToResponse(tmpCustomer)
	}

	if customerResponse.Id == 0 {
		return models.SettingCustomerResponse{}, 404, ErrCustomerNotFound
	}

	return customerResponse, 0, nil
}
