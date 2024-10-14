package accounting

import (
	"database/sql"
	"errors"
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

type AccountingJournalEntryService struct {
	DB *sql.DB
}

func (service *AccountingJournalEntryService) JournalEntries(qp *utils.QueryParams) ([]accounting.AccountingJournalEntryResponse, int, int, error) {
	bqbQuery := bqb.New(`WITH limited_journal_entries AS (
    SELECT
		*
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
		slog.Error(fmt.Sprintf("%v", err))
		return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
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
			slog.Error(fmt.Sprintf("%v", err))
			return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
		}

		if lastJournal.Id != nil && *lastJournal.Id != *tmpJournal.Id {
			journalResponse := accounting.AccountingJournalToResponse(lastJournal, nil)
			journalEntriesResponse = append(journalEntriesResponse, accounting.AccountingJournalEntryToResponse(lastJournalEntry, &journalEntryLinesResponse, &journalResponse))

			journalEntryLinesResponse = make([]accounting.AccountingJournalEntryLineResponse, 0)
		}

		lastJournalEntry = tmpJournalEntry
		accountResponse := accounting.AccountingAccountToResponse(tmpAccount)
		journalEntryLinesResponse = append(journalEntryLinesResponse, accounting.AccountingJournalEntryLineToResponse(tmpJournalEntryLine, &accountResponse))
	}

	if lastJournalEntry.Id != nil {
		journalResponse := accounting.AccountingJournalToResponse(lastJournal, nil)
		journalEntriesResponse = append(journalEntriesResponse, accounting.AccountingJournalEntryToResponse(lastJournalEntry, &journalEntryLinesResponse, &journalResponse))
	}

	bqbQuery = bqb.New(`SELECT COUNT(*) FROM "accounting.journal_entry"`)
	qp.FilterIntoBqb(bqbQuery)

	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
	}

	var total int
	err = service.DB.QueryRow(query, params...).Scan(&total)
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
	}

	return journalEntriesResponse, total, http.StatusOK, nil
}

func (service *AccountingJournalEntryService) JournalEntry(id string) (accounting.AccountingJournalEntryResponse, int, error) {
	bqbQuery := bqb.New(`WITH limited_journal_entries AS (
		SELECT
			*
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
		slog.Error(fmt.Sprintf("%v", err))
		return accounting.AccountingJournalEntryResponse{}, http.StatusInternalServerError, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return accounting.AccountingJournalEntryResponse{}, http.StatusInternalServerError, utils.ErrInternalServer
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
			slog.Error(fmt.Sprintf("%v", err))
			return accounting.AccountingJournalEntryResponse{}, http.StatusInternalServerError, utils.ErrInternalServer
		}

		accountResponse := accounting.AccountingAccountToResponse(tmpAccount)
		journalEntryLinesResponse = append(journalEntryLinesResponse, accounting.AccountingJournalEntryLineToResponse(tmpJournalEntryLine, &accountResponse))
	}

	if journalEntry.Id == nil {
		return accounting.AccountingJournalEntryResponse{}, http.StatusNotFound, utils.ErrJournalEntryNotFound
	}

	journalResponse := accounting.AccountingJournalToResponse(journal, nil)
	return accounting.AccountingJournalEntryToResponse(journalEntry, &journalEntryLinesResponse, &journalResponse), http.StatusOK, nil
}

func (service *AccountingJournalEntryService) CreateJournalEntry(ctx *utils.CtxValue, journalEntry *accounting.AccountingJournalEntryCreateRequest) (int, error) {
	commonModel := models.CommonModel{}
	commonModel.PrepareForCreate(*ctx.User.Id, *ctx.User.Id)

	// Check if lines are empty debit and credit
	for _, line := range journalEntry.Lines {
		if line.AmountDebit == 0 && line.AmountCredit == 0 {
			return http.StatusBadRequest, utils.ErrBothDebitAndCreditZero
		}
	}

	tx, err := service.DB.Begin()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}
	defer tx.Rollback()

	bqbQuery := bqb.New(`INSERT INTO "accounting.journal_entry"
	(name, date, note, status, accounting_journal_id, cid, ctime, mid, mtime)
	VALUES
	(?, ?, ?, ?, ?, ?, ?, ?, ?)
	RETURNING id`, journalEntry.Name, journalEntry.Date, journalEntry.Note, journalEntry.Status, journalEntry.JournalId, commonModel.CId, commonModel.CTime, commonModel.MId, commonModel.MTime)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	var journalEntryId int
	err = tx.QueryRow(query, params...).Scan(&journalEntryId)
	if err != nil {
		switch err.(*pq.Error).Constraint {
		case database.KEY_ACCOUNTING_JOURNAL_ENTRY_NAME:
			return http.StatusConflict, utils.ErrJournalEntryNameExists
		case database.FK_ACCOUNTING_JOURNAL_ID:
			return http.StatusBadRequest, utils.ErrJournalNotFound
		}

		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
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
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if _, err = tx.Exec(query, params...); err != nil {
		switch err.(*pq.Error).Constraint {
		case database.FK_ACCOUNTING_ACCOUNT_ID:
			return http.StatusNotFound, utils.ErrAccountNotFound
		}

		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	err = tx.Commit()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	return http.StatusCreated, nil
}

func (service *AccountingJournalEntryService) UpdateJournalEntry(ctx *utils.CtxValue, id string, journalEntry *accounting.AccountingJournalEntryUpdateRequest) (int, error) {
	commonModel := models.CommonModel{}
	commonModel.PrepareForCreate(*ctx.User.Id, *ctx.User.Id)

	tx, err := service.DB.Begin()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}
	defer tx.Rollback()

	bqbQuery := bqb.New(`SELECT status FROM "accounting.journal_entry" WHERE id = ?`, id)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	var status string
	err = tx.QueryRow(query, params...).Scan(&status)
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if status == models.AccountingJournalEntryStatusPosted || status == models.AccountingJournalEntryStatusCancelled {
		return http.StatusBadRequest, errors.New("cannot update journal entry with status posted or cancelled")
	}

	bqbQuery = bqb.New(`UPDATE "accounting.journal_entry" SET mid = ?, mtime = ?`, commonModel.MId, commonModel.MTime)
	if journalEntry.Name != nil {
		bqbQuery.Space(`name = ?`, *journalEntry.Name)
	}
	if journalEntry.Date != nil {
		bqbQuery.Space(`date = ?`, *journalEntry.Date)
	}
	if journalEntry.Note != nil {
		bqbQuery.Space(`note = ?`, *journalEntry.Note)
	}
	if journalEntry.Status != nil {
		bqbQuery.Space(`status = ?`, *journalEntry.Status)
	}
	if journalEntry.JournalId != nil {
		bqbQuery.Space(`accounting_journal_id = ?`, *journalEntry.JournalId)
	}
	bqbQuery.Space(`WHERE id = ?`, id)

	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if _, err = tx.Exec(query, params...); err != nil {
		switch err.(*pq.Error).Constraint {
		case database.KEY_ACCOUNTING_JOURNAL_ENTRY_NAME:
			return http.StatusConflict, utils.ErrJournalEntryNameExists
		case database.FK_ACCOUNTING_JOURNAL_ID:
			return http.StatusBadRequest, utils.ErrJournalNotFound
		}

		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if journalEntry.AddLines != nil {
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
			slog.Error(fmt.Sprintf("%v", err))
			return http.StatusInternalServerError, utils.ErrInternalServer
		}

		if _, err := tx.Exec(query, params...); err != nil {
			switch err.(*pq.Error).Constraint {
			case database.CHK_ACCOUNTING_JOURANL_ENTRY_LINE_AMOUNT:
				return http.StatusBadRequest, utils.ErrBothDebitAndCreditZero
			}

			slog.Error(fmt.Sprintf("%v", err))
			return http.StatusInternalServerError, utils.ErrInternalServer
		}
	}

	if journalEntry.UpdateLines != nil {
		for _, line := range *journalEntry.UpdateLines {
			bqbQuery := bqb.New(`UPDATE "accounting.journal_entry_line" SET mid = ?, mtime = ?`, commonModel.MId, commonModel.MTime)
			if line.Sequence != nil {
				bqbQuery.Space(`SET sequence = ?`, *line.Sequence)
			}
			if line.Name != nil {
				bqbQuery.Space(`SET name = ?`, *line.Name)
			}
			if line.AmountDebit != nil {
				bqbQuery.Space(`SET amount_debit = ?`, *line.AmountDebit)
			}
			if line.AmountCredit != nil {
				bqbQuery.Space(`SET amount_credit = ?`, *line.AmountCredit)
			}
			bqbQuery.Space(`WHERE id = ? AND accounting_journal_entry_id = ?`, line.Id, id)

			query, params, err := bqbQuery.ToPgsql()
			if err != nil {
				slog.Error(fmt.Sprintf("%v", err))
				return http.StatusInternalServerError, utils.ErrInternalServer
			}

			if _, err := tx.Exec(query, params...); err != nil {
				switch err.(*pq.Error).Constraint {
				case database.CHK_ACCOUNTING_JOURANL_ENTRY_LINE_AMOUNT:
					return http.StatusBadRequest, utils.ErrBothDebitAndCreditZero
				}

				slog.Error(fmt.Sprintf("%v", err))
				return http.StatusInternalServerError, utils.ErrInternalServer
			}
		}
	}

	if journalEntry.DeleteLines != nil {
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
			slog.Error(fmt.Sprintf("%v", err))
			return http.StatusInternalServerError, utils.ErrInternalServer
		}

		if _, err := tx.Exec(query, params...); err != nil {
			switch err.(*pq.Error).Constraint {
			case database.CHK_ACCOUNTING_JOURANL_ENTRY_LINE_AMOUNT:
				return http.StatusBadRequest, utils.ErrBothDebitAndCreditZero
			}

			slog.Error(fmt.Sprintf("%v", err))
			return http.StatusInternalServerError, utils.ErrInternalServer
		}
	}

	if err := tx.Commit(); err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	return http.StatusOK, nil
}

