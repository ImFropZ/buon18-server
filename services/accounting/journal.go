package accounting

import (
	"database/sql"
	"log"
	"server/models/accounting"
	"server/utils"

	"github.com/nullism/bqb"
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

		journals = append(journals, accounting.AccountingJournalToResponse(journal, account))
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