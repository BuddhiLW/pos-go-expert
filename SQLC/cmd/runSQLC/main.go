package main

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	ctx := context.Background()
	dbConn, err := sql.Open("mysql", "root:root@tcp(localhost:3333)/courses")

	if err != nil {
		panic(err)
	}
	defer dbConn.Close()

	// queries := db.New(dbConn)
	// Create
	// err = queries.CreateCategory(ctx, db.CreateCategoryParams{
	// 	ID:   uuid.New().String(),
	// 	Name: "Backend",
	// 	Description: sql.NullString{
	// 		String: "Backend description",
	// 		Valid:  true,
	// 	},
	// })

	// if err != nil {
	// 	panic(err)
	// }

	// categories, err := queries.ListCategories(ctx)
	// if err != nil {
	// 	panic(err)
	// }

	// for _, category := range categories {
	// 	println(category.ID, category.Name)
	// }

	// UPDATE
	// err = queries.UpdateCategory(ctx, db.UpdateCategoryParams{
	// 	ID:   "81f8ff0d-3c12-4717-8800-db13ab6a1f31",
	// 	Name: "Backend updated",
	// 	Description: sql.NullString{
	// 		String: "Backend description updated",
	// 		Valid:  true,
	// 	},
	// })

	// DELETE
	// err = queries.DeleteCategory(ctx, "81f8ff0d-3c12-4717-8800-db13ab6a1f31")

	// categories, err := queries.ListCategories(ctx)
	// if err != nil {
	// 	panic(err)
	// }

	// for _, category := range categories {
	// 	println(category.ID, category.Name)
	// }
}
