package accounting

import (
	"database/sql"
	"errors"
	"log"
	"sync"

	"system.buon18.com/m/database"
	"system.buon18.com/m/models"
	"system.buon18.com/m/models/accounting"
	"system.buon18.com/m/utils"

	"github.com/lib/pq"
	"github.com/nullism/bqb"
)

var (
	ErrJournalEntryNotFound             = errors.New("journal entry not found")
	ErrAccountingJournalEntryNameExists = errors.New("accounting journal entry name already exists")
	ErrBothDebitAndCreditZero           = errors.New("amount debit and credit cannot be zero")
	ErrJournalEntryNotAllowToBeDeleted  = errors.New("journal entries in posted or cancelled status are not allow to be deleted")
)

type AccountingJournalEntryService struct {
	DB *sql.DB
}

func (service *AccountingJournalEntryService) JournalEntries(qp *utils.QueryParams) ([]accounting.AccountingJournalEntryResponse, int, int, error) {
	bqbQuery := bqb.New(`WITH limited_journal_entries AS (
    SELECT
        "accounting.journal_entry".id,
		"accounting.journal_entry".name,
		"accounting.journal_entry".date,
		"accounting.journal_entry".note,
		"accounting.journal_entry".status,
		"accounting.journal_entry".accounting_journal_id
    FROM
        "accounting.journal_entry"`)

	qp.FilterIntoBqb(bqbQuery)
	qp.PaginationIntoBqb(bqbQuery)

	bqbQuery.Space(`)
	SELECT
		"limited_journal_entries".id,
		"limited_journal_entries".name,
		"limited_journal_entries".date,
		"limited_journal_entries".note,
		"limited_journal_entries".status,
		"accounting.journal_entry_line".id,
		"accounting.journal_entry_line".sequence,
		"accounting.journal_entry_line".name,
		"accounting.journal_entry_line".amount_debit,
		"accounting.journal_entry_line".amount_credit,
		"accounting.account".id,
		"accounting.account".name,
		"accounting.account".code,
		"accounting.account".typ,
		"accounting.journal".id,
		"accounting.journal".code,
		"accounting.journal".name,
		"accounting.journal".typ
	FROM
		"limited_journal_entries"
	INNER JOIN "accounting.journal_entry_line" ON "accounting.journal_entry_line".accounting_journal_entry_id = "limited_journal_entries".id
	INNER JOIN "accounting.account" ON "accounting.account".id = "accounting.journal_entry_line".accounting_account_id
	INNER JOIN "accounting.journal" ON "accounting.journal".id = "limited_journal_entries".accounting_journal_id`)

	qp.OrderByIntoBqb(bqbQuery, `"limited_journal_entries".id ASC, "accounting.journal_entry_line".sequence ASC`)

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

	journalEntriesResponse := make([]accounting.AccountingJournalEntryResponse, 0)
	lastJournalEntry := accounting.AccountingJournalEntry{}
	lastJournal := accounting.AccountingJournal{}
	journalEntryLinesResponse := make([]accounting.AccountingJournalEntryLineResponse, 0)
	for rows.Next() {
		tmpJournalEntry := accounting.AccountingJournalEntry{}
		tmpJournalEntryLine := accounting.AccountingJournalEntryLine{}
		tmpAccount := accounting.AccountingAccount{}
		tmpJournal := accounting.AccountingJournal{}
		err := rows.Scan(
			&tmpJournalEntry.Id,
			&tmpJournalEntry.Name,
			&tmpJournalEntry.Date,
			&tmpJournalEntry.Note,
			&tmpJournalEntry.Status,
			&tmpJournalEntryLine.Id,
			&tmpJournalEntryLine.Sequence,
			&tmpJournalEntryLine.Name,
			&tmpJournalEntryLine.AmountDebit,
			&tmpJournalEntryLine.AmountCredit,
			&tmpAccount.Id,
			&tmpAccount.Name,
			&tmpAccount.Code,
			&tmpAccount.Typ,
			&tmpJournal.Id,
			&tmpJournal.Code,
			&tmpJournal.Name,
			&tmpJournal.Typ,
		)
		if err != nil {
			log.Printf("%v", err)
			return nil, 0, 500, utils.ErrInternalServer
		}

		if lastJournal.Id != tmpJournal.Id {
			if lastJournalEntry.Id != 0 {
				journalResponse := accounting.AccountingJournalToResponse(lastJournal, nil)
				journalEntriesResponse = append(journalEntriesResponse, accounting.AccountingJournalEntryToResponse(lastJournalEntry, journalEntryLinesResponse, journalResponse))
			}

			// Reset and append new data
			lastJournalEntry = tmpJournalEntry
			lastJournal = tmpJournal
			journalEntryLinesResponse = make([]accounting.AccountingJournalEntryLineResponse, 0)
			accountResponse := accounting.AccountingAccountToResponse(tmpAccount)
			journalEntryLinesResponse = append(journalEntryLinesResponse, accounting.AccountingJournalEntryLineToResponse(tmpJournalEntryLine, accountResponse))
			continue
		}

		if lastJournalEntry.Id != tmpJournalEntry.Id {
			lastJournalEntry = tmpJournalEntry
		}

		accountResponse := accounting.AccountingAccountToResponse(tmpAccount)
		journalEntryLinesResponse = append(journalEntryLinesResponse, accounting.AccountingJournalEntryLineToResponse(tmpJournalEntryLine, accountResponse))
	}

	if lastJournalEntry.Id != 0 {
		journalResponse := accounting.AccountingJournalToResponse(lastJournal, nil)
		journalEntriesResponse = append(journalEntriesResponse, accounting.AccountingJournalEntryToResponse(lastJournalEntry, journalEntryLinesResponse, journalResponse))
	}

	bqbQuery = bqb.New(`SELECT COUNT(*) FROM "accounting.journal_entry"`)
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

	return journalEntriesResponse, total, 200, nil
}

