package main

import (
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	caching "logo-api/caching"
	"logo-api/emojipedia"
	"logo-api/image"
	. "logo-api/structs"
	"os"
	"strconv"
	"strings"
)

func main() {
	app := iris.New()
	app.Get("/", hello)
	app.Get("/generate", generate)

	// Generate cache folder if it does not exist
	if _, err := os.Stat("cache"); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Creating cache folder")
		err := os.Mkdir("cache", os.ModePerm)
		if err != nil {
		}
	}

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
	color := ctx.URLParamDefault("color", "ffffff")
	sizeParam := ctx.URLParamDefault("size", "256")
	emojiType := ctx.URLParamDefault("type", "apple")

	currentLogo := Logo{Emoji: emojiName, Color: color, Platform: emojiType, Size: sizeParam}

	cached := caching.IsCached(currentLogo)

	if cached {
		fmt.Println("Using cached version")
		_ = ctx.ServeFile(fmt.Sprintf("cache/%s.png", caching.GetName(currentLogo)))
		return
	}

	fmt.Println("Is cached?", cached)

	// There's probably a better way to do this, so please feel free to improve it!
	switch strings.ToLower(emojiType) {
	case "apple":
	case "android":
	case "discord":
	default:
		{
			showError(ctx, "Invalid emoji type. Accepted types: apple, android, discord")
			return
		}
	}

	if len(color) < 6 {
		showError(ctx, "Color must be in hex format e.g. ad5ff2")
		return
	}

	fmt.Println(fmt.Sprintf("Searching for emoji with name %s, type %s", emojiName, emojiType))
	emoji, err := emojipedia.Search(currentLogo)

	if err != nil {
		if err == emojipedia.ErrNoEmoji {
			showError(ctx, "Couldn't find an emoji with the search term provided")
			return
		} else if err == emojipedia.ErrNoUrl {
			showError(ctx, "Couldn't fetch image URL for emoji")
			return
		} else {
			fmt.Println(err)
		}
	}

	size, err := strconv.Atoi(sizeParam)
	err = image.Generate(emoji, currentLogo, size)

	if err != nil {
		fmt.Println(err)
		showError(ctx, "Error generating, check the URL and try again.")
		return
	} else {
		err = ctx.ServeFile(fmt.Sprintf("cache/%s.png", caching.GetName(currentLogo)))
	}
}
