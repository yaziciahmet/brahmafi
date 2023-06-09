package api

import (
	"brahma/common/config"
	"brahma/common/logger"
	"brahma/internal/core"

	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
)

type Server struct {
	log           logger.Logger
	app           *fiber.App
	brahmaService *core.BrahmaService
	port          string
}

func NewServer(log logger.Logger, brahmaService *core.BrahmaService, config *config.ApiConfig) *Server {
	app := fiber.New()
	server := &Server{
		log:           log,
		app:           app,
		brahmaService: brahmaService,
		port:          config.Port,
	}

	app.Use(
		fiberLogger.New(fiberLogger.Config{
			Output: log.GetWriter(),
		}),
	)

	app.Get("/v1/api/pool/:poolId", server.getPoolSnapshotByBlockHandler)
	app.Get("/v1/api/pool/:poolId/historic", server.getPoolSnapshotsHandler)

	return server
}

func (s *Server) Listen() error {
	return s.app.Listen(s.port)
}

func (s *Server) getPoolSnapshotByBlockHandler(c *fiber.Ctx) error {
	poolId := c.Params("poolId")
	if poolId == "" {
		return c.JSON(ErrorResponse{Message: "missing pool id"})
	}

	block := c.Query("block")
	if block == "" {
		block = "latest"
	}

	snapshot, err := s.brahmaService.GetPoolSnapshotByBlock(poolId, block)
	if err != nil {
		s.log.Error("error occurred on GetPoolSnapshotByBlock", "err", err)
		return c.JSON(ErrorResponse{Message: "internal server error"})
	}

	return c.JSON(SuccessResponse{Data: snapshot})
}

func (s *Server) getPoolSnapshotsHandler(c *fiber.Ctx) error {
	poolId := c.Params("poolId")
	if poolId == "" {
		return c.JSON(ErrorResponse{Message: "missing pool id"})
	}

	snapshots, err := s.brahmaService.GetPoolSnapshots(poolId)
	if err != nil {
		s.log.Error("error occurred on GetPoolSnapshots", "err", err)
		return c.JSON(ErrorResponse{Message: "internal server error"})
	}

	return c.JSON(SuccessResponse{Data: snapshots})
}
