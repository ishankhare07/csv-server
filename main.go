package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"io"
	"os"
	"strings"
)

func getCSVData(key string) (string, error) {
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

		if strings.TrimSpace(record[0]) == key {
			return strings.TrimSpace(record[1]), nil
		}
	}

	return "", errors.New("Not Found!")
}

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")

	app.Use(recover.New())
	app.Use(logger.New())

	app.Handle("GET", "/{key}", func(ctx iris.Context) {
		value, err := getCSVData(ctx.Params().Get("key"))

		if err != nil {
			ctx.StatusCode(iris.StatusNotFound)
			ctx.JSON(iris.Map{
				"error": err.Error(),
			})
		} else {
			ctx.JSON(iris.Map{
				"key":   ctx.Params().Get("key"),
				"value": value,
			})
		}
	})

	app.Run(iris.Addr("0.0.0.0:8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
