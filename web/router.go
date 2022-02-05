package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/web/handlers"
	"github.com/yzx9/otodo/web/middlewares"
)

func setupRouter(e *gin.Engine) {
	r := e.Group("/api")

	// Public routes
	{
		// Ping test
		r.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})

		// Auth
		r.POST("/sessions", handlers.PostSessionHandler)
		r.DELETE("/sessions", handlers.DeleteSessionHandler)

		r.POST("/sessions/token", handlers.PostSessionTokenHandler)
	}

	// Authorized routes
	r = r.Group("/", middlewares.JwtAuthMiddleware())
	{
		// Auth
		r.GET("/sessions", handlers.GetSessionHandler)

		// File
		// r.MaxMultipartMemory = MaxFileSize // 限制 Gin 上传文件时最大内存 (默认 32 MiB)
		r.POST("/files", handlers.PostFileHandler)
		r.GET("/files/:id", handlers.GetFileHandler)

		// Todo
		r.POST("/todos", handlers.PostTodoHandler)
		r.PUT("/todos/:id", handlers.PutTodoHandler)
		r.PATCH("/todos/:id", handlers.PatchTodoHandler)
		r.GET("/todos/:id", handlers.GetTodoHandler)
		r.DELETE("/todos/:id", handlers.DeleteTodoHanlder)

		r.GET("/todos/todo-lists/:id", handlers.GetTodosFromTodoListHandler)

		// Todo List
		r.GET("/todo-lists", handlers.GetTodoListsHandler)
	}
}
