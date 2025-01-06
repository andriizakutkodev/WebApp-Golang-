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

// Middlewares
func AuthMiddleware(cfg *config.Config) g.HandlerFunc {
	return func(ctx *g.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" || !isValidToken(token, cfg) {
			ctx.JSON(http.StatusUnauthorized, g.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func isValidToken(token string, cfg *config.Config) bool {
	isValid, _ := jwt.ValidateToken(token, cfg)
	return isValid
}

// Main function to register routes
func RegisterRoutes(ge *g.Engine, db *gorm.DB, cfg *config.Config) {
	// Authentication routes: login and register
	authGroup := ge.Group("/auth")
	{
		authGroup.POST("/register", func(ctx *g.Context) { handleRegister(ctx, db, cfg) })
		authGroup.POST("/login", func(ctx *g.Context) { handleLogin(ctx, db, cfg) })
	}

	// User notes routes
	notesGroup := ge.Group("/notes").Use(AuthMiddleware(cfg))
	{
		notesGroup.GET("/", func(ctx *g.Context) { handleGetAllNotes(ctx, db) })
		notesGroup.POST("/create", func(ctx *g.Context) {})
		notesGroup.PUT("/update", func(ctx *g.Context) {})
		notesGroup.DELETE("/delete", func(ctx *g.Context) {})
	}
}

// Auth handlers

func handleRegister(ctx *g.Context, db *gorm.DB, cfg *config.Config) {
	user := models.User{}

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, g.H{"error": err.Error()})
		return
	}

	var count int64
	db.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, g.H{"error": "email has already taken"})
		return
	}

	hashedPassword := hasher.HashPassword(user.Password)
	user.Password = hashedPassword

	result := db.Create(&user)

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

func handleLogin(ctx *g.Context, db *gorm.DB, cfg *config.Config) {
	user := models.User{}

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, g.H{"error": err.Error()})
		return
	}

	var count int64
	db.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)
	if count == 0 {
		ctx.JSON(http.StatusBadRequest, g.H{"error": "invalid credentials"})
		return
	}

	passwordToVerify := user.Password
	db.First(&user, "email = ?", user.Email)

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

// Notes handlers

func handleGetAllNotes(ctx *g.Context, db *gorm.DB) {
	var notes []models.Note

	db.Find(&notes)

	ctx.JSON(http.StatusOK, notes)
}
