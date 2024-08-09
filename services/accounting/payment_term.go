package accounting

import (
	"database/sql"
	"log"
	"server/models/accounting"
	"server/utils"

	"github.com/nullism/bqb"
)

type AccountingPaymentTermService struct {
	DB *sql.DB
}

func (service *AccountingPaymentTermService) PaymentTerms(qp *utils.QueryParams) ([]accounting.AccountingPaymentTermResponse, int, int, error) {
	bqbQuery := bqb.New(`
	WITH "limited_payment_terms" AS (
		SELECT
			id,
			name,
			description
		FROM
			"accounting.payment_term"`)
	qp.FilterIntoBqb(bqbQuery)
	qp.PaginationIntoBqb(bqbQuery)

	bqbQuery.Space(`)
	SELECT 
		"limited_payment_terms".id,
		"limited_payment_terms".name,
		"limited_payment_terms".description,
		"accounting.payment_term_line".id,
		"accounting.payment_term_line".sequence,
		"accounting.payment_term_line".value_amount_percent,
		"accounting.payment_term_line".number_of_days
	FROM
		"limited_payment_terms"
	INNER JOIN "accounting.payment_term_line" ON "limited_payment_terms".id = "accounting.payment_term_line".accounting_payment_term_id`)

	qp.OrderByIntoBqb(bqbQuery, `"limited_payment_terms".id ASC, "accounting.payment_term_line".sequence ASC`)

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

	paymentTermsResponse := make([]accounting.AccountingPaymentTermResponse, 0)
	lastPaymentTerm := accounting.AccountingPaymentTerm{}
	paymentTermLines := make([]accounting.AccountingPaymentTermLine, 0)
	for rows.Next() {
		var tmpPaymentTerm accounting.AccountingPaymentTerm
		var tmpPaymentTermLine accounting.AccountingPaymentTermLine
		err = rows.Scan(
			&tmpPaymentTerm.Id,
			&tmpPaymentTerm.Name,
			&tmpPaymentTerm.Description,
			&tmpPaymentTermLine.Id,
			&tmpPaymentTermLine.Sequence,
			&tmpPaymentTermLine.ValueAmountPercent,
			&tmpPaymentTermLine.NumberOfDays,
		)
		if err != nil {
			log.Printf("%v", err)
			return nil, 0, 500, utils.ErrInternalServer
		}

		if lastPaymentTerm.Id != tmpPaymentTerm.Id && lastPaymentTerm.Id != 0 {
			paymentTermsResponse = append(paymentTermsResponse, accounting.AccountingPaymentTermToResponse(lastPaymentTerm, paymentTermLines))
			lastPaymentTerm = tmpPaymentTerm
			paymentTermLines = make([]accounting.AccountingPaymentTermLine, 0)
			paymentTermLines = append(paymentTermLines, tmpPaymentTermLine)
			continue
		}

		if lastPaymentTerm.Id == 0 {
			lastPaymentTerm = tmpPaymentTerm
		}

		paymentTermLines = append(paymentTermLines, tmpPaymentTermLine)
	}

	if lastPaymentTerm.Id != 0 {
		paymentTermsResponse = append(paymentTermsResponse, accounting.AccountingPaymentTermToResponse(lastPaymentTerm, paymentTermLines))
	}

	bqbQuery = bqb.New(`SELECT COUNT(*) FROM "accounting.payment_term"`)
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

	return paymentTermsResponse, total, 200, nil
}
