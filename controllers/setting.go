package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"system.buon18.com/m/models"
	"system.buon18.com/m/models/setting"
	"system.buon18.com/m/services"
	"system.buon18.com/m/utils"
)

type SettingHandler struct {
	DB            *sql.DB
	ServiceFacade *services.ServiceFacade
}

func (handler *SettingHandler) Users(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	qp := utils.NewQueryParams().
		PrepareFilters(setting.SettingUser{}, r, `"setting.user"`).
		PrepareSorts(setting.SettingUser{}, r, `"limited_users"`).
		PrepareLimitAndOffset(r)

	users, total, statusCode, err := handler.ServiceFacade.SettingUserService.Users(qp)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(statusCode, msg, clientErr, nil))
		return
	}

	w.Header().Set("X-Total-Count", fmt.Sprintf("%d", total))
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "", map[string]interface{}{
		"users": users,
	}))
}

func (handler *SettingHandler) User(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	user, statusCode, err := handler.ServiceFacade.SettingUserService.User(id)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(statusCode, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "", map[string]interface{}{
		"user": user,
	}))
}

func (handler *SettingHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(utils.CtxKey{}).(*utils.CtxValue)
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[setting.SettingUserCreateRequest](r, w, false)
	if !ok {
		return
	}

	statusCode, err := handler.ServiceFacade.SettingUserService.CreateUser(ctx, &req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(statusCode, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "created", nil))
}

func (handler *SettingHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(utils.CtxKey{}).(*utils.CtxValue)
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[setting.SettingUserUpdateRequest](r, w, true)
	if !ok {
		return
	}

	id := mux.Vars(r)["id"]

	// -- Prevent update role and password of system user
	if id == "1" && (req.RoleId != nil || req.Password != nil) {
		w.WriteHeader(403)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(403, "unable to update the system user's role or password", "Forbidden", nil))
		return
	}

	statusCode, err := handler.ServiceFacade.SettingUserService.UpdateUser(ctx, id, &req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(statusCode, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "updated", nil))
}

func (handler *SettingHandler) DeleteUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[models.CommonDelete](r, w, false)
	if !ok {
		return
	}

	// Check if the system user is in the list
	for _, id := range req.Ids {
		if id == 1 {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(utils.NewErrorResponse(http.StatusForbidden, "unable to delete the system user", "Forbidden", nil))
			return
		}
	}

	statusCode, err := handler.ServiceFacade.SettingUserService.DeleteUsers(&req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(statusCode, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "deleted", nil))
}

func (handler *SettingHandler) Customers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	qp := utils.NewQueryParams().
		PrepareFilters(setting.SettingCustomer{}, r, `"setting.customer"`).
		PrepareSorts(setting.SettingCustomer{}, r, `"limited_customers"`).
		PrepareLimitAndOffset(r)

	customers, total, statusCode, err := handler.ServiceFacade.SettingCustomerService.Customers(qp)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(statusCode, msg, clientErr, nil))
		return
	}

	w.Header().Set("X-Total-Count", fmt.Sprintf("%d", total))
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "", map[string]interface{}{
		"customers": customers,
	}))
}

func (handler *SettingHandler) Customer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	customer, statusCode, err := handler.ServiceFacade.SettingCustomerService.Customer(id)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(statusCode, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "", map[string]interface{}{
		"customer": customer,
	}))
}

func (handler *SettingHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(utils.CtxKey{}).(*utils.CtxValue)
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[setting.SettingCustomerCreateRequest](r, w, false)
	if !ok {
		return
	}

	statusCode, err := handler.ServiceFacade.SettingCustomerService.CreateCustomer(ctx, &req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(statusCode, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "created", nil))
}

func (handler *SettingHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(utils.CtxKey{}).(*utils.CtxValue)
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[setting.SettingCustomerUpdateRequest](r, w, true)
	if !ok {
		return
	}

	id := mux.Vars(r)["id"]
	statusCode, err := handler.ServiceFacade.SettingCustomerService.UpdateCustomer(ctx, id, &req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(statusCode, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "updated", nil))
}

func (handler *SettingHandler) DeleteCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[models.CommonDelete](r, w, false)
	if !ok {
		return
	}

	statusCode, err := handler.ServiceFacade.SettingCustomerService.DeleteCustomers(&req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(statusCode, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "deleted", nil))
}

