package main

import "github.com/skip2/go-qrcode"

func main() {
    // 扫完二维码直接去这个网站.
    qrcode.WriteFile("http://www.flysnow.org/",qrcode.Medium,256,"./blog_qrcode.png")
}
