
so
-------------

English | [中文](README.zh.md)

`so How to do it in shell?`

Have you ever encountered a situation where you want to accomplish something with a command but can't remember the parameters? So can help you solve this problem.

## Examples

so crop image into 256x256 tile

```
convert image.png -crop 256x256 +repage output_%d.png
```

![](https://so.pingflash.com/demo/01-en.gif)

so make a video from images, the filename is like 000001.png

```
ffmpeg -framerate 24 -i %06d.png -c:v libx265 -pix_fmt yuv420p output.mp4
```

![](https://so.pingflash.com/demo/02-en.gif)

so find out how many files created 10 days ago

```
find . -type f -ctime +10 | wc -l
```
![](https://so.pingflash.com/demo/03-en.gif)

You can also let AI explain what the parameters of the command mean:

![](https://so.pingflash.com/demo/04-en.png)


## Installation

### For Linux Users

```
curl -sSL https://so.pingflash.com/install.sh | sh
```

### For Mac Users

Since I don't have a Mac, I couldn't test the installation script. Please go to the release page to download and manually copy the binary file to your PATH.

## How It Works

so uses AI to generate commands. Your questions will be sent to AI, and then AI will tell you the command.

## Customization

By default, so provides free AI services, but you can also use your own API.

```
so config
```

Any API compatible with OpenAI's Chat interface can be used.

To prevent abuse, the free AI server limits each IP to 100 requests per day. Based on my personal experience, this amount is quite sufficient. If you have any other suggestions, please feel free to open an issue.

