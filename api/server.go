package api

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"xdb/store"
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

func (s *NodeHttpServer) getCollectionDetailsHandler(c *fiber.Ctx) error {
	collectionName := c.Params("name")

	//TODO: Extract this part to an handler or a service separated from the controller part
	data, err := s.store.Get(collectionName)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Collection not found"})
	}

	var decodedData interface{}
	if err := json.Unmarshal(data, &decodedData); err != nil {
		decodedData = string(data)
	}

	//---------------------------------------------------------------------

	return c.JSON(fiber.Map{
		"data":    decodedData,
		"indexes": make([]string, 0), //TODO: handle indexes
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
	s.app.Get("/collections/:name", s.getCollectionDetailsHandler)
	log.Println("API listening on ", s.addr)
	return s.app.Listen(s.addr)
}
