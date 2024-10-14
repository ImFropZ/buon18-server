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

type AccountingJournalService struct {
	DB *sql.DB
}

func (service *AccountingJournalService) Journals(qp *utils.QueryParams) ([]accounting.AccountingJournalResponse, int, int, error) {
	bqbQuery := bqb.New(`WITH "limited_journals" AS (
		SELECT
			*
		FROM
			"accounting.journal"`)

	qp.FilterIntoBqb(bqbQuery)
	qp.PaginationIntoBqb(bqbQuery)

	bqbQuery.Space(`)
	SELECT
		"limited_journals".id,
		"limited_journals".code,
		"limited_journals".name,
		"limited_journals".typ,
		"accounting.account".id,
		"accounting.account".code,
		"accounting.account".name,
		"accounting.account".typ
	FROM
		"limited_journals"
	INNER JOIN "accounting.account" ON "accounting.account".id = "limited_journals".accounting_account_id`)

	qp.OrderByIntoBqb(bqbQuery, `"limited_journals".id ASC, "limited_journals".code ASC`)

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

	var journals []accounting.AccountingJournalResponse
	for rows.Next() {
		var journal accounting.AccountingJournal
		var account accounting.AccountingAccount
		err := rows.Scan(&journal.Id, &journal.Code, &journal.Name, &journal.Typ, &account.Id, &account.Code, &account.Name, &account.Typ)
		if err != nil {
			slog.Error(fmt.Sprintf("%v", err))
			return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
		}

		accountResponse := accounting.AccountingAccountToResponse(account)
		journals = append(journals, accounting.AccountingJournalToResponse(journal, &accountResponse))
	}

	bqbQuery = bqb.New(`SELECT COUNT(*) FROM "accounting.journal"`)
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

	return journals, total, http.StatusOK, nil
}

func (service *AccountingJournalService) Journal(id string) (accounting.AccountingJournalResponse, int, error) {
	bqbQuery := bqb.New(`WITH "limited_journals" AS (
		SELECT
			*
		FROM
			"accounting.journal"
		WHERE
			"accounting.journal".id = ?
	)
	SELECT
		"limited_journals".id,
		"limited_journals".code,
		"limited_journals".name,
		"limited_journals".typ,
		"accounting.account".id,
		"accounting.account".code,
		"accounting.account".name,
		"accounting.account".typ
	FROM
		"limited_journals"
	INNER JOIN "accounting.account" ON "accounting.account".id = "limited_journals".accounting_account_id`, id)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return accounting.AccountingJournalResponse{}, http.StatusInternalServerError, utils.ErrInternalServer
	}

	var journal accounting.AccountingJournal
	var account accounting.AccountingAccount
	if err := service.DB.QueryRow(query, params...).Scan(&journal.Id, &journal.Code, &journal.Name, &journal.Typ, &account.Id, &account.Code, &account.Name, &account.Typ); err != nil {
		if err == sql.ErrNoRows {
			return accounting.AccountingJournalResponse{}, http.StatusNotFound, utils.ErrJournalNotFound
		}

		slog.Error(fmt.Sprintf("%v", err))
		return accounting.AccountingJournalResponse{}, http.StatusInternalServerError, utils.ErrInternalServer
	}

	if journal.Id == nil {
		return accounting.AccountingJournalResponse{}, http.StatusNotFound, utils.ErrJournalNotFound
	}

	accountResponse := accounting.AccountingAccountToResponse(account)
	return accounting.AccountingJournalToResponse(journal, &accountResponse), http.StatusOK, nil
}

func (service *AccountingJournalService) CreateJournal(ctx *utils.CtxValue, journal *accounting.AccountingJournalCreateRequest) (int, error) {
	commonModel := models.CommonModel{}
	commonModel.PrepareForCreate(*ctx.User.Id, *ctx.User.Id)

	bqbQuery := bqb.New(`INSERT INTO "accounting.journal"
	(code, name, typ, accounting_account_id, cid, ctime, mid, mtime)
	VALUES
	(?, ?, ?, ?, ?, ?, ?, ?)`, journal.Code, journal.Name, journal.Typ, journal.AccountId, commonModel.CId, commonModel.CTime, commonModel.MId, commonModel.MTime)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if _, err := service.DB.Exec(query, params...); err != nil {
		switch err.(*pq.Error).Constraint {
		case database.FK_ACCOUNTING_ACCOUNT_ID:
			return http.StatusNotFound, utils.ErrAccountNotFound
		case database.KEY_ACCOUNTING_JOURNAL_CODE:
			return http.StatusConflict, utils.ErrJournalCodeExists
		}

		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	return http.StatusCreated, nil
}

func (service *AccountingJournalService) UpdateJournal(ctx *utils.CtxValue, id string, journal *accounting.AccountingJournalUpdateRequest) (int, error) {
	commonModel := models.CommonModel{}
	commonModel.PrepareForUpdate(*ctx.User.Id)

	bqbQuery := bqb.New(`UPDATE "accounting.journal" SET mid = ?, mtime = ?`, commonModel.MId, commonModel.MTime)
	if journal.Code != nil {
		bqbQuery.Comma(`code = ?`, *journal.Code)
	}
	if journal.Name != nil {
		bqbQuery.Comma(`name = ?`, *journal.Name)
	}
	if journal.Typ != nil {
		bqbQuery.Comma(`typ = ?`, *journal.Typ)
	}
	if journal.AccountId != nil {
		bqbQuery.Comma(`accounting_account_id = ?`, *journal.AccountId)
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
		case database.FK_ACCOUNTING_ACCOUNT_ID:
			return http.StatusNotFound, utils.ErrAccountNotFound
		case database.KEY_ACCOUNTING_JOURNAL_CODE:
			return http.StatusConflict, utils.ErrJournalCodeExists
		}

		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if n, _ := result.RowsAffected(); n == 0 {
		return http.StatusNotFound, utils.ErrJournalNotFound
	}

	return http.StatusOK, nil
}

func (service *AccountingJournalService) DeleteJournals(req *models.CommonDelete) (int, error) {
	bqbQuery := bqb.New(`DELETE FROM "accounting.journal" WHERE id in (`)
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
		case database.FK_ACCOUNTING_JOURNAL_ID:
			return http.StatusForbidden, utils.ErrResourceInUsed
		}

		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if n, _ := result.RowsAffected(); n == 0 {
		return http.StatusNotFound, utils.ErrJournalNotFound
	}

	return http.StatusOK, nil
}
