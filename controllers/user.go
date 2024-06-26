package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"server/database"
	"server/models"
	"server/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/nullism/bqb"
)

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Deleted  string `json:"deleted"`
}

type UserHandler struct {
	DB *sql.DB
}

func (handler *UserHandler) First(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid user Id. user Id should be an integer"))
		return
	}

	// -- Prepare sql query
	query, params, err := bqb.New("SELECT id, name, email, role FROM \"user\" WHERE id = ?", id).ToPgsql()
	if err != nil {
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	var user models.User
	if row := handler.DB.QueryRow(query, params...); row.Err() != nil {
		log.Printf("Error finding user in database: %v\n", row.Err())
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Role); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, utils.NewErrorResponse(404, "user not found"))
			return
		}

		log.Printf("Error scanning user from database: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, "success", user.ToResponse()))
}

func (handler *UserHandler) List(c *gin.Context) {
	paginationQueryParams := utils.PaginationQueryParams{
		Offset: 0,
		Limit:  10,
	}

	// -- Parse query params
	paginationQueryParams.Parse(c)

	// -- Prepare sql query
	query, params, err := bqb.New("SELECT id, name, email, role FROM \"user\" ORDER BY id LIMIT ? OFFSET ?", paginationQueryParams.Limit, paginationQueryParams.Offset).ToPgsql()
	if err != nil {
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Query users from database
	var users []models.User = make([]models.User, 0)
	if rows, err := handler.DB.Query(query, params...); err != nil {
		log.Printf("Error finding users in database: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Role); err != nil {
				log.Printf("Error scanning user from database: %v\n", err)
				c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
				return
			}
			users = append(users, user)
		}
	}

	c.JSON(200, utils.NewResponse(200, "success", models.UsersToResponse(users)))
}

