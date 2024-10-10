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
        AllowOrigins:     []string{"http://localhost:3000"}, // Укажи URL фронтенда
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
    }))

    // Теперь все маршруты начинаются с /users
    users := router.Group("/users")
    {
        users.POST("/register", Register)  // Регистрация
        users.POST("/login", Login)        // Логин
    }
}

func Register(c *gin.Context) {
    var input models.User

    // Привязываем JSON из запроса к структуре User
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
        return
    }

    // Сохранение пользователя в базе данных
    if err := config.DB.Create(&input).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to register user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func Login(c *gin.Context) {
    var input models.User
    var user models.User

    // Привязываем JSON из запроса к структуре User
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
        return
    }

    // Ищем пользователя по email
    if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
        return
    }

    // Проверяем пароль
    if user.Password != input.Password {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
        return
    }

    // Если всё верно
    c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
