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

	app.Run(iris.Addr("0.0.0.0:80"), iris.WithoutServerError(iris.ErrServerClosed))
}
