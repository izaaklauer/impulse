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

func (s *Server) listChambersHandler(ctx *gin.Context) {
    chambers, err := s.engine.List(ctx)
    if err != nil {
        ctx.AbortWithError(500, err)
        return
    }
    
    ctx.JSON(200, chambers)
}

func (s *Server) createChamberHandler(ctx *gin.Context) {
    var spec chamber.Spec
    if err := ctx.ShouldBindJSON(&spec); err != nil {
        ctx.AbortWithError(400, err)
        return
    }
    
    if err := s.engine.Create(ctx, spec); err != nil {
        ctx.AbortWithError(500, err)
        return
    }
    
    ctx.Status(201)
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
