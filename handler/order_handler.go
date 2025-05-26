package handler

import (
	"backend/domain"
	"backend/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type OrderHandler struct {
	service domain.OrderService
}

func NewOrderHandler(service domain.OrderService) *OrderHandler {
	return &OrderHandler{service: service}

}

func (h *OrderHandler) GetOrder(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	utils.Logger.Println("🔄 [GetAllOrders] Start fetching unpaid orders for user")

	orders, err := h.service.GetUnpaidOrdersByUserID(userID)
	if err != nil {
		utils.Logger.Printf("❌ [GetAllOrders] Failed to fetch unpaid orders: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	utils.Logger.Println("✅ [GetAllOrders] Successfully fetched unpaid orders")
	return c.JSON(orders)
}
func (h *OrderHandler) GetOrdersByStatus(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	status := c.Query("status", "") // รับ status จาก query string เช่น ?status=pending

	utils.Logger.Printf("🔄 [GetOrdersByStatus] Fetching orders with status '%s'", status)

	orders, err := h.service.GetOrdersByUserIDAndStatus(userID, status)
	if err != nil {
		utils.Logger.Printf("❌ [GetOrdersByStatus] Failed to fetch orders: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	utils.Logger.Println("✅ [GetOrdersByStatus] Successfully fetched orders")
	return c.JSON(orders)
}

func (h *OrderHandler) GetAllOrders(c *fiber.Ctx) error {
	utils.Logger.Println("🔄 [GetAllOrders] Start fetching all orders")
	orders, err := h.service.GetAllOrders()
	if err != nil {
		utils.Logger.Printf("❌ [GetAllOrders] Failed to fetch orders: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	utils.Logger.Println("✅ [GetAllOrders] Successfully fetched all orders")
	return c.JSON(orders)
}

func (h *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	utils.Logger.Printf("🔄 [GetOrderByID] Start fetching order ID: %s", idParam)
	idUint64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.Logger.Printf("❌ [GetOrderByID] Invalid order ID format: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid order ID"})
	}
	id := uint(idUint64)

	order, err := h.service.GetOrderByID(id)
	if err != nil {
		utils.Logger.Printf("❌ [GetOrderByID] Order not found: %v", err)
		return c.Status(404).JSON(fiber.Map{"error": "Order not found"})
	}

	utils.Logger.Printf("✅ [GetOrderByID] Successfully fetched order ID: %d", id)
	return c.JSON(order)
}

func (h *OrderHandler) UpdateOrder(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		utils.Logger.Println("❌ [UpdateOrder] User ID not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userIDUint, ok := userID.(uint)
	if !ok {
		utils.Logger.Println("❌ [UpdateOrder] Invalid user ID in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user ID in context"})
	}

	idParam := c.Params("id")
	utils.Logger.Printf("🔄 [UpdateOrder] Start updating order ID: %s by user ID: %d", idParam, userIDUint)
	idUint64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.Logger.Printf("❌ [UpdateOrder] Invalid order ID format: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}
	orderID := uint(idUint64)

	existingOrder, err := h.service.GetOrderByID(orderID)
	if err != nil {
		utils.Logger.Printf("❌ [UpdateOrder] Order not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
	}

	if existingOrder.UserID != userIDUint {
		utils.Logger.Printf("❌ [UpdateOrder] User ID %d does not have permission to update order ID %d", userIDUint, orderID)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You do not have permission to update this order"})
	}

	var order domain.Order
	if err := c.BodyParser(&order); err != nil {
		utils.Logger.Printf("❌ [UpdateOrder] Invalid request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	updatedOrder, err := h.service.UpdateOrder(orderID, order)
	if err != nil {
		utils.Logger.Printf("❌ [UpdateOrder] Failed to update order ID %d: %v", orderID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	utils.Logger.Printf("✅ [UpdateOrder] Successfully updated order ID: %d", orderID)
	return c.JSON(updatedOrder)
}

func (h *OrderHandler) DeleteOrder(c *fiber.Ctx) error {
	idParam := c.Params("id")
	utils.Logger.Printf("🔄 [DeleteOrder] Start deleting order ID: %s", idParam)
	idUint64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.Logger.Printf("❌ [DeleteOrder] Invalid order ID format: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}
	id := uint(idUint64)

	if err := h.service.DeleteOrder(id); err != nil {
		utils.Logger.Printf("❌ [DeleteOrder] Failed to delete order ID %d: %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	utils.Logger.Printf("✅ [DeleteOrder] Successfully deleted order ID: %d", id)
	return c.JSON(fiber.Map{"message": "Order deleted successfully"})
}

func (h *OrderHandler) MarkOrderAsPaid(c *fiber.Ctx) error {
	// ดึง userID จาก context
	userIDValue := c.Locals("user_id")
	userID, ok := userIDValue.(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	utils.Logger.Printf("🔄 [MarkOrderAsPaid] Start marking order for user ID: %d as paid", userID)

	err := h.service.MarkOrderAsPaidByUserID(userID)
	if err != nil {
		utils.Logger.Printf("❌ [MarkOrderAsPaid] Failed to mark order as paid for user %d: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	utils.Logger.Printf("✅ [MarkOrderAsPaid] Successfully marked order as paid for user %d", userID)
	return c.JSON(fiber.Map{"message": "Order marked as paid"})
}

func (h *OrderHandler) CancelOrder(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	utils.Logger.Printf("🔄 [CancelOrder] User %d requested to cancel order", userID)

	err := h.service.CancelOrderByUserID(userID)
	if err != nil {
		utils.Logger.Printf("❌ [CancelOrder] Failed: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.Logger.Printf("✅ [CancelOrder] Successfully canceled order for user %d", userID)
	return c.JSON(fiber.Map{"message": "Order canceled successfully"})
}

func (h *OrderHandler) GetRevenueByCategory(c *fiber.Ctx) error {
	status := c.Query("status")
	results, err := h.service.GetRevenueByCategory(status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(results)
}
