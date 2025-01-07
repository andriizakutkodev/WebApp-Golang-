package routes

import (
	"net/http"
	"strings"
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
		token := getTokenFromContext(ctx)
		if token == "" || !isValidToken(token, cfg) {
			ctx.JSON(http.StatusUnauthorized, g.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func getTokenFromContext(ctx *g.Context) string {
	token := ctx.GetHeader("Authorization")

	if strings.HasPrefix(strings.ToLower(token), "bearer ") {
		token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))
	}

	return token
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
		notesGroup.GET("/", func(ctx *g.Context) { handleGetAllNotes(ctx, db, cfg) })
		notesGroup.POST("/create", func(ctx *g.Context) { handleCreateNoteForUser(ctx, db, cfg) })
		notesGroup.PUT("/update", func(ctx *g.Context) { handleUpdateNoteForUser(ctx, db, cfg) })
		notesGroup.DELETE("/delete", func(ctx *g.Context) { handleDeleteNoteForUser(ctx, db, cfg) })
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

func handleGetAllNotes(ctx *g.Context, db *gorm.DB, cfg *config.Config) {
	var notes []models.Note

	token := getTokenFromContext(ctx)
	userId := jwt.GetUserIdFromToken(token, cfg)

	db.Model(models.Note{}).Where("user_id = ?", userId).Find(&notes)

	ctx.JSON(http.StatusOK, notes)
}

func handleCreateNoteForUser(ctx *g.Context, db *gorm.DB, cfg *config.Config) {
	var note models.Note

	if err := ctx.ShouldBindJSON(&note); err != nil {
		ctx.JSON(http.StatusBadRequest, g.H{"error": err.Error()})
		return
	}

	note.UserID = jwt.GetUserIdFromToken(getTokenFromContext(ctx), cfg)

	result := db.Create(&note)

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusBadRequest, g.H{"error": result.Error})
		return
	}

	ctx.JSON(http.StatusOK, g.H{"note_id": note.ID})
}

func handleUpdateNoteForUser(ctx *g.Context, db *gorm.DB, cfg *config.Config) {
	var note models.Note
	var noteToUpdate models.Note

	if err := ctx.ShouldBindJSON(&note); err != nil {
		ctx.JSON(http.StatusBadRequest, g.H{"error": err.Error()})
		return
	}

	token := getTokenFromContext(ctx)
	userId := jwt.GetUserIdFromToken(token, cfg)

	if note.UserID != userId {
		ctx.JSON(http.StatusBadRequest, g.H{"error": "you don't have rights to update this note"})
		return
	}

	db.First(&noteToUpdate, note.ID)

	noteToUpdate.Title = note.Title
	noteToUpdate.Body = note.Body

	result := db.Save(&noteToUpdate)

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusBadRequest, g.H{"error": result.Error})
		return
	}

	ctx.JSON(http.StatusOK, g.H{"result": "note has been updated"})
}

func handleDeleteNoteForUser(ctx *g.Context, db *gorm.DB, cfg *config.Config) {
	noteId := ctx.DefaultQuery("note_id", "")

	if noteId == "" {
		ctx.JSON(http.StatusBadRequest, g.H{"error": "you should provide note id in query params"})
		return
	}

	var note models.Note

	db.First(&note, noteId)

	token := getTokenFromContext(ctx)
	userId := jwt.GetUserIdFromToken(token, cfg)

	if note.UserID != userId {
		ctx.JSON(http.StatusBadRequest, g.H{"error": "you don't have rights to delete this note"})
		return
	}

	result := db.Delete(&note)

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusBadRequest, g.H{"error": result.Error})
		return
	}

	ctx.JSON(http.StatusOK, g.H{"result": "note has been removed"})
}
