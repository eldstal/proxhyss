package proxhyss;

//#cgo pkg-config: opencv
//#include "detect.hpp"
import "C"

import (
  "github.com/jeromelesaux/facedetection/facedetector"
  "github.com/fogleman/gg"
  "github.com/nfnt/resize"
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


func oldDetectFaces(img image.Image) []image.Rectangle {

    // Detect face rectangles (spelling error in API)
    f := facedetector.NewFaceDectectorFromImage(img)

    rects := f.GetFaces()

    var ret = make([]image.Rectangle, 0, len(rects))

    for _,r := range rects {
      ret = append(ret, image.Rect(r.X, r.Y, r.X+r.Width, r.Y+r.Height))
    }

    return ret
}


func detectFaces(img image.Image) []image.Rectangle {
  maxfaces := 10
  ret := make([]image.Rectangle, 0, maxfaces)

  found_faces := make([]C.struct_Face, maxfaces);

  C.TestFunc()
  nfaces := int(C.FindFaces(nil,
                C.int(img.Bounds().Dy()), C.int(img.Bounds().Dx()),
                &found_faces[0], C.int(len(found_faces))))


  for f:=0; f<nfaces; f++ {
    fmt.Printf("(%v,%v) (%v,%v)\n",
                found_faces[f].x1, found_faces[f].y1,
                found_faces[f].x2, found_faces[f].y2)
  }


  return ret
}

func (self HatDB) ApplyHats(img image.Image) image.Image {

    rects := oldDetectFaces(img);

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


