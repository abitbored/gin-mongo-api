package controllers

import (
	"context"
	"gin-mongo-api/configs"
	"gin-mongo-api/models"
	"gin-mongo-api/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var bookCollection *mongo.Collection = configs.GetCollection(configs.DB, "books")
var validate = validator.New()

func CreateBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var book models.Book
		defer cancel()

		if err := c.BindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, responses.BookResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if validationErr := validate.Struct(&book); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.BookResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newBook := models.Book{
			ID:       primitive.NewObjectID(),
			Title:    book.Title,
			Author:   book.Author,
			Synopsis: book.Synopsis,
		}

		result, err := bookCollection.InsertOne(ctx, newBook)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.BookResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.BookResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		id := c.Param("id")
		var book models.Book
		defer cancel()

		objID, _ := primitive.ObjectIDFromHex(id)

		err := bookCollection.FindOne(ctx, bson.M{"id": objID}).Decode(&book)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.BookResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.BookResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": book}})
	}
}

func EditBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		id := c.Param("id")
		var book models.Book
		defer cancel()

		objID, _ := primitive.ObjectIDFromHex(id)

		if err := c.BindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, responses.BookResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if validationErr := validate.Struct(&book); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.BookResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{"title": book.Title, "author": book.Author, "Synopsis": book.Synopsis}
		result, err := bookCollection.UpdateOne(ctx, bson.M{"id": objID}, bson.M{"$set": update})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.BookResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		var updatedBook models.Book
		if result.MatchedCount == 1 {
			err := bookCollection.FindOne(ctx, bson.M{"id": objID}).Decode(&updatedBook)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.BookResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.BookResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedBook}})
	}
}

func DeleteBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		id := c.Param("id")
		defer cancel()

		objID, _ := primitive.ObjectIDFromHex(id)

		result, err := bookCollection.DeleteOne(ctx, bson.M{"id": objID})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.BookResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.BookResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Book with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			responses.BookResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "Book successfully deleted!"}},
		)
	}
}

func GetAllBooks() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var books []models.Book
		defer cancel()

		results, err := bookCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.BookResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleBook models.Book
			if err = results.Decode(&singleBook); err != nil {
				c.JSON(http.StatusInternalServerError, responses.BookResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			books = append(books, singleBook)
		}

		c.JSON(http.StatusOK,
			responses.BookResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": books}},
		)
	}
}