func (handler *SettingHandler) Roles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	qp := utils.NewQueryParams().
		PrepareFilters(setting.SettingRole{}, r, `"setting.role"`).
		PrepareSorts(setting.SettingRole{}, r, `"limited_roles"`).
		PrepareLimitAndOffset(r)

	roles, total, statusCode, err := handler.ServiceFacade.SettingRoleService.Roles(qp)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(statusCode, msg, clientErr, nil))
		return
	}

	w.Header().Set("X-Total-Count", fmt.Sprintf("%d", total))
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "", map[string]interface{}{
		"roles": roles,
	}))
}

func (handler *SettingHandler) Role(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	role, statusCode, err := handler.ServiceFacade.SettingRoleService.Role(id)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(statusCode, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "", map[string]interface{}{
		"role": role,
	}))
}

func (handler *SettingHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(utils.CtxKey{}).(*utils.CtxValue)
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[setting.SettingRoleCreateRequest](r, w, false)
	if !ok {
		return
	}

	statusCode, err := handler.ServiceFacade.SettingRoleService.CreateRole(ctx, &req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(statusCode, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "created", nil))
}

func (handler *SettingHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(utils.CtxKey{}).(*utils.CtxValue)
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	if id == "1" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(403, "unable to update the system role", "Forbidden", nil))
		return
	}

	// -- Parse request
	req, ok := utils.ValidateRequest[setting.SettingRoleUpdateRequest](r, w, true)
	if !ok {
		return
	}

	if req.AddPermissionIds != nil || req.RemovePermissionIds != nil {
		needPermission := false

		if req.AddPermissionIds != nil {
			for _, permissionId := range *req.AddPermissionIds {
				if utils.ContainsString(utils.FULL_PERMISSION_IDS, utils.IntToStr(int(permissionId))) {
					needPermission = true
				}
			}
		}
		if req.RemovePermissionIds != nil {
			for _, permissionId := range *req.RemovePermissionIds {
				if utils.ContainsString(utils.FULL_PERMISSION_IDS, utils.IntToStr(int(permissionId))) {
					needPermission = true
				}
			}
		}

		if needPermission {
			hasPermission := false

			if req.AddPermissionIds != nil {
				for _, permissionId := range *req.AddPermissionIds {
					if utils.ContainsString(utils.FULL_PERMISSION_IDS, utils.IntToStr(int(permissionId))) {
						for _, permission := range *ctx.Permissions {

							if permission.Name != nil && utils.ContainsString([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCESS}, *permission.Name) {
								hasPermission = true
							}
						}
					}
				}
			}
			if req.RemovePermissionIds != nil {
				for _, permissionId := range *req.RemovePermissionIds {
					if utils.ContainsString(utils.FULL_PERMISSION_IDS, utils.IntToStr(int(permissionId))) {
						for _, permission := range *ctx.Permissions {
							if permission.Name != nil && utils.ContainsString([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCESS}, *permission.Name) {
								hasPermission = true
							}
						}
					}
				}
			}

			if !hasPermission {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(utils.NewErrorResponse(403, "unable to update role with full permission", "Forbidden", nil))
				return
			}
		}
	}

	statusCode, err := handler.ServiceFacade.SettingRoleService.UpdateRole(ctx, id, &req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(statusCode, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "updated", nil))
}

func (handler *SettingHandler) DeleteRoles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[models.CommonDelete](r, w, false)
	if !ok {
		return
	}

	// Check if the system role is in the list
	for _, roleId := range req.Ids {
		if roleId == 1 {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(utils.NewErrorResponse(http.StatusForbidden, "unable to delete the system role", "Forbidden", nil))
			return
		}
	}

	statusCode, err := handler.ServiceFacade.SettingRoleService.DeleteRoles(&req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(statusCode, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "deleted", nil))
}

func (handler *SettingHandler) Permissions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	qp := utils.NewQueryParams().
		PrepareFilters(setting.SettingPermission{}, r, `"setting.permission"`).
		PrepareSorts(setting.SettingPermission{}, r, `"limited_permissions"`).
		PrepareLimitAndOffset(r)

	permissions, total, statusCode, err := handler.ServiceFacade.SettingPermissionService.Permissions(qp)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(statusCode, msg, clientErr, nil))
		return
	}

	w.Header().Set("X-Total-Count", fmt.Sprintf("%d", total))
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "", map[string]interface{}{
		"permissions": permissions,
	}))
}
