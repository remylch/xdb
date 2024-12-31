package api

import (
	"encoding/json"
	"log"
	"xdb/internal/store"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

const (
	DefaultAPIAddr = ":8080"
)

type NodeHttpServer struct {
	store *store.XDBStore
	addr  string
	app   *fiber.App
}

func NewHttpServer(store *store.XDBStore, addr string) *NodeHttpServer {
	if addr == "" {
		addr = DefaultAPIAddr
	}

	return &NodeHttpServer{
		store: store,
		addr:  addr,
		app:   fiber.New(),
	}
}

func (s *NodeHttpServer) healthcheckHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "healthy"})
}

func (s *NodeHttpServer) queryData(c *fiber.Ctx) error {
	query := c.Query("query")

	data, err := s.store.Get(query)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err})
	}

	var decodedData interface{}

	if err := json.Unmarshal(data, &decodedData); err != nil {
		decodedData = string(data)
	}

	return c.JSON(fiber.Map{
		"data": decodedData,
	})
}

func (s *NodeHttpServer) getCollectionsHandler(c *fiber.Ctx) error {
	collections := s.store.GetCollections()
	return c.JSON(fiber.Map{
		"location":    s.store.DataDir,
		"collections": collections,
	})
}

func (s *NodeHttpServer) getXDBApis(c *fiber.Ctx) error {
	apiPaths := make(map[string]interface{})

	for _, route := range s.app.GetRoutes() {
		if _, exist := apiPaths[route.Path]; !exist {
			apiPaths[route.Path] = map[string]interface{}{
				"method": route.Method,
				"params": route.Params,
			}
		}
	}

	return c.JSON(fiber.Map{
		"apis": apiPaths,
	})
}

func (s *NodeHttpServer) Start() error {
	s.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET",
	}))
	s.app.Get("/apis", s.getXDBApis)
	s.app.Get("/health", s.healthcheckHandler)
	s.app.Get("/collections", s.getCollectionsHandler)
	s.app.Get("/collections/data", s.queryData)
	log.Println("API listening on ", s.addr)
	return s.app.Listen(s.addr)
}
