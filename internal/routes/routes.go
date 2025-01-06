package routes

import (
	"net/http"
	"webapp/internal/config"
	"webapp/internal/jwt"
	"webapp/internal/models"
	hasher "webapp/internal/password_hasher"

	g "github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(ge *g.Engine, storage *gorm.DB, cfg *config.Config) {
	authGroup := ge.Group("/auth")
	{
		authGroup.POST("/register", func(ctx *g.Context) { handleRegister(ctx, storage, cfg) })
		authGroup.POST("/login", func(ctx *g.Context) { handleLogin(ctx, storage, cfg) })
	}
}

func handleRegister(ctx *g.Context, storage *gorm.DB, cfg *config.Config) {
	user := models.User{}

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, g.H{"error": err.Error()})
		return
	}

	var count int64
	storage.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, g.H{"error": "email has already taken"})
		return
	}

	hashedPassword := hasher.HashPassword(user.Password)
	user.Password = hashedPassword

	result := storage.Create(&user)

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusBadRequest, g.H{"error": result.Error})
		return
	}

	jwtToken, err := jwt.GenerateToken(&user, cfg)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, g.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, g.H{"token": jwtToken})
}

func handleLogin(ctx *g.Context, storage *gorm.DB, cfg *config.Config) {
	user := models.User{}

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, g.H{"error": err.Error()})
		return
	}

	var count int64
	storage.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)
	if count == 0 {
		ctx.JSON(http.StatusBadRequest, g.H{"error": "invalid credentials"})
		return
	}

	passwordToVerify := user.Password
	storage.First(&user, "email = ?", user.Email)

	isPasswordValid := hasher.VerifyPassword(user.Password, passwordToVerify)

	if !isPasswordValid {
		ctx.JSON(http.StatusBadRequest, g.H{"error": "invalid credentials"})
		return
	}

	jwtToken, err := jwt.GenerateToken(&user, cfg)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, g.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, g.H{"token": jwtToken})
}
