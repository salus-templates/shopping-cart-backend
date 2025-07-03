// ProductsService/Program.cs

using ProductService;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
builder.Services.AddEndpointsApiExplorer();

// Configure CORS to allow requests from any origin
builder.Services.AddCors(options =>
{
    options.AddPolicy("AllowAnyOrigin",
        policyBuilder =>
        {
            policyBuilder.AllowAnyOrigin() // Allow any host to access
                   .AllowAnyHeader()
                   .AllowAnyMethod();
        });
});

var app = builder.Build();

// Use CORS policy
app.UseCors("AllowAnyOrigin"); // Apply the new policy


// Static list of products (moved from React app)
var products = new List<Product>
{
    new Product { Id = "prod1", Name = "Wireless Headphones", Price = 99.99m, ImageUrl = "https://placehold.co/300x200/FFD700/000000?text=Headphones", Description = "High-fidelity sound with noise cancellation." },
    new Product { Id = "prod2", Name = "Smartwatch", Price = 199.99m, ImageUrl = "https://placehold.co/300x200/87CEEB/000000?text=Smartwatch", Description = "Track your fitness and receive notifications." },
    new Product { Id = "prod3", Name = "Portable Bluetooth Speaker", Price = 49.99m, ImageUrl = "https://placehold.co/300x200/98FB98/000000?text=Speaker", Description = "Compact and powerful sound on the go." },
    new Product { Id = "prod4", Name = "Ergonomic Office Chair", Price = 249.99m, ImageUrl = "https://placehold.co/300x200/DDA0DD/000000?text=Chair", Description = "Comfortable and supportive for long working hours." },
    new Product { Id = "prod5", Name = "4K UHD Monitor", Price = 399.99m, ImageUrl = "https://placehold.co/300x200/F08080/000000?text=Monitor", Description = "Stunning visuals for work and entertainment." },
    new Product { Id = "prod6", Name = "Gaming Keyboard", Price = 79.99m, ImageUrl = "https://placehold.co/300x200/ADD8E6/000000?text=Keyboard", Description = "Mechanical keyboard with RGB lighting." },
    new Product { Id = "prod7", Name = "Gaming Mouse", Price = 39.99m, ImageUrl = "https://placehold.co/300x200/FFB6C1/000000?text=Mouse", Description = "High-precision sensor for competitive gaming." },
    new Product { Id = "prod8", Name = "Webcam 1080p", Price = 59.99m, ImageUrl = "https://placehold.co/300x200/DAA520/000000?text=Webcam", Description = "Full HD video calls and streaming." },
    new Product { Id = "prod9", Name = "External SSD 1TB", Price = 129.99m, ImageUrl = "https://placehold.co/300x200/B0C4DE/000000?text=SSD", Description = "Fast and portable storage solution." },
    new Product { Id = "prod10", Name = "USB-C Hub", Price = 29.99m, ImageUrl = "https://placehold.co/300x200/F4A460/000000?text=USB-C+Hub", Description = "Expand your laptop's connectivity with multiple ports." },
    new Product { Id = "prod11", Name = "Noise-Cancelling Earbuds", Price = 129.99m, ImageUrl = "https://placehold.co/300x200/C0C0C0/000000?text=Earbuds", Description = "Compact earbuds with active noise cancellation." },
    new Product { Id = "prod12", Name = "Smart Home Hub", Price = 89.99m, ImageUrl = "https://placehold.co/300x200/D8BFD8/000000?text=Smart+Hub", Description = "Control all your smart devices from one place." },
    new Product { Id = "prod13", Name = "Robot Vacuum Cleaner", Price = 299.99m, ImageUrl = "https://placehold.co/300x200/AFEEEE/000000?text=Vacuum", Description = "Automated cleaning for a spotless home." },
    new Product { Id = "prod14", Name = "Digital Camera", Price = 499.99m, ImageUrl = "https://placehold.co/300x200/F5DEB3/000000?text=Camera", Description = "Capture stunning photos and videos." },
    new Product { Id = "prod15", Name = "Portable Projector", Price = 179.99m, ImageUrl = "https://placehold.co/300x200/9ACD32/000000?text=Projector", Description = "Enjoy movies anywhere with a compact projector." },
    new Product { Id = "prod16", Name = "Fitness Tracker", Price = 69.99m, ImageUrl = "https://placehold.co/300x200/FFA07A/000000?text=Fitness+Tracker", Description = "Monitor your activity, heart rate, and sleep." },
    new Product { Id = "prod17", Name = "Electric Toothbrush", Price = 45.99m, ImageUrl = "https://placehold.co/300x200/BDB76B/000000?text=Toothbrush", Description = "Advanced cleaning for healthier gums." },
    new Product { Id = "prod18", Name = "Air Fryer", Price = 89.99m, ImageUrl = "https://placehold.co/300x200/E0FFFF/000000?text=Air+Fryer", Description = "Cook healthier meals with less oil." },
    new Product { Id = "prod19", Name = "Coffee Maker", Price = 75.99m, ImageUrl = "https://placehold.co/300x200/D2B48C/000000?text=Coffee+Maker", Description = "Brew your perfect cup of coffee every morning." },
    new Product { Id = "prod20", Name = "Smart Light Bulbs (2-pack)", Price = 25.99m, ImageUrl = "https://placehold.co/300x200/F0F8FF/000000?text=Smart+Bulbs", Description = "Control your lighting with your voice or app." },
};

// Endpoint to get all products
app.MapGet("/all-products", () => products)
    .WithName("GetAllProducts");

app.Run("http://0.0.0.0:8080");
