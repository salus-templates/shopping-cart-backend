namespace ProductService;

// Define the Product model
public class Product
{
    public string Id { get; set; }
    public string Name { get; set; }
    public decimal Price { get; set; }
    public string ImageUrl { get; set; }
    public string Description { get; set; }
    public int Stock { get; set; }
}

// Order item received from the client
public class OrderItemRequest
{
    public string Id { get; set; }
    public int Quantity { get; set; }
}

// Full order request from the client
public class PlaceOrderRequest
{
    public List<OrderItemRequest> Items { get; set; }
    public decimal TotalAmount { get; set; }
    public string DeliveryAddress { get; set; }
    public string OrderDate { get; set; }
}

// Response for order placement
public class PlaceOrderResponse
{
    public bool Success { get; set; }
    public string Message { get; set; }
    public string OrderId { get; set; }
    public List<string> OutOfStockItems { get; set; } // New: List of items that caused failure
}
