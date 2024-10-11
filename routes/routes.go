package routes

import (
    "kuplyu_app/config"
    "kuplyu_app/models"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "github.com/golang-jwt/jwt"
    "net/http"
    "time"
    "fmt"
)

// Секретный ключ для подписи JWT токенов
var SecretKey = "your_secret_key"

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

    // Маршруты, которые требуют аутентификации
    authorized := router.Group("/")
    authorized.Use(AuthMiddleware())  // Применяем JWT middleware
    {
        authorized.GET("/requests", GetRequests)          // Получение всех заявок
        authorized.POST("/requests/create", CreateRequest) // Создание новой заявки
    }
}

// JWT middleware для защиты маршрутов
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")

        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is missing"})
            c.Abort()
            return
        }

        // Удаляем "Bearer " из строки токена
        tokenString = tokenString[len("Bearer "):]

        // Парсинг и проверка токена
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(SecretKey), nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
            c.Abort()
            return
        }



        c.Next()
    }
}

// Регистрация нового пользователя
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

// Логин пользователя и генерация JWT токена
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

    // Генерация JWT токена
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "exp":     time.Now().Add(time.Hour * 72).Unix(), // Токен действует 72 часа
    })

    // Подпись токена
    tokenString, err := token.SignedString([]byte(SecretKey))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token"})
        return
    }

    // Отправляем токен в ответе
    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// Получение всех заявок
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

// Создание новой заявки с использованием user_id из JWT
func CreateRequest(c *gin.Context) {
    var request models.Request

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
        return
    }


    // Сохранение заявки в базе данных
    if err := config.DB.Create(&request).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create request", "error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Request created successfully"})
}
