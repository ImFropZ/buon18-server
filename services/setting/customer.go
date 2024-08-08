package setting

import (
	"database/sql"
	"fmt"
	"log"
	"server/models"
	"server/utils"

	"github.com/nullism/bqb"
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

	if len(qp.Fitlers) > 0 {
		bqbQuery.Space("WHERE")
		for index, filter := range qp.Fitlers {
			bqbQuery.Space(fmt.Sprintf("%s %s ?", filter.Field, utils.MAPPED_FILTER_OPERATORS_TO_SQL[filter.Operator]), filter.Value)
			if index < len(qp.Fitlers)-1 {
				bqbQuery.Space("AND")
			}
		}
	}

	if len(qp.OrderBy) > 0 {
		bqbQuery.Space("ORDER BY")
		for index, sort := range qp.OrderBy {
			bqbQuery.Space(sort)
			if index < len(qp.Fitlers)-1 {
				bqbQuery.Space(",")
			}
		}
	} else {
		bqbQuery.Space(`ORDER BY "setting.customer".id ASC`)
	}

	bqbQuery.Space(`OFFSET ? LIMIT ?`, qp.Pagination.Offset, qp.Pagination.Limit)

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
	if len(qp.Fitlers) > 0 {
		bqbQuery.Space("WHERE")
		for index, filter := range qp.Fitlers {
			bqbQuery.Space(fmt.Sprintf("%s %s ?", filter.Field, utils.MAPPED_FILTER_OPERATORS_TO_SQL[filter.Operator]), filter.Value)
			if index < len(qp.Fitlers)-1 {
				bqbQuery.Space("AND")
			}
		}
	}

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