func (service *AccountingJournalEntryService) JournalEntry(id string) (accounting.AccountingJournalEntryResponse, int, error) {
	bqbQuery := bqb.New(`WITH limited_journal_entries AS (
		SELECT
			"accounting.journal_entry".id,
			"accounting.journal_entry".name,
			"accounting.journal_entry".date,
			"accounting.journal_entry".note,
			"accounting.journal_entry".status,
			"accounting.journal_entry".accounting_journal_id
		FROM
			"accounting.journal_entry"
		WHERE
			"accounting.journal_entry".id = ?
		)
		SELECT
			"limited_journal_entries".id,
			"limited_journal_entries".name,
			"limited_journal_entries".date,
			"limited_journal_entries".note,
			"limited_journal_entries".status,
			"accounting.journal_entry_line".id,
			"accounting.journal_entry_line".sequence,
			"accounting.journal_entry_line".name,
			"accounting.journal_entry_line".amount_debit,
			"accounting.journal_entry_line".amount_credit,
			"accounting.account".id,
			"accounting.account".name,
			"accounting.account".code,
			"accounting.account".typ,
			"accounting.journal".id,
			"accounting.journal".code,
			"accounting.journal".name,
			"accounting.journal".typ
		FROM
			"limited_journal_entries"
		INNER JOIN "accounting.journal_entry_line" ON "accounting.journal_entry_line".accounting_journal_entry_id = "limited_journal_entries".id
		INNER JOIN "accounting.account" ON "accounting.account".id = "accounting.journal_entry_line".accounting_account_id
		INNER JOIN "accounting.journal" ON "accounting.journal".id = "limited_journal_entries".accounting_journal_id
		ORDER BY "limited_journal_entries".id ASC, "accounting.journal_entry_line".sequence ASC`, id)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("%v", err)
		return accounting.AccountingJournalEntryResponse{}, 500, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		log.Printf("%v", err)
		return accounting.AccountingJournalEntryResponse{}, 500, utils.ErrInternalServer
	}

	journalEntry := accounting.AccountingJournalEntry{}
	journal := accounting.AccountingJournal{}
	journalEntryLinesResponse := make([]accounting.AccountingJournalEntryLineResponse, 0)
	for rows.Next() {
		tmpJournalEntryLine := accounting.AccountingJournalEntryLine{}
		tmpAccount := accounting.AccountingAccount{}
		err := rows.Scan(
			&journalEntry.Id,
			&journalEntry.Name,
			&journalEntry.Date,
			&journalEntry.Note,
			&journalEntry.Status,
			&tmpJournalEntryLine.Id,
			&tmpJournalEntryLine.Sequence,
			&tmpJournalEntryLine.Name,
			&tmpJournalEntryLine.AmountDebit,
			&tmpJournalEntryLine.AmountCredit,
			&tmpAccount.Id,
			&tmpAccount.Name,
			&tmpAccount.Code,
			&tmpAccount.Typ,
			&journal.Id,
			&journal.Code,
			&journal.Name,
			&journal.Typ,
		)
		if err != nil {
			log.Printf("%v", err)
			return accounting.AccountingJournalEntryResponse{}, 500, utils.ErrInternalServer
		}

		accountResponse := accounting.AccountingAccountToResponse(tmpAccount)
		journalEntryLinesResponse = append(journalEntryLinesResponse, accounting.AccountingJournalEntryLineToResponse(tmpJournalEntryLine, accountResponse))
	}

	if journalEntry.Id == 0 {
		return accounting.AccountingJournalEntryResponse{}, 404, ErrJournalEntryNotFound
	}

	journalResponse := accounting.AccountingJournalToResponse(journal, nil)

	return accounting.AccountingJournalEntryToResponse(journalEntry, journalEntryLinesResponse, journalResponse), 200, nil
}

