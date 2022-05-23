# logo-api

An API to generate logos very quickly.

The base of this API is made from Alex's [logo generator CLI](https://gitlab.com/honour/logo-generator), I took it and made an API for it, while adding some needed features.

## How to use

Generation URL is `https://logo.jaack.host/generate` <br>

**Query params are required for the API to work** <br><br>
`emoji` - The emoji to use for this logo, must be a valid emoji on [emojipedia](https://emojipedia.org/) <br>
`color` - The hex color to use as the background, without the `#`. (Optional, defaults to white / ffffff)  <br>
`type` - The platform the emoji is from. Accepted platforms: Apple, Android, Discord. (Optional, defaults to Apple)

## Examples

`https://logo.jaack.host/generate?emoji=cookie` <br>
`https://logo.jaack.host/generate?emoji=open book&color=4287f5`

## Download

You can choose to download the API and run it locally if you would like. Click [here](https://github.com/astridlol/logo-api/releases) to find all the releases.
