package qrcode

import (
   code "github.com/skip2/go-qrcode"

)

/*
content表示要生成二维码的内容，可以是任意字符串。
level表示二维码的容错级别，取值有Low、Medium、High、Highest。
size表示生成图片的width和height，像素单位。
filename表示生成的文件名路径。
 */
func WriteFile(content string, level code.RecoveryLevel, size int, filename string) error  {
    return code.WriteFile(content,level,size,filename)

}


//有时候我们不想直接生成一个PNG文件存储，我们想对PNG图片做一些处理，比如缩放了，旋转了，或者网络传输了等，
//基于此，我们可以使用Encode函数，生成一个PNG 图片的字节流，这样我们就可以进行各种处理了。
//用法和WriteFile函数差不多，只不过返回的是一个[]byte字节数组，这样我们就可以对这个字节数组进行处理了。
func Encode(content string, level code.RecoveryLevel, size int) ([]byte, error){
    return code.Encode(content,level,size)

}





/*

自定义二维码
除了以上两种快捷方式，该库还为我们提供了对二维码的自定义方式，比如我们可以自定义二维码的前景色和背景色等。qrcode.New函数可以返回一个*QRCode，我们可以对*QRCode设置，实现对二维码的自定义。

比如我们设置背景色为绿色，前景色为白色的二维码


func main() {
	qr,err:=qrcode.New("http://www.flysnow.org/",qrcode.Medium)
	if err != nil {
		log.Fatal(err)
	} else {
		qr.BackgroundColor = color.RGBA{50,205,50,255}
		qr.ForegroundColor = color.White
		qr.WriteFile(256,"./blog_qrcode.png")
	}
}
指定*QRCode的BackgroundColor和ForegroundColor即可。然后调用WriteFile方法生成这个二维码文件。

func New(content string, level RecoveryLevel) (*QRCode, error)

// A QRCode represents a valid encoded QRCode.
type QRCode struct {
	// Original content encoded.
	Content string

	// QR Code type.
	Level         RecoveryLevel
	VersionNumber int

	// User settable drawing options.
	ForegroundColor color.Color
	BackgroundColor color.Color
}

func New(content string, level RecoveryLevel) (*QRCode, error)

// A QRCode represents a valid encoded QRCode.
type QRCode struct {
	// Original content encoded.
	Content string

	// QR Code type.
	Level         RecoveryLevel
	VersionNumber int

	// User settable drawing options.
	ForegroundColor color.Color
	BackgroundColor color.Color
}
 */