func (service *AccountingJournalEntryService) CreateJournalEntry(ctx *utils.CtxW, journalEntry *accounting.AccountingJournalEntryCreateRequest) (int, error) {
	commonModel := models.CommonModel{}
	commonModel.PrepareForCreate(ctx.User.Id, ctx.User.Id)

	// Check if lines are empty debit and credit
	for _, line := range journalEntry.Lines {
		if line.AmountDebit == 0 && line.AmountCredit == 0 {
			return 400, ErrBothDebitAndCreditZero
		}
	}

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	bqbQuery := bqb.New(`INSERT INTO "accounting.journal_entry"
	(name, date, note, status, accounting_journal_id, cid, ctime, mid, mtime)
	VALUES
	(?, ?, ?, ?, ?, ?, ?, ?, ?)
	RETURNING id`, journalEntry.Name, journalEntry.Date, journalEntry.Note, journalEntry.Status, journalEntry.JournalId, commonModel.CId, commonModel.CTime, commonModel.MId, commonModel.MTime)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	var journalEntryId int
	err = tx.QueryRow(query, params...).Scan(&journalEntryId)
	if err != nil {
		tx.Rollback()

		switch err.(*pq.Error).Constraint {
		case database.KEY_ACCOUNTING_JOURNAL_ENTRY_NAME:
			return 409, ErrAccountingJournalEntryNameExists
		case database.FK_ACCOUNTING_JOURNAL_ID:
			return 400, ErrJournalNotFound
		}

		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	bqbQuery = bqb.New(`INSERT INTO "accounting.journal_entry_line"
	(sequence, name, amount_debit, amount_credit, accounting_journal_entry_id, accounting_account_id, cid, ctime, mid, mtime)
	VALUES`)
	for index, line := range journalEntry.Lines {
		bqbQuery.Space(`(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, line.Sequence, line.Name, line.AmountDebit, line.AmountCredit, journalEntryId, line.AccountId, commonModel.CId, commonModel.CTime, commonModel.MId, commonModel.MTime)
		if index != len(journalEntry.Lines)-1 {
			bqbQuery.Space(`,`)
		}
	}

	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	_, err = tx.Exec(query, params...)
	if err != nil {
		tx.Rollback()

		switch err.(*pq.Error).Constraint {
		case database.FK_ACCOUNTING_ACCOUNT_ID:
			return 404, ErrAccountNotFound
		}

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

func (service *AccountingJournalEntryService) UpdateJournalEntry(ctx *utils.CtxW, id string, journalEntry *accounting.AccountingJournalEntryUpdateRequest) (int, error) {
	commonModel := models.CommonModel{}
	commonModel.PrepareForCreate(ctx.User.Id, ctx.User.Id)

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	bqbQuery := bqb.New(`SELECT status FROM "accounting.journal_entry" WHERE id = ?`, id)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	var status string
	err = tx.QueryRow(query, params...).Scan(&status)
	if err != nil {
		tx.Rollback()
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	if status == models.AccountingJournalEntryStatusPosted || status == models.AccountingJournalEntryStatusCancelled {
		tx.Rollback()
		return 400, errors.New("cannot update journal entry with status posted or cancelled")
	}

	bqbQuery = bqb.New(`UPDATE "accounting.journal_entry" SET mid = ?, mtime = ?`, commonModel.MId, commonModel.MTime)
	utils.PrepareUpdateBqbQuery(bqbQuery, journalEntry)
	bqbQuery.Space(`WHERE id = ?`, id)

	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	_, err = tx.Exec(query, params...)
	if err != nil {
		tx.Rollback()

		switch err.(*pq.Error).Constraint {
		case database.KEY_ACCOUNTING_JOURNAL_ENTRY_NAME:
			return 409, ErrAccountingJournalEntryNameExists
		case database.FK_ACCOUNTING_JOURNAL_ID:
			return 400, ErrJournalNotFound
		}

		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	errorChan := make(chan error)
	var wg sync.WaitGroup

	if journalEntry.AddLines != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bqbQuery := bqb.New(`INSERT INTO "accounting.journal_entry_line"
			(sequence, name, amount_debit, amount_credit, accounting_journal_entry_id, accounting_account_id, cid, ctime, mid, mtime)
			VALUES`)
			for index, line := range *journalEntry.AddLines {
				bqbQuery.Space(`(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, line.Sequence, line.Name, line.AmountDebit, line.AmountCredit, id, line.AccountId, commonModel.CId, commonModel.CTime, commonModel.MId, commonModel.MTime)
				if index != len(*journalEntry.AddLines)-1 {
					bqbQuery.Space(`,`)
				}
			}

			query, params, err := bqbQuery.ToPgsql()
			if err != nil {
				errorChan <- err
				return
			}

			_, err = tx.Exec(query, params...)
			if err != nil {
				errorChan <- err
				return
			}
		}()
	}

	if journalEntry.UpdateLines != nil {
		for _, line := range *journalEntry.UpdateLines {
			wg.Add(1)
			go func() {
				defer wg.Done()
				bqbQuery := bqb.New(`UPDATE "accounting.journal_entry_line" SET mid = ?, mtime = ?`, commonModel.MId, commonModel.MTime)
				utils.PrepareUpdateBqbQuery(bqbQuery, &line)
				bqbQuery.Space(`WHERE id = ? AND accounting_journal_entry_id = ?`, line.Id, id)

				query, params, err := bqbQuery.ToPgsql()
				if err != nil {
					errorChan <- err
					return
				}

				_, err = tx.Exec(query, params...)
				if err != nil {
					errorChan <- err
					return
				}
			}()
		}
	}

	if journalEntry.DeleteLines != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bqbQuery := bqb.New(`DELETE FROM "accounting.journal_entry_line" WHERE id IN (`)
			for index, lineId := range *journalEntry.DeleteLines {
				if index == 0 {
					bqbQuery.Space(`?`, lineId)
				} else {
					bqbQuery.Comma(`?`, lineId)
				}
			}
			bqbQuery.Space(`) AND accounting_journal_entry_id = ?`, id)

			query, params, err := bqbQuery.ToPgsql()
			if err != nil {
				errorChan <- err
				return
			}

			_, err = tx.Exec(query, params...)
			if err != nil {
				errorChan <- err
				return
			}
		}()
	}

	hasError := false
	var errorMessage error
	go func() {
		for err := range errorChan {
			switch err.(*pq.Error).Constraint {
			case database.CHK_ACCOUNTING_JOURANL_ENTRY_LINE_AMOUNT:
				errorMessage = ErrBothDebitAndCreditZero
			case database.FK_ACCOUNTING_ACCOUNT_ID:
				errorMessage = ErrAccountNotFound
			}
			if !hasError {
				hasError = true
				tx.Rollback()
			}
		}
	}()

	wg.Wait()
	close(errorChan)

	if hasError {
		switch errorMessage {
		case ErrBothDebitAndCreditZero:
			return 400, ErrBothDebitAndCreditZero
		case ErrAccountNotFound:
			return 400, ErrAccountNotFound
		}

		return 500, utils.ErrInternalServer
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	return 200, nil
}

func (service *AccountingJournalEntryService) DeleteJournalEntry(id string) (int, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	bqbQuery := bqb.New(`SELECT status FROM "accounting.journal_entry" WHERE id = ?`, id)
	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	var status string
	err = tx.QueryRow(query, params...).Scan(&status)
	if err != nil {
		tx.Rollback()

		if err == sql.ErrNoRows {
			return 404, ErrJournalEntryNotFound
		}

		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	if status == models.AccountingJournalEntryStatusPosted || status == models.AccountingJournalEntryStatusCancelled {
		tx.Rollback()
		return 400, ErrJournalEntryNotAllowToBeDeleted
	}

	bqbQuery = bqb.New(`DELETE FROM "accounting.journal_entry_line" WHERE accounting_journal_entry_id = ?`, id)
	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	_, err = tx.Exec(query, params...)
	if err != nil {
		tx.Rollback()
		log.Printf("%v", err)
		return 500, utils.ErrInternalServer
	}

	bqbQuery = bqb.New(`DELETE FROM "accounting.journal_entry" WHERE id = ?`, id)
	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("%v", err)
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

	return 200, nil
}
