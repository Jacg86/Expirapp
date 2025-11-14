package server

/*
Este archivo contiene la implementación del servidor Gin
*/

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/*
# Server representa el servidor Gin
* httpServer: el servidor HTTP subyacente
* engine: el motor Gin
* db: la conexión a la base de datos
* config: la configuración del servidor
*/

type Server struct {
	httpServer *http.Server
	engine     *gin.Engine
	db         *gorm.DB
	config     Config
}

/*
# New crea una nueva instancia del servidor Gin
* db: la conexión a la base de datos
* cfg: la configuración del servidor
*/
func New(db *gorm.DB, cfg Config) *Server {
	gin.SetMode(cfg.Mode)

	engine := gin.New()

	/*
		# Crea una instancia de un servidorr gin
	*/
	server := &Server{
		engine: engine,
		db:     db,
		config: cfg,
		httpServer: &http.Server{
			Addr:         cfg.Port,
			Handler:      engine,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
	}

	server.setupMiddlewares()

	server.setupRouter()

	return server
}

func (s *Server) setupMiddlewares() {
	/*
		# evita que el servidor reviente
		# si hay un panic en alguna parte del código
	*/
	s.engine.Use(gin.Recovery())

	/*
		# registra las solicitudes entrantes
		# solo para depurar
	*/
	s.engine.Use(gin.Logger())
}

/*
# corsMiddleware maneja las solicitudes CORS
*/
func corsMiddleware() gin.HandlerFunc {
	/*
		# permite solicitudes CORS desde cualquier origen
		# con los métodos y encabezados especificados
	*/
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (s *Server) setupRouter() {

	s.engine.GET("/health", s.healthCheck)

	v1 := s.engine.Group("/api/v1")
	{
		v1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "API v1 is working!",
				"status":  "running",
			})
		})
	}

}

func (s *Server) healthCheck(c *gin.Context) {

	sqlDB, err := s.db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "database connection failed"})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "database ping failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "dabase is ok"})

}

func (s *Server) Start() error {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Starting server on %s", s.config.Port)

		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", s.config.Port, err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting")
	return nil
}

/*
# para debug
*/
func (s *Server) Engine() *gin.Engine {
	return s.engine
}
