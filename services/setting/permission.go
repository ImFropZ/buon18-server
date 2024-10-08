package setting

import (
	"database/sql"
	"log"
	"net/http"

	"system.buon18.com/m/models/setting"
	"system.buon18.com/m/utils"

	"github.com/nullism/bqb"
)

type SettingPermissionService struct {
	DB *sql.DB
}

func (service *SettingPermissionService) Permissions(qp *utils.QueryParams) ([]setting.SettingPermissionResponse, int, int, error) {
	bqbQuery := bqb.New(`SELECT "setting.permission".id, "setting.permission".name FROM "setting.permission"`)
	qp.FilterIntoBqb(bqbQuery)
	qp.OrderByIntoBqb(bqbQuery, `"setting.permission".id ASC`)
	qp.PaginationIntoBqb(bqbQuery)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("%s", err)
		return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		log.Printf("%s", err)
		return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
	}

	permissions := make([]setting.SettingPermissionResponse, 0)
	for rows.Next() {
		var permission setting.SettingPermission
		err := rows.Scan(&permission.Id, &permission.Name)
		if err != nil {
			log.Printf("%s", err)
			return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
		}
		permissions = append(permissions, setting.SettingPermissionToResponse(permission))
	}

	bqbQuery = bqb.New(`SELECT COUNT(*) FROM "setting.permission"`)
	qp.FilterIntoBqb(bqbQuery)

	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("%s", err)
		return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
	}

	var total int
	err = service.DB.QueryRow(query, params...).Scan(&total)
	if err != nil {
		log.Printf("%s", err)
		return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
	}

	return permissions, total, http.StatusOK, nil
}
