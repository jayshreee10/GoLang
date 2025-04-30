package controllers

import (
	"html/template"
	"net/http"
	"go-crud/models"
)

// Multiply is a helper function for the template
func multiply(a, b interface{}) float64 {
	var aVal, bVal float64
	
	switch a := a.(type) {
	case float64:
		aVal = a
	case int:
		aVal = float64(a)
	}
	
	switch b := b.(type) {
	case float64:
		bVal = b
	case int:
		bVal = float64(b)
	}
	
	return aVal * bVal
}

// Define the template with functions
var adminTmpl = template.Must(template.New("admin").Funcs(template.FuncMap{
	"multiply": multiply,
}).Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>SQLite Admin</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        h1, h2 { color: #333; }
        table { border-collapse: collapse; width: 100%; margin-bottom: 20px; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
        tr:nth-child(even) { background-color: #f9f9f9; }
    </style>
</head>
<body>
    <h1>SQLite Database Admin</h1>
    
    <h2>Users</h2>
    <table>
        <tr>
            <th>ID</th>
            <th>Email</th>
            <th>Created At</th>
        </tr>
        {{range .Users}}
        <tr>
            <td>{{.ID}}</td>
            <td>{{.Email}}</td>
            <td>{{.CreatedAt}}</td>
        </tr>
        {{end}}
    </table>
    
    <h2>Products</h2>
    <table>
        <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Status</th>
            <th>Price</th>
            <th>Created At</th>
            <th>Updated At</th>
        </tr>
        {{range .Products}}
        <tr>
            <td>{{.ID}}</td>
            <td>{{.Name}}</td>
            <td>{{.Status}}</td>
            <td>${{printf "%.2f" .Price}}</td>
            <td>{{.CreatedAt}}</td>
            <td>{{.UpdatedAt}}</td>
        </tr>
        {{end}}
    </table>
    
    <h2>Roles</h2>
    <table>
        <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Description</th>
            <th>Created At</th>
        </tr>
        {{range .Roles}}
        <tr>
            <td>{{.ID}}</td>
            <td>{{.Name}}</td>
            <td>{{.Description}}</td>
            <td>{{.CreatedAt}}</td>
        </tr>
        {{end}}
    </table>
    
    <h2>Orders</h2>
    <table>
        <tr>
            <th>ID</th>
            <th>User ID</th>
            <th>Address ID</th>
            <th>Total Amount</th>
            <th>Status</th>
            <th>Created At</th>
            <th>Updated At</th>
        </tr>
        {{range .Orders}}
        <tr>
            <td>{{.ID}}</td>
            <td>{{.UserID}}</td>
            <td>{{if .AddressID}}{{.AddressID}}{{else}}NULL{{end}}</td>
            <td>${{printf "%.2f" .TotalAmount}}</td>
            <td>{{.Status}}</td>
            <td>{{.CreatedAt}}</td>
            <td>{{.UpdatedAt}}</td>
        </tr>
        {{end}}
    </table>
    
    <h2>Order Items</h2>
    <table>
        <tr>
            <th>ID</th>
            <th>Order ID</th>
            <th>Product ID</th>
            <th>Product Name</th>
            <th>Quantity</th>
            <th>Price</th>
            <th>Total</th>
        </tr>
        {{range .OrderItems}}
        <tr>
            <td>{{.ID}}</td>
            <td>{{.OrderID}}</td>
            <td>{{.ProductID}}</td>
            <td>{{.Product.Name}}</td>
            <td>{{.Quantity}}</td>
            <td>${{printf "%.2f" .Price}}</td>
            <td>${{printf "%.2f" (multiply .Price .Quantity)}}</td>
        </tr>
        {{end}}
    </table>
    
    <h2>Addresses</h2>
    <table>
        <tr>
            <th>ID</th>
            <th>User ID</th>
            <th>Street Line 1</th>
            <th>Street Line 2</th>
            <th>City</th>
            <th>State</th>
            <th>Postal Code</th>
            <th>Country</th>
            <th>Default</th>
            <th>Created At</th>
            <th>Updated At</th>
        </tr>
        {{range .Addresses}}
        <tr>
            <td>{{.ID}}</td>
            <td>{{.UserID}}</td>
            <td>{{.StreetLine1}}</td>
            <td>{{.StreetLine2}}</td>
            <td>{{.City}}</td>
            <td>{{.State}}</td>
            <td>{{.PostalCode}}</td>
            <td>{{.Country}}</td>
            <td>{{if .IsDefault}}Yes{{else}}No{{end}}</td>
            <td>{{.CreatedAt}}</td>
            <td>{{.UpdatedAt}}</td>
        </tr>
        {{end}}
    </table>
</body>
</html>
`))

// AdminDashboard renders an admin dashboard to view database tables
func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	// Get users
	users, err := models.GetUsers(100)
	if err != nil {
		http.Error(w, "Error fetching users: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Get products
	products, err := models.GetProducts(100)
	if err != nil {
		http.Error(w, "Error fetching products: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Get roles
	roles, err := models.GetRoles(100)
	if err != nil {
		http.Error(w, "Error fetching roles: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Get orders
	orders, err := models.GetOrders(100)
	if err != nil {
		http.Error(w, "Error fetching orders: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Get all order items
	var allOrderItems []models.OrderItem
	for _, order := range orders {
		for _, item := range order.Items {
			allOrderItems = append(allOrderItems, item)
		}
	}
	
	// Get addresses
	addresses, err := models.GetAddresses(100)
	if err != nil {
		http.Error(w, "Error fetching addresses: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Render template
	data := struct {
		Users      []models.User
		Products   []models.Product
		Roles      []models.Role
		Orders     []models.Order
		OrderItems []models.OrderItem
		Addresses  []models.Address
	}{
		Users:      users,
		Products:   products,
		Roles:      roles,
		Orders:     orders,
		OrderItems: allOrderItems,
		Addresses:  addresses,
	}
	
	w.Header().Set("Content-Type", "text/html")
	adminTmpl.Execute(w, data)
}