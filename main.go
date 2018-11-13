package main

import (
	"encoding/csv"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"io"
	"os"
	"strings"
)

func loadCSVData() map[string]string {
	csvData := make(map[string]string)

	file, err := os.Open("Corpus.csv")

	if err != nil {
		panic(fmt.Sprintf("%s", err))
	}

	defer file.Close()

	r := csv.NewReader(file)

	// skip first line
	record, err := r.Read()

	for record, err = r.Read(); err != io.EOF; record, err = r.Read() {
		if err != nil {
			panic(fmt.Sprintf("err in reading csv: %s", err))
		}

		csvData[strings.TrimSpace(record[0])] = strings.TrimSpace(record[1])
	}

	return csvData
}

type Record struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func main() {
	csvData := loadCSVData()

	app := iris.New()
	app.Logger().SetLevel("debug")

	app.Use(recover.New())
	app.Use(logger.New())

	app.Handle("GET", "/{key}", func(ctx iris.Context) {
		value, ok := csvData[ctx.Params().Get("key")]

		if ok {
			ctx.JSON(iris.Map{
				"key":   ctx.Params().Get("key"),
				"value": value,
			})
		} else {
			ctx.StatusCode(iris.StatusNotFound)
			ctx.JSON(iris.Map{
				"error": "Not Found!",
			})
		}
	})

	app.Handle("POST", "/", func(ctx iris.Context) {
		var r Record
		if err := ctx.ReadJSON(&r); err != nil {
			fmt.Println(err.Error())
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{
				"error": err.Error(),
			})
			return
		}

		// check if key already exists
		if _, ok := csvData[r.Key]; ok {
			ctx.StatusCode(iris.StatusNotAcceptable)
			ctx.JSON(iris.Map{
				"error": "Key already exists",
			})
			return
		}

		csvData[r.Key] = r.Value
		ctx.StatusCode(iris.StatusCreated)
		ctx.JSON(iris.Map{
			"key":   r.Key,
			"value": csvData[r.Key],
		})
	})

	app.Handle("PATCH", "/", func(ctx iris.Context) {
		var r Record

		if err := ctx.ReadJSON(&r); err != nil {
			fmt.Println(err.Error())
		}

		_, ok := csvData[r.Key]

		if !ok {
			ctx.StatusCode(iris.StatusNotFound)
			ctx.JSON(iris.Map{
				"error": fmt.Sprintf("Key %s not found!", r.Key),
			})
			return
		}

		csvData[r.Key] = r.Value
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(iris.Map{
			"key":   r.Key,
			"value": csvData[r.Key],
		})
	})

	app.Handle("DELETE", "/{key}", func(ctx iris.Context) {
		key := ctx.Params().Get("key")
		_, ok := csvData[key]

		if !ok {
			ctx.StatusCode(iris.StatusNotFound)
			ctx.JSON(iris.Map{
				"error": fmt.Sprintf("Key %s does not exists", key),
			})
			return
		}

		delete(csvData, key)
		ctx.StatusCode(iris.StatusNoContent)
		return
	})

	app.Run(iris.Addr("0.0.0.0:80"), iris.WithoutServerError(iris.ErrServerClosed))
}
