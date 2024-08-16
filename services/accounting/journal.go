package accounting

import (
	"database/sql"
	"errors"
	"log"
	"server/database"
	"server/models"
	"server/models/accounting"
	"server/utils"

	"github.com/lib/pq"
	"github.com/nullism/bqb"
)

var (
	ErrJournalNotFound                    = errors.New("journal not found")
	ErrAccountingJournalCodeExists        = errors.New("accounting journal code already exists")
	ErrUnableToDeleteCurrentlyUsedJournal = errors.New("unable to delete currently used journal")
)

type AccountingJournalService struct {
	DB *sql.DB
}

func (service *AccountingJournalService) Journals(qp *utils.QueryParams) ([]accounting.AccountingJournalResponse, int, int, error) {
	bqbQuery := bqb.New(`WITH "limited_journals" AS (
		SELECT
			"accounting.journal".id,
			"accounting.journal".code,
			"accounting.journal".name,
			"accounting.journal".typ,
			"accounting.journal".accounting_account_id
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
		log.Printf("%v", err)
		return nil, 0, 500, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		log.Printf("%v", err)
		return nil, 0, 500, utils.ErrInternalServer
	}

	var journals []accounting.AccountingJournalResponse
	for rows.Next() {
		var journal accounting.AccountingJournal
		var account accounting.AccountingAccount
		err := rows.Scan(&journal.Id, &journal.Code, &journal.Name, &journal.Typ, &account.Id, &account.Code, &account.Name, &account.Typ)
		if err != nil {
			log.Printf("%v", err)
			return nil, 0, 500, utils.ErrInternalServer
		}

		accountResponse := accounting.AccountingAccountToResponse(account)
		journals = append(journals, accounting.AccountingJournalToResponse(journal, &accountResponse))
	}

	bqbQuery = bqb.New(`SELECT COUNT(*) FROM "accounting.journal"`)
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

	return journals, total, 200, nil
}

func (service *AccountingJournalService) Journal(id string) (accounting.AccountingJournalResponse, int, error) {
	bqbQuery := bqb.New(`WITH "limited_journals" AS (
		SELECT
			"accounting.journal".id,
			"accounting.journal".code,
			"accounting.journal".name,
			"accounting.journal".typ,
			"accounting.journal".accounting_account_id
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
		log.Printf("%v", err)
		return accounting.AccountingJournalResponse{}, 500, utils.ErrInternalServer
	}

	var journal accounting.AccountingJournal
	var account accounting.AccountingAccount
	err = service.DB.QueryRow(query, params...).Scan(&journal.Id, &journal.Code, &journal.Name, &journal.Typ, &account.Id, &account.Code, &account.Name, &account.Typ)
	if err != nil {
		if err == sql.ErrNoRows {
			return accounting.AccountingJournalResponse{}, 404, ErrJournalNotFound
		}

		log.Printf("%v", err)
		return accounting.AccountingJournalResponse{}, 500, utils.ErrInternalServer
	}

	if journal.Id == 0 {
		return accounting.AccountingJournalResponse{}, 404, ErrJournalNotFound
	}

	accountResponse := accounting.AccountingAccountToResponse(account)

	return accounting.AccountingJournalToResponse(journal, &accountResponse), 200, nil
}

func (service *AccountingJournalService) CreateJournal(ctx *utils.CtxW, journal *accounting.AccountingJournalCreateRequest) (int, error) {
	commonModel := models.CommonModel{}
	commonModel.PrepareForCreate(ctx.User.Id, ctx.User.Id)

	bqbQuery := bqb.New(`INSERT INTO "accounting.journal"
	(code, name, typ, accounting_account_id, cid, ctime, mid, mtime)
	VALUES
	(?, ?, ?, ?, ?, ?, ?, ?)`, journal.Code, journal.Name, journal.Typ, journal.AccountId, commonModel.CId, commonModel.CTime, commonModel.MId, commonModel.MTime)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	_, err = service.DB.Exec(query, params...)
	if err != nil {
		switch err.(*pq.Error).Constraint {
		case database.FK_ACCOUNTING_ACCOUNT_ID:
			return 400, ErrAccountNotFound
		case database.KEY_ACCOUNTING_JOURNAL_CODE:
			return 409, ErrAccountingJournalCodeExists
		}

		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	return 201, nil
}

func (service *AccountingJournalService) UpdateJournal(ctx *utils.CtxW, id string, journal *accounting.AccountingJournalUpdateRequest) (int, error) {
	commonModel := models.CommonModel{}
	commonModel.PrepareForUpdate(ctx.User.Id)

	bqbQuery := bqb.New(`UPDATE "accounting.journal" SET mid = ?, mtime = ?`, commonModel.MId, commonModel.MTime)
	utils.PrepareUpdateBqbQuery(bqbQuery, journal)
	bqbQuery.Space(`WHERE id = ?`, id)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	result, err := service.DB.Exec(query, params...)
	if err != nil {
		switch err.(*pq.Error).Constraint {
		case database.FK_ACCOUNTING_ACCOUNT_ID:
			return 400, ErrAccountNotFound
		case database.KEY_ACCOUNTING_JOURNAL_CODE:
			return 409, ErrAccountingJournalCodeExists
		}

		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	if n, _ := result.RowsAffected(); n == 0 {
		return 404, ErrJournalNotFound
	}

	return 200, nil
}

func (service *AccountingJournalService) DeleteJournal(id string) (int, error) {
	bqbQuery := bqb.New(`DELETE FROM "accounting.journal" WHERE id = ?`, id)
	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	result, err := service.DB.Exec(query, params...)
	if err != nil {
		switch err.(*pq.Error).Constraint {
		case database.FK_ACCOUNTING_JOURNAL_ID:
			return 400, ErrUnableToDeleteCurrentlyUsedJournal
		}

		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	if n, _ := result.RowsAffected(); n == 0 {
		return 404, ErrJournalNotFound
	}

	return 200, nil
}
