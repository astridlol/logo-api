package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"logo-api/emojipedia"
	"logo-api/image"
	"strconv"
	"strings"
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
	emojiType := ctx.URLParamDefault("type", "apple")

	// There's probably a better way to do this, so please feel free to improve it!
	switch strings.ToLower(emojiType) {
	case "apple":
	case "android":
	case "discord":
	default:
		{
			showError(ctx, "Invalid emoji type. Accepted types: Apple, android, discord")
			return
		}
	}

	if len(color) < 6 {
		showError(ctx, "Color must be in hex format e.g. ad5ff2")
	}

	fmt.Println(fmt.Sprintf("Searching for emoji with name %s, type %s", emojiName, emojiType))
	emoji, err := emojipedia.Search(emojiName, emojiType)

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
		showError(ctx, "Error generating, check the URL and try again.")
	} else {
		err = ctx.ServeFile("./output.png")

		fmt.Println("Logo saved to output.png!")
	}
}
