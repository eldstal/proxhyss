package main;

import (
  "github.com/eldstal/proxhyss"
  "gocv.io/x/gocv"

  "os"
  "image"
)

func main() {

  window := gocv.NewWindow("Proxhyss")
	defer window.Close()
  
  hats := proxhyss.InitHats("hats")

  for _, file := range os.Args[1:] {
    infile, err := os.Open(file)
    if err != nil {
      panic(err.Error())
    }
    defer infile.Close()

    img, _, err := image.Decode(infile)
    if err != nil {
      panic(err.Error())
    }

    hatted,_ := hats.ApplyHats(img)

    mat,_ := gocv.ImageToMatRGBA(hatted)
    window.IMShow(mat)
    window.WaitKey(2000)
  }
}
