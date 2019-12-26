package proxhyss;

import (
  "github.com/fogleman/gg"
  "github.com/nfnt/resize"
  "gocv.io/x/gocv"

  "fmt"
  "image"
  //"image/color"

  "path/filepath"
  "io/ioutil"
  "strings"
)

type HatDB struct {
  hats []image.Image
}

func InitHats(directory string) *HatDB {
  self := new(HatDB)

  files,_ := ioutil.ReadDir(directory)

  for _,file := range files {
    path := filepath.Join(directory, file.Name())

    if (!file.IsDir() &&
        strings.HasSuffix(path, ".png")) {
      pixmap,err := gg.LoadPNG(path)
      if err == nil {
        self.hats = append(self.hats, pixmap)
        fmt.Printf("Loaded hat %v\n", path)
      } else {
        fmt.Printf("Unable to load hat %v: %v\n", path, err)
      }
    }
  }

  return self
}

func (db HatDB) getScaledHat(width uint) (image.Image) {
  orig_hat := &db.hats[0];     // TODO: Select hat at random
  return resize.Resize(width, 0, *orig_hat, resize.NearestNeighbor)
}


func detectFaces(img image.Image) []image.Rectangle {

  xmlFile := "haarcascade_frontalface_default.xml"


  // prepare image matrix
  mat,err := gocv.ImageToMatRGBA(img)
  if err != nil {
    fmt.Printf("Error converting image to matrix")
    ret := make([]image.Rectangle, 0, 1)
    return ret
  }
  defer mat.Close()

  // load classifier to recognize faces
  classifier := gocv.NewCascadeClassifier()
  defer classifier.Close()

  if !classifier.Load(xmlFile) {
    fmt.Printf("Error reading cascade file: %v\n", xmlFile)
    ret := make([]image.Rectangle, 0, 1)
    return ret
  }

  // detect faces
  rects := classifier.DetectMultiScale(mat)
  fmt.Printf("found %d faces\n", len(rects))

  return rects
}

func (self HatDB) ApplyHats(img image.Image) image.Image {

    rects := detectFaces(img);

    draw := gg.NewContextForImage(img)
    draw.SetLineWidth(3)
    draw.SetRGB(0,1,0)

    for _, head := range rects {
      fmt.Printf("%v\n", head)

      /*
      draw.DrawRectangle(float64(head.Min.X),
                         float64(head.Min.Y),
                         float64(head.Max.X-head.Min.X),
                         float64(head.Max.Y-head.Min.Y))
      draw.Stroke()
      */

      // Select a hat and scale it to the head's size
      head_width := head.Max.X - head.Min.X
      pix := self.getScaledHat(uint(head_width * 3))

      // Fit the hat to the head
      draw.DrawImageAnchored(pix, (head.Min.X + head.Max.X)/2, head.Min.Y, 0.5, 0.5)


    }


  return draw.Image()
}


