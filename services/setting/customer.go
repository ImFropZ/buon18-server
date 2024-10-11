package setting

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"system.buon18.com/m/database"
	"system.buon18.com/m/models"
	"system.buon18.com/m/models/setting"
	"system.buon18.com/m/utils"

	"github.com/lib/pq"
	"github.com/nullism/bqb"
)

type SettingCustomerService struct {
	DB *sql.DB
}

func (service *SettingCustomerService) Customers(qp *utils.QueryParams) ([]setting.SettingCustomerResponse, int, int, error) {
	bqbQuery := bqb.New(`
	SELECT
		"setting.customer".id,
		"setting.customer".full_name,
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
		slog.Error(fmt.Sprintf("%s", err))
		return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		slog.Error(fmt.Sprintf("%s", err))
		return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
	}

	customersResponse := make([]setting.SettingCustomerResponse, 0)
	for rows.Next() {
		tmpCustomer := setting.SettingCustomer{}
		err := rows.Scan(&tmpCustomer.Id, &tmpCustomer.FullName, &tmpCustomer.Gender, &tmpCustomer.Email, &tmpCustomer.Phone, &tmpCustomer.AdditionalInformation)
		if err != nil {
			slog.Error(fmt.Sprintf("%s", err))
			return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
		}

		customersResponse = append(customersResponse, setting.SettingCustomerToResponse(tmpCustomer))
	}

	bqbQuery = bqb.New(`SELECT COUNT(*) FROM "setting.customer"`)
	qp.FilterIntoBqb(bqbQuery)

	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%s", err))
		return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
	}

	var total int
	if err := service.DB.QueryRow(query, params...).Scan(&total); err != nil {
		slog.Error(fmt.Sprintf("%s", err))
		return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
	}

	return customersResponse, total, http.StatusOK, nil
}

func (service *SettingCustomerService) Customer(id string) (setting.SettingCustomerResponse, int, error) {
	bqbQuery := bqb.New(`
	SELECT
		"setting.customer".id,
		"setting.customer".full_name,
		"setting.customer".gender,
		"setting.customer".email,
		"setting.customer".phone,
		"setting.customer".additional_information
	FROM "setting.customer" WHERE "setting.customer".id = ?`, id)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%s", err))
		return setting.SettingCustomerResponse{}, http.StatusInternalServerError, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		slog.Error(fmt.Sprintf("%s", err))
		return setting.SettingCustomerResponse{}, http.StatusInternalServerError, utils.ErrInternalServer
	}

	var customer setting.SettingCustomer
	for rows.Next() {
		err := rows.Scan(&customer.Id, &customer.FullName, &customer.Gender, &customer.Email, &customer.Phone, &customer.AdditionalInformation)
		if err != nil {
			slog.Error(fmt.Sprintf("%s", err))
			return setting.SettingCustomerResponse{}, http.StatusInternalServerError, utils.ErrInternalServer
		}
	}

	if customer.Id == nil {
		return setting.SettingCustomerResponse{}, http.StatusNotFound, utils.ErrCustomerNotFound
	}

	return setting.SettingCustomerToResponse(customer), http.StatusOK, nil
}

func (service *SettingCustomerService) CreateCustomer(ctx *utils.CtxValue, customer *setting.SettingCustomerCreateRequest) (int, error) {
	commonModel := models.CommonModel{}
	commonModel.PrepareForCreate(*ctx.User.Id, *ctx.User.Id)

	bqbQuery := bqb.New(`INSERT INTO "setting.customer" (
		full_name,
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
		slog.Error(fmt.Sprintf("%s", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if _, err := service.DB.Exec(query, params...); err != nil {
		switch err.(*pq.Error).Constraint {
		case database.KEY_SETTING_CUSTOMER_EMAIL:
			return http.StatusConflict, utils.ErrCustomerEmailExists
		}

		slog.Error(fmt.Sprintf("%s", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	return http.StatusCreated, nil
}

func (service *SettingCustomerService) UpdateCustomer(ctx *utils.CtxValue, id string, customer *setting.SettingCustomerUpdateRequest) (int, error) {
	commonModel := models.CommonModel{}
	commonModel.PrepareForUpdate(*ctx.User.Id)

	bqbQuery := bqb.New(`UPDATE "setting.customer" SET mid = ?, mtime = ?`, commonModel.MId, commonModel.MTime)
	if customer.Email != nil {
		bqbQuery.Space(`SET email = ?`, *customer.Email)
	}
	if customer.Phone != nil {
		bqbQuery.Space(`SET phone = ?`, *customer.Phone)
	}
	if customer.AdditionalInformation != nil {
		bqbQuery.Space(`SET additional_information = ?`, *customer.AdditionalInformation)
	}
	bqbQuery.Space(` WHERE id = ?`, id)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	result, err := service.DB.Exec(query, params...)
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if n, err := result.RowsAffected(); err != nil || n == 0 {
		if n == 0 {
			return http.StatusNotFound, utils.ErrCustomerNotFound
		}

		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	return http.StatusOK, nil
}

func (service *SettingCustomerService) DeleteCustomer(id string) (int, error) {
	bqbQuery := bqb.New(`DELETE FROM "setting.customer" WHERE id = ?`, id)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	result, err := service.DB.Exec(query, params...)
	if err != nil {
		switch err.(*pq.Error).Constraint {
		case database.FK_SETTING_CUSTOMER_ID:
			return http.StatusConflict, utils.ErrUnableToDeleteCustomer
		}

		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if n, err := result.RowsAffected(); err != nil || n == 0 {
		if n == 0 {
			return http.StatusNotFound, utils.ErrCustomerNotFound
		}

		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	return http.StatusOK, nil
}
