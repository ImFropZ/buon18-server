package accounting

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/lib/pq"
	"github.com/nullism/bqb"
	"system.buon18.com/m/database"
	"system.buon18.com/m/models"
	"system.buon18.com/m/models/accounting"
	"system.buon18.com/m/utils"
)

type AccountingAccountService struct {
	DB *sql.DB
}

func (service *AccountingAccountService) Accounts(qp *utils.QueryParams) ([]accounting.AccountingAccountResponse, int, int, error) {
	bqbQuery := bqb.New(`SELECT
		"accounting.account".id,
		"accounting.account".code,
		"accounting.account".name,
		"accounting.account".typ
	FROM
		"accounting.account"`)

	qp.FilterIntoBqb(bqbQuery)
	qp.OrderByIntoBqb(bqbQuery, `"accounting.account".id ASC`)
	qp.PaginationIntoBqb(bqbQuery)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
	}

	accounts := []accounting.AccountingAccountResponse{}
	for rows.Next() {
		account := accounting.AccountingAccount{}
		if err := rows.Scan(&account.Id, &account.Code, &account.Name, &account.Typ); err != nil {
			slog.Error(fmt.Sprintf("%v", err))
			return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
		}

		accounts = append(accounts, accounting.AccountingAccountToResponse(account))
	}

	bqbQuery = bqb.New(`SELECT COUNT(*) FROM "accounting.account"`)
	qp.FilterIntoBqb(bqbQuery)

	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
	}

	var total int
	if err := service.DB.QueryRow(query, params...).Scan(&total); err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
	}

	return accounts, total, http.StatusOK, nil
}

func (service *AccountingAccountService) Account(id string) (accounting.AccountingAccountResponse, int, error) {
	bqbQuery := bqb.New(`SELECT
		"accounting.account".id,
		"accounting.account".code,
		"accounting.account".name,
		"accounting.account".typ
	FROM
		"accounting.account"
	WHERE
		"accounting.account".id = ?`, id)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return accounting.AccountingAccountResponse{}, http.StatusInternalServerError, utils.ErrInternalServer
	}

	account := accounting.AccountingAccount{}
	if err = service.DB.QueryRow(query, params...).Scan(&account.Id, &account.Code, &account.Name, &account.Typ); err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return accounting.AccountingAccountResponse{}, http.StatusNotFound, utils.ErrAccountNotFound
	}

	return accounting.AccountingAccountToResponse(account), http.StatusOK, nil
}

func (service *AccountingAccountService) CreateAccount(ctx *utils.CtxValue, account *accounting.AccountingAccountCreateRequest) (int, error) {
	commonModel := models.CommonModel{}
	commonModel.PrepareForCreate(*ctx.User.Id, *ctx.User.Id)

	bqbQuery := bqb.New(`INSERT INTO "accounting.account"
	(name, code, typ, cid, ctime, mid, mtime)
	VALUES
	(?, ?, ?, ?, ?, ?, ?)`, account.Name, account.Code, account.Typ, commonModel.CId, commonModel.CTime, commonModel.MId, commonModel.MTime)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if _, err = service.DB.Exec(query, params...); err != nil {
		switch err.(*pq.Error).Constraint {
		case database.KEY_ACCOUNTING_ACCOUNT_CODE:
			return http.StatusConflict, utils.ErrResourceInUsed
		}

		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	return http.StatusCreated, nil
}

func (service *AccountingAccountService) UpdateAccount(ctx *utils.CtxValue, id string, account *accounting.AccountingAccountUpdateRequest) (int, error) {
	commonModel := models.CommonModel{}
	commonModel.PrepareForUpdate(*ctx.User.Id)

	bqbQuery := bqb.New(`UPDATE "accounting.account" SET mid = ?, mtime = ?`, commonModel.MId, commonModel.MTime)
	if account.Name != nil {
		bqbQuery.Space(`SET name = ?`, *account.Name)
	}
	if account.Code != nil {
		bqbQuery.Space(`SET code = ?`, *account.Code)
	}
	if account.Typ != nil {
		bqbQuery.Space(`SET typ = ?`, *account.Typ)
	}
	bqbQuery.Space(`WHERE id = ?`, id)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	result, err := service.DB.Exec(query, params...)
	if err != nil {
		switch err.(*pq.Error).Constraint {
		case database.KEY_ACCOUNTING_ACCOUNT_CODE:
			return http.StatusConflict, utils.ErrAccountCodeExists
		}

		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if n, _ := result.RowsAffected(); n == 0 {
		return http.StatusNotFound, utils.ErrAccountNotFound
	}

	return http.StatusOK, nil
}

func (service *AccountingAccountService) DeleteAccounts(req *models.CommonDelete) (int, error) {
	bqbQuery := bqb.New(`DELETE FROM "accounting.account" WHERE id in (`)
	for i, id := range req.Ids {
		bqbQuery.Space(`?`, id)

		if i < len(req.Ids)-1 {
			bqbQuery.Comma("")
		}
	}
	bqbQuery.Space(`)`)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	result, err := service.DB.Exec(query, params...)
	if err != nil {
		switch err.(*pq.Error).Constraint {
		case database.FK_ACCOUNTING_ACCOUNT_ID:
			return http.StatusConflict, utils.ErrResourceInUsed
		}

		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if n, _ := result.RowsAffected(); n == 0 {
		return http.StatusNotFound, utils.ErrAccountNotFound
	}

	return http.StatusOK, nil
}
