package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"system.buon18.com/m/models/setting"
	"system.buon18.com/m/services"
	"system.buon18.com/m/utils"
)

type AuthHandler struct {
	DB            *sql.DB
	ServiceFacade *services.ServiceFacade
}

func (handler *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(utils.CtxKey{}).(*utils.CtxValue)
	w.Header().Set("Content-Type", "application/json")

	pResponse := make([]setting.SettingPermissionResponse, 0)
	for _, permission := range *ctx.Permissions {
		pResponse = append(pResponse, setting.SettingPermissionToResponse(permission))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.NewResponse(http.StatusOK, "", map[string]interface{}{
		"user": setting.SettingUserToResponse(*ctx.User, setting.SettingRoleToResponse(*ctx.Role, pResponse)),
	}))
}

func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[services.LoginRequest](r, w, false)
	if !ok {
		return
	}

	tokenAndRefreshToken, statusCode, err := handler.ServiceFacade.AuthService.Login(&req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "", tokenAndRefreshToken))
}

func (handler *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[services.RefreshTokenRequest](r, w, false)
	if !ok {
		return
	}

	tokenAndRefreshToken, statusCode, err := handler.ServiceFacade.AuthService.RefreshToken(&req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "", tokenAndRefreshToken))
}

func (handler *AuthHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(utils.CtxKey{}).(*utils.CtxValue)
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[services.UpdatePasswordRequest](r, w, false)
	if !ok {
		return
	}

	message, statusCode, err := handler.ServiceFacade.AuthService.UpdatePassword(ctx, &req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, message, nil))
}

func (handler *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(utils.CtxKey{}).(*utils.CtxValue)
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[services.UpdateProfileRequest](r, w, false)
	if !ok {
		return
	}

	message, statusCode, err := handler.ServiceFacade.AuthService.UpdateProfile(ctx, &req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, message, nil))
}
