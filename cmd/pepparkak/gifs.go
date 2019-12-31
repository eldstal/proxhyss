package proxhyss

import (
  "fmt"
  "github.com/sanzaru/go-giphy"
  "image"
  "os"
  "encoding/json"
  "net/http"
  "io"
  //"io/ioutil"
  "image/gif"
  "image/draw"
)

type GiphyConfig struct {
    Key string    `json:"api_key"`
}

var CONFIG *GiphyConfig
var GIPHY *libgiphy.Giphy

func GifSetup() error {
  file, err := os.Open("config/giphy.json")
  if err != nil {  return err }

  config := &GiphyConfig{}

  decoder := json.NewDecoder(file)
  err = decoder.Decode(config)
  if err != nil { return err }

  CONFIG = config
  GIPHY = libgiphy.NewGiphy(CONFIG.Key)
  return nil
}

func renderOnPrevious(src image.Paletted, old *image.RGBA, bgindex uint8) *image.RGBA {
  rgba := image.NewRGBA(src.Rect)
  for x := 0; x < src.Rect.Dx(); x++ {
    for y := 0; y < src.Rect.Dy(); y++ {
      idx := src.ColorIndexAt(x,y)
      if (idx == bgindex) {
        rgba.Set(x, y, old.At(x,y))
      } else {
        rgba.Set(x, y, src.Palette[idx])
      }
    }
  }
  return rgba
}

func renderFully(src image.Paletted) *image.RGBA {
  rgba := image.NewRGBA(src.Rect)
  for x := 0; x < src.Rect.Dx(); x++ {
    for y := 0; y < src.Rect.Dy(); y++ {
      idx := src.ColorIndexAt(x,y)
      rgba.Set(x, y, src.Palette[idx])
    }
  }
  return rgba
}

func GifRender(GIF *gif.GIF) []image.Image {

  ret := make([]image.Image,len(GIF.Image), len(GIF.Image))

  overpaintImage := image.NewRGBA(GIF.Image[0].Rect)
  draw.Draw(overpaintImage, overpaintImage.Bounds(), GIF.Image[0], image.ZP, draw.Src)

  for i, srcImg := range GIF.Image {
      draw.Draw(overpaintImage, overpaintImage.Bounds(), srcImg, image.ZP, draw.Over)

      frame := image.NewRGBA(GIF.Image[0].Rect)
      draw.Draw(frame, frame.Bounds(), overpaintImage, image.ZP, draw.Over)
      ret[i] = frame
  }

  return ret

}

func GifGet(tag string) []image.Image {

  var src io.Reader
  ret := make([]image.Image,0,1)

  if (true) {
    metadata, err := GIPHY.GetRandom(tag)
    if err != nil {
      fmt.Println("Giphy error:", err)
      return ret
    }

    url := metadata.Data.Image_original_url
    fmt.Printf("Downloading %v\n", url)

    resp, err := http.Get(url)
    if err != nil {
      fmt.Println("Unable to download %+v", url)
      return ret
    }
    defer resp.Body.Close()
    src = resp.Body
  } else {
    file, err := os.Open("/tmp/giphy.gif")
    if err != nil {
      fmt.Println("Unable to load /tmp/giphy.gif")
      return ret
    }
    defer file.Close()
    src = file
  }


  GIF,err := gif.DecodeAll(src)
  if err != nil {
    fmt.Println("Unable to decode GIF")
    return ret
  }

  ret = GifRender(GIF)

  return ret

}
