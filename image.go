package main

import (
	"encoding/base64"
	"github.com/nfnt/resize"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {

	Bigger()
	Decode()
}

func Decode() {
	//读原图片
	ff, _ := os.Open("test_resized.png")
	defer ff.Close()
	sourcebuffer := make([]byte, 500000)
	n, _ := ff.Read(sourcebuffer)
	//base64压缩
	sourcestring := base64.StdEncoding.EncodeToString(sourcebuffer[:n])

	//写入临时文件
	ioutil.WriteFile("a.png.txt", []byte(sourcestring), 0667)
	//读取临时文件
	cc, _ := ioutil.ReadFile("a.png.txt")

	//解压
	dist, _ := base64.StdEncoding.DecodeString(string(cc))
	//写入新文件
	f, _ := os.OpenFile("xx.png", os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer f.Close()
	f.Write(dist)
}

func Bigger() {
	// open "test.jpg"
	file, err := os.Open("a.png")
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(600, 0, img, resize.Lanczos3)

	out, err := os.Create("test_resized.png")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)
}

func Decoder() {
	data := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAKAAAAARBAMAAACod7rOAAAAG1BMVEX///8AAP/f3/+/v/9/f/9fX/8fH/8/P/+fn/83sGquAAACUElEQVQ4jdWTT1PbMBDFH7Js+Wi50HCUBya5OtOOuSadOuTohEI4Ovw/GmggRwHF6cfurhwSQr9Aq8x4HHn3p/d2V8C/v8TXkxTK0psX06NjplrrVNGjRX+V3vG0LoAfwCuHhxXnaL3gKOCAfsCt5lU2wMHts+ytgf45amyl/OUQKlU290pJwHPcJ4ZPKBgYbXUsAoqaIUmXyrwlEBnEbA3E6TvgoyFgWBEwLMNSVug7KXocLZRTKGPI6m8g1IEDRs3GCqgqtiwLsqlsYEkOMHdliubKKexaFum7U46WwJHWxxQvm7O3pxEDjQMSgBRmRkyrrpkA7XGTW4o4dkBxSWHt3geFA6MqZzlkyyKmjL4D9isHPOX36f4ZuDMjy5kiOnaWu65Dl2YTeHWBIwf0t52blcJ6wu3PqYW7cnRAzd9Bp24URknnpQxS6vs95fUS2izegOFDx0A9EyKpV0BXw8ybGaqh3xOxpA7HOMYwWCpUC3sm7gDJwpEQdQUMMuz9Gn2maZiwTgcM0iZMFqRQxH5BmukDzeIbsL2wrzUVL7Eq/QDMM3l9mlFceCbOGTiuoV542gj7PUhzHA4MAWk8GNhYbrWo7Sxu5jRuAL/R2LibcmJAXNG6qyGHXG1ZYpaYHO1dEJAGeL5WSGfe3bgQPBDwfQ2xBNZUEXGdejtkWVwUrl2/Pw2Q71/WY2ztXQg2wEDJwKCVZj08fckfJwUrJIebwPDebd6EBQHRd0X0p8WwuztLsbB01RK+0gR80paAPpmZe/rnM7g3fPqVxX+x/gC4nIDbXy+SzAAAAABJRU5ErkJggg=="
	// The actual image starts after the ","
	i := strings.Index(data, ",")
	if i < 0 {
		log.Fatal("no comma")
	}
	// pass reader to NewDecoder
	// base 64 数据
	//src := "iVBORw0KGgoAAAANSUhEUgAAAKAAAAARBAMAAACod7rOAAAAG1BMVEX///8AAP/f3/+/v/9/f/9fX/8fH/8/P/+fn/83sGquAAACUElEQVQ4jdWTT1PbMBDFH7Js+Wi50HCUBya5OtOOuSadOuTohEI4Ovw/GmggRwHF6cfurhwSQr9Aq8x4HHn3p/d2V8C/v8TXkxTK0psX06NjplrrVNGjRX+V3vG0LoAfwCuHhxXnaL3gKOCAfsCt5lU2wMHts+ytgf45amyl/OUQKlU290pJwHPcJ4ZPKBgYbXUsAoqaIUmXyrwlEBnEbA3E6TvgoyFgWBEwLMNSVug7KXocLZRTKGPI6m8g1IEDRs3GCqgqtiwLsqlsYEkOMHdliubKKexaFum7U46WwJHWxxQvm7O3pxEDjQMSgBRmRkyrrpkA7XGTW4o4dkBxSWHt3geFA6MqZzlkyyKmjL4D9isHPOX36f4ZuDMjy5kiOnaWu65Dl2YTeHWBIwf0t52blcJ6wu3PqYW7cnRAzd9Bp24URknnpQxS6vs95fUS2izegOFDx0A9EyKpV0BXw8ybGaqh3xOxpA7HOMYwWCpUC3sm7gDJwpEQdQUMMuz9Gn2maZiwTgcM0iZMFqRQxH5BmukDzeIbsL2wrzUVL7Eq/QDMM3l9mlFceCbOGTiuoV542gj7PUhzHA4MAWk8GNhYbrWo7Sxu5jRuAL/R2LibcmJAXNG6qyGHXG1ZYpaYHO1dEJAGeL5WSGfe3bgQPBDwfQ2xBNZUEXGdejtkWVwUrl2/Pw2Q71/WY2ztXQg2wEDJwKCVZj08fckfJwUrJIebwPDebd6EBQHRd0X0p8WwuztLsbB01RK+0gR80paAPpmZe/rnM7g3fPqVxX+x/gC4nIDbXy+SzAAAAABJRU5ErkJggg=="
	reader := strings.NewReader(data[i+1:])
	decoder := base64.NewDecoder(base64.StdEncoding, reader)
	// 以流式解码

	f, err := os.Create("a.png")
	if err != nil {
		log.Println(err)
	}
	bs, _ := ioutil.ReadAll(decoder)

	f.Write(bs)
}
