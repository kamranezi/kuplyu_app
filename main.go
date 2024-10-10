package main

import (
    "kuplyu_app/config"
    "kuplyu_app/routes"
    "github.com/gin-gonic/gin"
)

func main() {
    // Подключение к базе данных
    config.ConnectDatabase()

    // Инициализация роутера
    router := gin.Default()

    // Инициализация маршрутов
    routes.InitRoutes(router)

    // Запуск сервера
    router.Run(":8080")
}
