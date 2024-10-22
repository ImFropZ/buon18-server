package utils

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

var (
	ErrBadRequest     = errors.New("bad request")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrNotFound       = errors.New("not found")
	ErrForbidden      = errors.New("forbidden")
	ErrConflict       = errors.New("conflict")
	ErrInternalServer = errors.New("internal server error")

	ErrUserAccountNotFound    = errors.New("contact your administrator to create an account")
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
	ErrInvalidRefreshToken    = errors.New("invalid refresh token")
	ErrInvalidOldPassword     = errors.New("invalid old password")
	ErrUserNotFound           = errors.New("user not found")

	ErrUserEmailExists = errors.New("user email already exists")
	ErrUpdateUserPwd   = errors.New("unable to update user password")

	ErrUnableToDeleteSystemRole = errors.New("unable to delete the system role")

	ErrRoleNotFound                 = errors.New("role not found")
	ErrCreateRoleWithFullPermission = errors.New("unable to create role with full permission")

	ErrCustomerNotFound       = errors.New("customer not found")
	ErrCustomerEmailExists    = errors.New("customer email already exists")
	ErrUnableToDeleteCustomer = errors.New("unable to delete customer")

	ErrPermissionNotFound = errors.New("permission not found")

	ErrQuotationNotFound       = errors.New("quotation not found")
	ErrQuotationNameExists     = errors.New("quotation name already exists")
	ErrUnableToUpdateQuotation = errors.New("quotation in sales order or cancelled status are not allowed to be updated")
	ErrUnableToDeleteQuotation = errors.New("quotations in sales order or cancelled status are not allowed to be deleted")

	ErrOrderNotFound   = errors.New("order not found")
	ErrOrderNameExists = errors.New("order name already exists")

	ErrPaymentTermNotFound   = errors.New("payment term not found")
	ErrPaymentTermNameExists = errors.New("payment term name already exists")

	ErrAccountNotFound   = errors.New("account not found")
	ErrAccountCodeExists = errors.New("account code already exists")

	ErrJournalNotFound   = errors.New("journal not found")
	ErrJournalCodeExists = errors.New("journal code already exists")

	ErrJournalEntryNotFound       = errors.New("journal entry not found")
	ErrJournalEntryNameExists     = errors.New("journal entry name already exists")
	ErrBothDebitAndCreditZero     = errors.New("amount debit and credit cannot be zero")
	ErrUnableToDeleteJournalEntry = errors.New("journal entries in posted or cancelled status are not allow to be deleted")

	ErrResourceInUsed = errors.New("resource in used")
)

func ServerToClientError(err error) (msg string, clientErr string, code int) {
	msg = err.Error()
	switch err {
	case ErrBadRequest, ErrInvalidEmailOrPassword, ErrInvalidOldPassword,
		ErrInvalidRefreshToken, ErrBothDebitAndCreditZero:
		clientErr = "Bad Request"
		code = http.StatusBadRequest
		return
	case ErrUnauthorized:
		clientErr = "Unauthorized"
		code = http.StatusUnauthorized
		return
	case ErrNotFound, ErrUserAccountNotFound, ErrUserNotFound,
		ErrRoleNotFound, ErrCustomerNotFound, ErrQuotationNotFound,
		ErrOrderNotFound, ErrPermissionNotFound, ErrPaymentTermNotFound,
		ErrAccountNotFound, ErrJournalNotFound, ErrJournalEntryNotFound:
		clientErr = "Not Found"
		code = http.StatusNotFound
		return
	case ErrForbidden, ErrUpdateUserPwd, ErrUnableToDeleteSystemRole,
		ErrUnableToDeleteCustomer, ErrCreateRoleWithFullPermission, ErrUnableToDeleteQuotation,
		ErrUnableToDeleteJournalEntry, ErrUnableToUpdateQuotation:
		clientErr = "Forbidden"
		code = http.StatusForbidden
		return
	case ErrConflict, ErrUserEmailExists, ErrCustomerEmailExists,
		ErrQuotationNameExists, ErrOrderNameExists, ErrAccountCodeExists,
		ErrPaymentTermNameExists, ErrJournalCodeExists, ErrJournalEntryNameExists:
		clientErr = "Conflict"
		code = http.StatusConflict
		return
	}

	slog.Error(fmt.Sprintf("unhandled error: %v", err))
	return "", "internal server error", http.StatusInternalServerError
}
