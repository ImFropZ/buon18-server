package setting

import (
	"database/sql"
	"errors"
	"log"
	"server/database"
	"server/models"
	"server/models/setting"
	"server/utils"

	"github.com/lib/pq"
	"github.com/nullism/bqb"
)

var (
	ErrCustomerNotFound    = errors.New("customer not found")
	ErrCustomerEmailExists = errors.New("customer email already exists")
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

	return customersResponse, total, 200, nil
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

	return setting.SettingCustomerToResponse(customer), 200, nil
}

func (service *SettingCustomerService) CreateCustomer(ctx *utils.CtxW, customer *setting.SettingCustomerCreateRequest) (int, error) {
	commonModel := models.CommonModel{}
	commonModel.PrepareForCreate(ctx.User.Id, ctx.User.Id)

	bqbQuery := bqb.New(`INSERT INTO "setting.customer" (
		fullname,
		gender,
		email,
		phone,
		additional_information,
		cid,
		ctime,
		mid,
		mtime
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, customer.FullName, customer.Gender, customer.Email, customer.Phone, customer.AdditionalInformation, commonModel.CId, commonModel.CTime, commonModel.MId, commonModel.MTime)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("%s", err)
		return 500, utils.ErrInternalServer
	}

	_, err = service.DB.Exec(query, params...)
	if err != nil {
		switch err.(*pq.Error).Constraint {
		case database.KEY_SETTING_CUSTOMER_EMAIL:
			return 409, ErrCustomerEmailExists
		}

		log.Printf("%s", err)
		return 500, utils.ErrInternalServer
	}

	return 201, nil
}

func (service *SettingCustomerService) UpdateCustomer(ctx *utils.CtxW, id string, customer *setting.SettingCustomerUpdateRequest) (int, error) {
	commonModel := models.CommonModel{}
	commonModel.PrepareForUpdate(ctx.User.Id)

	bqbQuery := bqb.New(`UPDATE "setting.customer" SET mid = ?, mtime = ?`, commonModel.MId, commonModel.MTime)
	utils.PrepareUpdateBqbQuery(bqbQuery, customer)
	bqbQuery.Space(` WHERE id = ?`, id)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	result, err := service.DB.Exec(query, params...)
	if err != nil {
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	if n, err := result.RowsAffected(); err != nil || n == 0 {
		if n == 0 {
			return 404, ErrCustomerNotFound
		}
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	return 200, nil
}