func (service *AccountingJournalEntryService) DeleteJournalEntries(req *models.CommonDelete) (int, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}
	defer tx.Rollback()

	bqbQuery := bqb.New(fmt.Sprintf(`SELECT COUNT(*) FROM "accounting.journal_entry" WHERE (status != '%s' OR status != '%s') AND id in (`, models.AccountingJournalEntryStatusPosted, models.AccountingJournalEntryStatusCancelled))
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

	var total int
	if err := tx.QueryRow(query, params...).Scan(&total); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, utils.ErrJournalEntryNotFound
		}

		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if total != len(req.Ids) {
		return http.StatusForbidden, utils.ErrUnableToDeleteJournalEntry
	}

	bqbQuery = bqb.New(`DELETE FROM "accounting.journal_entry_line" WHERE accounting_journal_entry_id in (`)
	for i, id := range req.Ids {
		bqbQuery.Space(`?`, id)

		if i < len(req.Ids)-1 {
			bqbQuery.Comma("")
		}
	}
	bqbQuery.Space(`)`)

	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if _, err := tx.Exec(query, params...); err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	bqbQuery = bqb.New(`DELETE FROM "accounting.journal_entry" WHERE id in (`)
	for i, id := range req.Ids {
		bqbQuery.Space(`?`, id)

		if i < len(req.Ids)-1 {
			bqbQuery.Comma("")
		}
	}
	bqbQuery.Space(`)`)

	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if _, err := tx.Exec(query, params...); err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if err := tx.Commit(); err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	return http.StatusOK, nil
}