func (handler *UserHandler) Create(c *gin.Context) {
	// -- Get email
	var userId uint
	if id, ok := c.Get("user_id"); !ok {
		log.Printf("Error getting user Id from context: %v\n", errors.New("user Id not found in context"))
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		userId = id.(uint)
	}

	// -- Bind request
	var request CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request. request should contain name, email, password, and role fields"))
		return
	}

	// -- Validate role
	role, ok := utils.ValidateRole(request.Role)
	if !ok {
		c.JSON(400, utils.NewErrorResponse(400, "invalid role"))
		return
	}

	// -- Hash password
	hashedPwd, err := utils.HashPwd(request.Password)
	if err != nil {
		log.Printf("Error hashing password: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Prepare user model
	user := models.User{
		Name:  request.Name,
		Email: request.Email,
		Pwd:   hashedPwd,
		Role:  role,
	}

	// -- Prepare for create
	if err := user.PrepareForCreate(userId, userId); err != nil {
		log.Printf("Error preparing create fields for user: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Prepare sql query
	query, params, err := bqb.New(`INSERT INTO "user" (name, email, pwd, role, cid, ctime, mid, mtime) 
	VALUES 
		(?, ?, ?, ?, ?, ?, ?, ?)
	RETURNING id`, user.Name, user.Email, user.Pwd, user.Role, user.CId, user.CTime, user.MId, user.MTime).ToPgsql()
	if err != nil {
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Open transaction
	tx, err := handler.DB.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Create user
	var createdUserId uint
	if row := tx.QueryRow(query, params...); row.Err() != nil {
		tx.Rollback()
		if pqErr, ok := row.Err().(*pq.Error); ok {
			if pqErr.Code == pq.ErrorCode(database.PQ_ERROR_CODES[database.DUPLICATE]) {
				c.JSON(400, utils.NewErrorResponse(400, "email already exists"))
				return
			}
		}

		log.Printf("Error creating new user: %v\n", row.Err())
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		if err := row.Scan(&createdUserId); err != nil {
			tx.Rollback()
			log.Printf("Error scanning created user from database: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}

		if err := tx.Commit(); err != nil {
			log.Printf("Error committing transaction: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}
	}

	c.JSON(201, utils.NewResponse(201, fmt.Sprintf("user %d created", createdUserId), ""))
}

func (handler *UserHandler) Update(c *gin.Context) {
	// -- Get user id
	var userId uint
	if id, ok := c.Get("user_id"); !ok {
		log.Printf("Error getting user Id from context: %v\n", errors.New("user Id not found in context"))
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		userId = id.(uint)
	}

	// -- Get update user Id
	updateUserId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid user Id. user Id should be an integer"))
		return
	}

	// -- Bind request
	var request UpdateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request. request should contain either one of name, email, password, and role fields"))
		return
	}

	// -- Validate role
	var role string
	if request.Role != "" {
		if roleStr, ok := utils.ValidateRole(request.Role); !ok {
			c.JSON(400, utils.NewErrorResponse(400, "invalid role"))
			return
		} else {
			role = roleStr
		}
	}

	// -- Prepare sql query
	query, params, err := bqb.New(`SELECT id, name, email, pwd, role, deleted FROM "user" WHERE id = ?`, updateUserId).ToPgsql()
	if err != nil {
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Begin transaction
	tx, err := handler.DB.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Query user by Id
	var user models.User
	if result := tx.QueryRow(query, params...); result.Err() != nil {
		tx.Rollback()
		log.Printf("Error finding user in database: %v\n", result.Err())
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else if err := result.Scan(&user.Id, &user.Name, &user.Email, &user.Pwd, &user.Role, &user.Deleted); err != nil {
		tx.Rollback()
		log.Printf("Error scanning user from database: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Check if user already deleted
	if user.Deleted {
		if lower := strings.ToLower(request.Deleted); lower != "false" || lower == "true" {
			tx.Rollback()
			c.JSON(400, utils.NewErrorResponse(400, "user already deleted"))
			return
		}
	}

	// -- Prepare for update
	if err := user.PrepareForUpdate(userId); err != nil {
		tx.Rollback()
		log.Printf("Error preparing update fields for user: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Prepare sql query
	query, _, err = bqb.New(`SELECT COUNT(id) FROM "user" WHERE role = 'Admin'`).ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Check if user is updating an only admin
	var count int64
	if row := tx.QueryRow(query); row.Err() != nil {
		tx.Rollback()
		log.Printf("Error finding admin role in database: %v\n", row.Err())
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else if err := row.Scan(&count); err != nil {
		tx.Rollback()
		log.Printf("Error scanning admin role from database: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}
	if count < 2 && user.Role == "Admin" && role != "Admin" {
		tx.Rollback()
		c.JSON(400, utils.NewErrorResponse(400, "cannot update the only admin"))
		return
	}

	// -- Update user
	bqbQuery := bqb.New(`UPDATE "user" SET`)
	if request.Name != "" && request.Name != user.Name {
		bqbQuery.Space(`name = ?,`, request.Name)
	}
	if request.Email != "" && request.Email != user.Email {
		bqbQuery.Space(`email = ?,`, request.Email)
	}
	if request.Password != "" {
		hashedPwd, err := utils.HashPwd(request.Password)
		if err != nil {
			tx.Rollback()
			log.Printf("Error hashing password: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}
		bqbQuery.Space(`pwd = ?,`, hashedPwd)
	}
	if role != "" && role != user.Role {
		bqbQuery.Space(`role = ?,`, role)
	}
	if lower := strings.ToLower(request.Deleted); lower == "true" || lower == "false" {
		bqbQuery.Space(`deleted = ?,`, lower == "true")
	}

	// -- Uppdate timestamp
	bqbQuery.Space(`mid = ?, mtime = ? WHERE id = ?`, user.MId, user.MTime, updateUserId)

	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Update user to database
	if result, err := tx.Exec(query, params...); err != nil {
		tx.Rollback()
		log.Printf("Error updating user: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		if n, err := result.RowsAffected(); err != nil {
			tx.Rollback()
			log.Printf("Error getting rows affected: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		} else if n == 0 {
			tx.Rollback()
			c.JSON(400, utils.NewErrorResponse(400, "user not updated"))
			return
		}
	}

	// -- Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, fmt.Sprintf("user %d updated", user.Id), ""))
}

func (handler *UserHandler) Delete(c *gin.Context) {
	// -- Get user id
	var userId uint
	if id, ok := c.Get("user_id"); !ok {
		log.Printf("Error getting user Id from context: %v\n", errors.New("user Id not found in context"))
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		userId = id.(uint)
	}

	// -- Get delete user Id
	targetUserId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid user Id. user Id should be an integer"))
		return
	}

	// -- Prepare sql query
	query, params, err := bqb.New(`SELECT id, email, deleted FROM "user" WHERE id = ?`, targetUserId).ToPgsql()
	if err != nil {
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Begin transaction
	tx, err := handler.DB.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Query user by Id
	var targetUser models.User
	if err := tx.QueryRow(query, params...).Scan(&targetUser.Id, &targetUser.Email, &targetUser.Deleted); err != nil {
		tx.Rollback()
		log.Printf("Error finding user in database: %v\n", err)
		c.JSON(404, utils.NewErrorResponse(404, "user not found"))
		return
	}

	// -- Check if user already deleted
	if targetUser.Deleted {
		c.JSON(400, utils.NewErrorResponse(400, "user already deleted"))
		return
	}

	// -- Check if user is deleting itself
	if targetUser.Id == userId {
		c.JSON(400, utils.NewErrorResponse(400, "user cannot delete itself"))
		return
	}

	// -- Prepare for update
	if err := targetUser.PrepareForUpdate(userId); err != nil {
		log.Printf("Error preparing update fields for user: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Prepare sql query
	query, _, err = bqb.New(`SELECT COUNT(id) FROM "user" WHERE role = 'Admin'`).ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Check if user is deleting an only admin
	var count int64
	if row := tx.QueryRow(query); row.Err() != nil {
		tx.Rollback()
		log.Printf("Error finding admin role in database: %v\n", row.Err())
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else if err := row.Scan(&count); err != nil {
		tx.Rollback()
		log.Printf("Error scanning admin role from database: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}
	if count < 2 && targetUser.Role == "Admin" {
		tx.Rollback()
		c.JSON(400, utils.NewErrorResponse(400, "cannot update the only admin"))
		return
	}

	// -- Prepare sql query
	query, params, err = bqb.New(`UPDATE "user" SET deleted = true WHERE id = ?`, targetUser.Id).ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Delete user
	if result, err := tx.Exec(query, params...); err != nil {
		tx.Rollback()
		log.Printf("Error deleting user: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		if n, err := result.RowsAffected(); err != nil {
			tx.Rollback()
			log.Printf("Error getting rows affected: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		} else if n == 0 {
			tx.Rollback()
			c.JSON(400, utils.NewErrorResponse(400, "user not deleted"))
			return
		}
	}

	// -- Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, fmt.Sprintf("user %d deleted", targetUser.Id), nil))
}
