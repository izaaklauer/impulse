package server

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "impulse/chamber"
    "impulse/engine"
    "log"
)

type Server struct{
    engine engine.Engine
}

func NewServer() (*Server, error) {
    engine, err := engine.NewContainerEngine() 
    if err != nil {
       return nil, fmt.Errorf("filed to instantiate engine: %v", err) 
    }
    
    return &Server{
        engine: engine,
    }, nil
}

func (s *Server) listChambersHandler(c *gin.Context) {
    chambers, err := s.engine.List()
    if err != nil {
        c.AbortWithError(500, err)
        return
    }
    
    c.JSON(200, chambers)
}

func (s *Server) createChamberHandler(c *gin.Context) {
    var spec chamber.Spec
    if err := c.ShouldBindJSON(&spec); err != nil {
        c.AbortWithError(400, err)
        return
    }
    
    if err := s.engine.Create(spec); err != nil {
        c.AbortWithError(500, err)
        return
    }
    
    c.Status(201)
}

func (s *Server) Serve() error {
    r := gin.Default()
    
    r.GET("/chambers", s.listChambersHandler)
    r.POST("/chambers", s.createChamberHandler)
    
    if err := r.Run(); err != nil {
        return fmt.Errorf("failed to run server: %v", err)
    }
    log.Printf("Server exiting")
    return nil
}
