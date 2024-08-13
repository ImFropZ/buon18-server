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
	ErrPaymentTermNotFound             = errors.New("payment term not found")
	ErrAccountingPaymentTermNameExists = errors.New("accounting payment term name already exists")
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
			paymentTermLinesResponse := make([]accounting.AccountingPaymentTermLineResponse, 0)
			for _, paymentTermLine := range paymentTermLines {
				paymentTermLinesResponse = append(paymentTermLinesResponse, accounting.AccountingPaymentTermLineToResponse(paymentTermLine))
			}
			paymentTermsResponse = append(paymentTermsResponse, accounting.AccountingPaymentTermToResponse(lastPaymentTerm, paymentTermLinesResponse))
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
		paymentTermLinesResponse := make([]accounting.AccountingPaymentTermLineResponse, 0)
		for _, paymentTermLine := range paymentTermLines {
			paymentTermLinesResponse = append(paymentTermLinesResponse, accounting.AccountingPaymentTermLineToResponse(paymentTermLine))
		}
		paymentTermsResponse = append(paymentTermsResponse, accounting.AccountingPaymentTermToResponse(lastPaymentTerm, paymentTermLinesResponse))
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

func (service *AccountingPaymentTermService) PaymentTerm(id string) (accounting.AccountingPaymentTermResponse, int, error) {
	bqbQuery := bqb.New(`
	WITH "limited_payment_terms" AS (
		SELECT
			id,
			name,
			description
		FROM
			"accounting.payment_term"
		WHERE id = ?)
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
	INNER JOIN "accounting.payment_term_line" ON "limited_payment_terms".id = "accounting.payment_term_line".accounting_payment_term_id
	ORDER BY "limited_payment_terms".id ASC, "accounting.payment_term_line".sequence ASC`, id)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("%v", err)
		return accounting.AccountingPaymentTermResponse{}, 500, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		log.Printf("%v", err)
		return accounting.AccountingPaymentTermResponse{}, 500, utils.ErrInternalServer
	}

	paymentTerm := accounting.AccountingPaymentTerm{}
	paymentTermLines := make([]accounting.AccountingPaymentTermLine, 0)
	for rows.Next() {
		tmpPaymentTermLine := accounting.AccountingPaymentTermLine{}
		err = rows.Scan(
			&paymentTerm.Id,
			&paymentTerm.Name,
			&paymentTerm.Description,
			&tmpPaymentTermLine.Id,
			&tmpPaymentTermLine.Sequence,
			&tmpPaymentTermLine.ValueAmountPercent,
			&tmpPaymentTermLine.NumberOfDays,
		)
		if err != nil {
			log.Printf("%v", err)
			return accounting.AccountingPaymentTermResponse{}, 500, utils.ErrInternalServer
		}

		paymentTermLines = append(paymentTermLines, tmpPaymentTermLine)
	}

	if paymentTerm.Id == 0 {
		return accounting.AccountingPaymentTermResponse{}, 404, ErrPaymentTermNotFound
	}

	paymentTermLinesResponse := make([]accounting.AccountingPaymentTermLineResponse, 0)
	for _, paymentTermLine := range paymentTermLines {
		paymentTermLinesResponse = append(paymentTermLinesResponse, accounting.AccountingPaymentTermLineToResponse(paymentTermLine))
	}

	return accounting.AccountingPaymentTermToResponse(paymentTerm, paymentTermLinesResponse), 200, nil
}

func (service *AccountingPaymentTermService) CreatePaymentTerm(ctx *utils.CtxW, paymentTerm *accounting.AccountingPaymentTermCreateRequest) (int, error) {
	commonModel := models.CommonModel{}
	commonModel.PrepareForCreate(ctx.User.Id, ctx.User.Id)

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	bqbQuery := bqb.New(`INSERT INTO "accounting.payment_term" 
	(name, description, cid, ctime, mid, mtime) 
	VALUES 
	(?, ?, ?, ?, ?, ?) RETURNING id`, paymentTerm.Name, paymentTerm.Description, commonModel.CId, commonModel.CTime, commonModel.MId, commonModel.MTime)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("%v", err)
		tx.Rollback()
		return 500, utils.ErrInternalServer
	}

	var paymentTermId int
	err = tx.QueryRow(query, params...).Scan(&paymentTermId)
	if err != nil {
		tx.Rollback()
		switch err.(*pq.Error).Constraint {
		case database.KEY_ACCOUNTING_PAYMENT_TERM_NAME:
			return 409, ErrAccountingPaymentTermNameExists
		}

		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	bqbQuery = bqb.New(`INSERT INTO "accounting.payment_term_line"
	(sequence, value_amount_percent, number_of_days, accounting_payment_term_id, cid, ctime, mid, mtime)
	VALUES`)
	for index, paymentTermLine := range paymentTerm.Lines {
		bqbQuery.Space(`(?, ?, ?, ?, ?, ?, ?, ?)`, paymentTermLine.Sequence, paymentTermLine.ValueAmountPercent, paymentTermLine.NumberOfDays, paymentTermId, commonModel.CId, commonModel.CTime, commonModel.MId, commonModel.MTime)
		if index != len(paymentTerm.Lines)-1 {
			bqbQuery.Space(`,`)
		}
	}

	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("%v", err)
		tx.Rollback()
		return 500, utils.ErrInternalServer
	}

	_, err = tx.Exec(query, params...)
	if err != nil {
		tx.Rollback()
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	return 201, nil
}
