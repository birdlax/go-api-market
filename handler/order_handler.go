package handler

import (
	"backend/domain"
	"backend/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}

}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var order domain.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if len(order.OrderItems) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Order must have at least 1 item"})
	}
	order.UserID = userID.(uint)
	order, err := h.service.CreateOrder(order)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"message":  "Order created successfully",
		"order_id": order.ID,
	})

}

func (h *OrderHandler) GetAllOrders(c *fiber.Ctx) error {
	orders, err := h.service.GetAllOrders()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(orders)

}

func (h *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	idUint64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid order ID"})
	}
	id := uint(idUint64)

	order, err := h.service.GetOrderByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Order not found"})
	}

	return c.JSON(order)
}

func (h *OrderHandler) UpdateOrder(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userIDUint, ok := userID.(uint)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user ID in context"})
	}

	idParam := c.Params("id")
	idUint64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}
	orderID := uint(idUint64)

	existingOrder, err := h.service.GetOrderByID(orderID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
	}

	if existingOrder.UserID != userIDUint {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You do not have permission to update this order"})
	}

	var order domain.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	updatedOrder, err := h.service.UpdateOrder(orderID, order)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(updatedOrder)
}

func (h *OrderHandler) DeleteOrder(c *fiber.Ctx) error {
	idParam := c.Params("id")
	idUint64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}
	id := uint(idUint64)

	if err := h.service.DeleteOrder(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Order deleted successfully"})
}

func (h *OrderHandler) MarkOrderAsPaid(c *fiber.Ctx) error {
	idParam := c.Params("id")
	idUint64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}
	id := uint(idUint64)

	err = h.service.MarkOrderAsPaid(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Order marked as paid"})
}
