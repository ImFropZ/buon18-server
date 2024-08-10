package accounting

import (
	"database/sql"
	"errors"
	"log"
	"server/models/accounting"
	"server/utils"

	"github.com/nullism/bqb"
)

var (
	ErrAccountNotFound = errors.New("account not found")
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
		log.Printf("%v", err)
		return nil, 0, 500, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		log.Printf("%v", err)
		return nil, 0, 500, utils.ErrInternalServer
	}

	accounts := []accounting.AccountingAccountResponse{}
	for rows.Next() {
		account := accounting.AccountingAccount{}
		err := rows.Scan(&account.Id, &account.Code, &account.Name, &account.Typ)
		if err != nil {
			log.Printf("%v", err)
			return nil, 0, 500, utils.ErrInternalServer
		}
		accounts = append(accounts, accounting.AccountingAccountToResponse(account))
	}

	bqbQuery = bqb.New(`SELECT COUNT(*) FROM "accounting.account"`)
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

	return accounts, total, 200, nil
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
		log.Printf("%v", err)
		return accounting.AccountingAccountResponse{}, 500, utils.ErrInternalServer
	}

	account := accounting.AccountingAccount{}
	err = service.DB.QueryRow(query, params...).Scan(&account.Id, &account.Code, &account.Name, &account.Typ)
	if err != nil {
		log.Printf("%v", err)
		return accounting.AccountingAccountResponse{}, 404, ErrAccountNotFound
	}

	return accounting.AccountingAccountToResponse(account), 200, nil
}
