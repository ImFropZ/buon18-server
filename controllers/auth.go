package controllers

import (
	"database/sql"
	"errors"
	"log"
	"server/models"
	"server/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nullism/bqb"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type AuthHandler struct {
	DB *sql.DB
}

func (handler *AuthHandler) Me(c *gin.Context) {
	// -- Get id
	var user_id uint
	if id, err := c.Get("user_id"); !err {
		log.Printf("Error getting user id: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		user_id = id.(uint)
	}

	// -- Prepare sql query
	query, params, err := bqb.New("SELECT name, email, role FROM \"user\" WHERE id = ?", user_id).ToPgsql()
	if err != nil {
		log.Printf("Error preparing query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Get user from db
	var user models.User
	if row := handler.DB.QueryRow(query, params...); row.Err() != nil {
		c.JSON(404, utils.NewErrorResponse(404, "user doesn't existed"))
		return
	} else {
		if err := row.Scan(&user.Name, &user.Email, &user.Role); err != nil {
			log.Printf("Error scanning user: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}
	}

	c.JSON(200, utils.NewResponse(200, "", gin.H{
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	}))
}

func (handler *AuthHandler) Login(c *gin.Context) {
	// -- Parse request
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// -- Prepare sql query
	query, params, err := bqb.New("SELECT email, COALESCE(pwd, ''), deleted FROM \"user\" WHERE email = ?", req.Email).ToPgsql()
	if err != nil {
		log.Printf("Error preparing query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Validate user
	var user models.User
	if row := handler.DB.QueryRow(query, params...); row.Err() != nil {
		log.Printf("Error querying user: %v\n", row.Err())
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		if err := row.Scan(&user.Email, &user.Pwd, &user.Deleted); err != nil {
			log.Printf("Error scanning user: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}
	}

	if user.Deleted {
		c.JSON(401, utils.NewErrorResponse(401, "Your account has been deleted"))
		return
	}

	log.Printf("User: %v\n", user)

	if user.Email != req.Email || (!utils.ComparePwd(req.Password, user.Pwd) && user.Pwd != "") {
		c.JSON(401, utils.NewErrorResponse(401, "Invalid email or password"))
		return
	}

	// -- Generate token
	token, err := utils.GenerateWebToken(utils.WebTokenClaims{
		Email: user.Email,
		Role:  user.Role,
	})

	if err != nil {
		log.Printf("Error generating web token: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(utils.RefreshTokenClaims{
		Email: user.Email,
	})
	if err != nil {
		log.Printf("Error generating refresh token: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, "", gin.H{
		"token":         token,
		"refresh_token": refreshToken,
	}))
}

func (handler *AuthHandler) RefreshToken(c *gin.Context) {
	// -- Parse request
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request. request should contain refresh_token field"))
		return
	}

	token, err := utils.RemoveBearer(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(401, utils.NewErrorResponse(401, "invalid token"))
		return
	}

	// -- Validate token
	claims, err := utils.ValidateWebToken(token)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			refreshClaims, refreshErr := utils.ValidateRefreshToken(req.RefreshToken)
			if refreshErr != nil {
				c.JSON(401, utils.NewErrorResponse(401, "invalid refresh token"))
				return
			}

			// -- Prepare sql query
			query, params, err := bqb.New("SELECT email, role FROM \"user\" WHERE email = ?", refreshClaims.Email).ToPgsql()
			if err != nil {
				log.Printf("Error preparing query: %v\n", err)
				c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
				return
			}

			// -- Check email
			var user models.User
			if row := handler.DB.QueryRow(query, params...); row.Err() != nil {
				log.Printf("Error querying user : %v\n", err)
				c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
				return
			} else {
				if err := row.Scan(&user.Email, &user.Role); err != nil {
					log.Printf("Error scanning user: %v\n", err)
					c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
					return
				}
			}

			// -- Generate new token
			token, err := utils.GenerateWebToken(utils.WebTokenClaims{
				Email: user.Email,
				Role:  user.Role,
			})
			if err != nil {
				log.Printf("Error generating web token: %v\n", err)
				c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
				return
			}

			// -- Response new token
			c.JSON(200, utils.NewResponse(200, "success", gin.H{
				"token": token,
			}))
			return
		}

		c.JSON(401, utils.NewErrorResponse(401, "invalid token"))
		return
	}

	token, err = utils.GenerateWebToken(claims)
	if err != nil {
		log.Printf("Error generating web token: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, "success", gin.H{
		"token": token,
	}))
}

func (handler *AuthHandler) UpdatePassword(c *gin.Context) {
	// -- Get user id
	var user_id uint
	if id, err := c.Get("user_id"); !err {
		log.Printf("Error getting user id: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		user_id = id.(uint)
	}

	// -- Parse request
	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request. request should contain old_password and new_password fields"))
		return
	}

	// -- Prepare sql query
	query, params, err := bqb.New("SELECT COALESCE(pwd, '') FROM \"user\" WHERE id = ?", user_id).ToPgsql()
	if err != nil {
		log.Printf("Error preparing query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Get user from db
	var user models.User
	if row := handler.DB.QueryRow(query, params...); row.Err() != nil {
		log.Printf("Error querying user: %v\n", row.Err())
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		if err := row.Scan(&user.Pwd); err != nil {
			log.Printf("Error scanning user: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}
	}

	// -- Update Pwd if user in db doesn't have pwd
	if user.Pwd != "" {
		// -- Compare pwd
		if ok := utils.ComparePwd(req.OldPassword, user.Pwd); !ok {
			c.JSON(400, utils.NewErrorResponse(400, "invalid old password"))
			return
		}
	}

	// -- Update pwd
	if hashedPwd, err := utils.HashPwd(req.NewPassword); err == nil {
		// -- Begin transaction
		var tx *sql.Tx
		if dbTx, err := handler.DB.Begin(); err != nil {
			log.Printf("Error beginning transaction: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		} else {
			tx = dbTx
		}

		// -- Prepare sql query
		query, params, err = bqb.New("UPDATE \"user\" SET pwd = ? WHERE id = ?", hashedPwd, user_id).ToPgsql()
		if err != nil {
			log.Printf("Error preparing query: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}

		// -- Update pwd
		if _, err := tx.Exec(query, params...); err != nil {
			tx.Rollback()

			log.Printf("Error updating password: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}

		// -- Commit transaction
		if err := tx.Commit(); err != nil {
			log.Printf("Error committing transaction: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}

		c.JSON(200, utils.NewResponse(200, "success", nil))
	} else {
		log.Printf("Error hashing password: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
	}
}
