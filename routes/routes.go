package routes

import (
    "kuplyu_app/config"
    "kuplyu_app/models"
    "github.com/gin-gonic/gin"
    "net/http"
    "github.com/gin-contrib/cors"
)


// Инициализация маршрутов
func InitRoutes(router *gin.Engine) {
    // Настройка CORS для междоменных запросов
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"}, // URL фронтенда
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
    }))

    // Маршруты для регистрации и логина
    router.POST("/register", Register)  // Регистрация
    router.POST("/login", Login)        // Логин

    // Маршруты для работы с заявками
    router.GET("/requests", GetRequests)    // Получение всех заявок
    router.POST("/requests/create", CreateRequest) // Создание новой заявки
}

// Регистрация
func Register(c *gin.Context) {
    var input models.User

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
        return
    }

    if err := config.DB.Create(&input).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to register user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

// Логин
func Login(c *gin.Context) {
    var input models.User
    var user models.User

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
        return
    }

    if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
        return
    }

    if user.Password != input.Password {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
        return
    }

    // Ответ с успешным входом и перенаправление на страницу заявок
    c.JSON(http.StatusOK, gin.H{"message": "Login successful", "redirect": "/requests"})
}

// Получение всех заявок из базы данных
func GetRequests(c *gin.Context) {
    var requests []models.Request
    c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
    c.Header("Pragma", "no-cache")
    // Получение всех заявок из базы данных
    if err := config.DB.Find(&requests).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get requests"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"requests": requests})
}

// Создание новой заявки
func CreateRequest(c *gin.Context) {
    var request models.Request

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
        return
    }

    // Сохранение заявки в базе данных
    if err := config.DB.Create(&request).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create request"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Request created successfully"})
}
