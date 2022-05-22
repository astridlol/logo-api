package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"logo-api/emojipedia"
	"logo-api/image"
	"os"
	"strconv"
)

func main() {
	app := iris.New()
	app.Get("/", hello)
	app.Get("/generate", generate)

	// Listen for requests on port 8080
	_ = app.Listen(":8080")
	// check for errors
}

// Routes

func hello(ctx iris.Context) {
	ctx.JSON(iris.Map{
		"success": true,
	})
}

func showError(ctx iris.Context, err string) {
	_, err2 := ctx.JSON(iris.Map{
		"success": false,
		"err":     err,
	})
	if err2 != nil {
		return
	}
}

func generate(ctx iris.Context) {
	emojiName := ctx.URLParamDefault("emoji", "cookie")
	color := ctx.URLParamDefault("color", "#ffffff")
	sizeParam := ctx.URLParamDefault("size", "256")

	if len(color) < 6 {
		showError(ctx, "Color must be in hex format e.g. ad5ff2")
	}

	fmt.Println("Searching for emoji...")
	emoji, err := emojipedia.Search(emojiName)

	if err != nil {
		if err == emojipedia.ErrNoEmoji {
			showError(ctx, "Couldn't find an emoji with the search term provided")
		} else if err == emojipedia.ErrNoUrl {
			showError(ctx, "Couldn't fetch image URL for emoji")
		} else {
			fmt.Println(err)
		}
	}

	size, err := strconv.Atoi(sizeParam)

	fmt.Println("Please wait, generating image...")
	err = image.Generate(emoji, color, size)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = ctx.ServeFile("./output.png")

	fmt.Println("Logo saved to output.png!")
}
