so
-------------

[English](README.md) | 中文

`so 这个命令要怎么写？`

你是否遇到过这样的情况，想使用一条命令完成某些事情，但却不记得参数了？

so 可以帮助你解决这个问题。



## 例子


so 怎样把 image.png 切分成 256x256 的图片？

```
convert image.png -crop 256x256 +repage output_%d.png
```

![](https://so.pingflash.com/demo/01-zh.gif)


so 用 ffmpeg 把这个目录下的图片拼接成视频，用 x265 编码。图片文件名是 image_0000001.png 这种格式的

```
ffmpeg -framerate 24 -i image_%07d.png -c:v libx265 -pix_fmt yuv420p output.mp4
```

![](https://so.pingflash.com/demo/02-zh.gif)


so 搜索并统计10天以前创建的文件的数量

```
find . -type f -ctime +10 | wc -l
```

![](https://so.pingflash.com/demo/03-zh.gif)


也可以让 AI 解释一下这条命令的参数都是什么意思：

![](https://so.pingflash.com/demo/04-zh.png)



## 安装

### Linux 用户

```
curl -sSL https://so.pingflash.com/install.sh | sh
```

### Mac 用户

因为我没有 Mac，没办法测试安装脚本，所以还请到 release 页面下载然后手动将二进制文件复制到你的 PATH 下。

## 原理

so 使用 AI 生成命令，你的疑问会被发送给 AI，然后 AI 会告诉你答案。



## 自定义

默认情况下，so 提供免费的 AI 服务，你也可以使用自己的 API.

```
so config
```

与 OpenAI 的 Chat 接口兼容的 API 均可使用。

为了避免滥用，免费的 AI 服务器限制每个 IP 每天 100 次请求。根据我个人的经验，这个数量是很充足的。如果有其他意见，请直接提 issue。

