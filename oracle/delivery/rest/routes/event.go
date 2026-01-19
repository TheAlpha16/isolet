package routes

import (
	eventHan "github.com/TheAlpha16/isolet/oracle/delivery/rest/handler/event"
	"github.com/TheAlpha16/isolet/oracle/utils"

	"github.com/gofiber/fiber/v2"
)

func RegisterEvent(router fiber.Router, eventHandler eventHan.EventHandler) {
	eventROuter := router.Group(utils.RouteEvent)
	eventROuter.Get(utils.RouteEventInfo, eventHandler.Info)
}
