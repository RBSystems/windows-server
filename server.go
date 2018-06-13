package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/byuoitav/common/events"
	"github.com/byuoitav/common/log"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	BIN_NAME          = "write"
	EVENT_ROUTER_NAME = "CP1"
	EVENT_ROUTER_PORT = "7000"
)

func main() {
	log.SetLevel("debug")

	// get hostname
	hostname, err := os.Hostname()
	if err != nil {
		log.L.Fatalf(fmt.Sprintf("failed to get hostname: %s", err))
	}

	// get building/room information
	building, room := "ITB", "1108"
	/* TODO add back in
	split := strings.Split(hostname, "-")
	if len(split) != 3 {
		log.L.Fatalf(fmt.Sprintf("hostname (%s) is invalid.", hostname))
	}
	building, room := split[0], split[1]
	*/

	// build event router address
	eventRouterAddress := fmt.Sprintf("%v-%v-%v:%v", building, room, EVENT_ROUTER_NAME, EVENT_ROUTER_PORT)

	// create event node
	filters := []string{}
	eventNode := events.NewEventNode(hostname, eventRouterAddress, filters)

	// router
	port := ":8000"
	router := echo.New()
	router.Pre(middleware.RemoveTrailingSlash())
	router.Use(middleware.CORS())

	// endpoints
	router.GET("/enable", func(c echo.Context) error {
		err := exec.Command(BIN_NAME, "--enable").Run()
		if err != nil {
			log.L.Errorf("command failed to execute: %s", err)
			//			eventNode.PublishEvent("error", events.Event{})
			return c.JSON(http.StatusOK, fmt.Sprintf("failed to enable: %s", err))
		}

		return c.JSON(http.StatusOK, "enabled")
	})

	router.GET("/disable", func(c echo.Context) error {
		err := exec.Command(BIN_NAME, "--disable").Run()
		if err != nil {
			log.L.Errorf("command failed to execute: %s", err)
			//			eventNode.PublishEvent("error", events.Event{})
			return c.JSON(http.StatusOK, fmt.Sprintf("failed to disable: %s", err))
		}

		return c.JSON(http.StatusOK, "disabled")
	})

	router.Start(port)
}
