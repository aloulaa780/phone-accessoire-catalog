

package main  

import (  
    "context"  
    "encoding/json"  
    "github.com/gofiber/fiber/v2"  
    "github.com/jackc/pgx/v5/pgxpool"  
    "log"  
    "net/http"  
)  

var db *pgxpool.Pool  

type Product struct {  
    ID          int     `json:"id"`  
    Name        string  `json:"name"`  
    Description string  `json:"description"`  
    Price       float64 `json:"price"`  
}  

func main() {  
    app := fiber.New()  

    var err error  
    // Remplacez les valeurs dans la chaîne de connexion par les vôtres  
    db, err = pgxpool.Connect(context.Background(), "postgres://username:password@localhost:5432/yourdbname")  
    if err != nil {  
        log.Fatal(err)  
    }  
    defer db.Close()  

    // Routes  
    app.Get("/products", getProducts)                     // Récupérer tous les produits  
    app.Get("/products/:id", getProduct)                  // Récupérer un produit par ID  
    app.Post("/products", createProduct)                  // Créer un nouveau produit  
    app.Put("/products/:id", updateProduct)               // Mettre à jour un produit  
    app.Delete("/products/:id", deleteProduct)            // Supprimer un produit  

    log.Fatal(app.Listen(":3000"))  
}  

// Récupérer tous les produits  
func getProducts(c *fiber.Ctx) error {  
    rows, err := db.Query(context.Background(), "SELECT * FROM products")  
    if err != nil {  
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())  
    }  
    defer rows.Close()  

    products := []Product{}  
    for rows.Next() {  
        var prod Product  
        if err := rows.Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price); err != nil {  
            return err  
        }  
        products = append(products, prod)  
    }  
    return c.JSON(products)  
}  

// Récupérer un produit par ID  
func getProduct(c *fiber.Ctx) error {  
    id := c.Params("id")  
    var prod Product  
    err := db.QueryRow(context.Background(), "SELECT * FROM products WHERE id=$1", id).Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price)  
    if err != nil {  
        return c.Status(http.StatusNotFound).SendString(err.Error())  
    }  
    return c.JSON(prod)  
}  

// Créer un nouveau produit  
func createProduct(c *fiber.Ctx) error {  
    var prod Product  
    if err := c.BodyParser(&prod); err != nil {  
        return c.Status(http.StatusBadRequest).SendString(err.Error())  
    }  

    // Insertion dans la base de données  
    _, err := db.Exec(context.Background(), "INSERT INTO products (name, description, price) VALUES ($1, $2, $3)", prod.Name, prod.Description, prod.Price)  
    if err != nil {  
        return c.Status(http.StatusInternalServerError).SendString(err.Error())  
    }  
    // Répondre avec le produit créé  
    return c.Status(http.StatusCreated).JSON(prod)  
}  

// Mettre à jour un produit  
func updateProduct(c *fiber.Ctx) error {  
    id := c.Params("id")  
    var prod Product  
    if err := c.BodyParser(&prod); err != nil {  
        return c.Status(http.StatusBadRequest).SendString(err.Error())  
    }  

    // Mise à jour dans la base de données  
    _, err := db.Exec(context.Background(), "UPDATE products SET name=$1, description=$2, price=$3 WHERE id=$4", prod.Name, prod.Description, prod.Price, id)  
    if err != nil {  
        return c.Status(http.StatusInternalServerError).SendString(err.Error())  
    }  
    return c.JSON(prod)  
}  

// Supprimer un produit  
func deleteProduct(c *fiber.Ctx) error {  
    id := c.Params("id")  
    _, err := db.Exec(context.Background(), "DELETE FROM products WHERE id=$1", id)  
    if err != nil {  
        return c.Status(http.StatusInternalServerError).SendString(err.Error())  
    }  
    return c.SendStatus(http.StatusNoContent)  
}