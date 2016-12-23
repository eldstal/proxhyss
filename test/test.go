package main;

import (
  "github.com/eldstal/proxhyss"
  "os"
  "image"
  "image/png"
)

func main() {

  
  hats := proxhyss.InitHats("../hats")

  for _, file := range os.Args[1:] {
    infile, err := os.Open(os.Args[1])
      if err != nil {
        panic(err.Error())
      }
    defer infile.Close()

    img, _, err := image.Decode(infile)
    if err != nil {
      panic(err.Error())
    }

    hatted := hats.ApplyHats(img)

    outfile, err := os.Create(file + "_faces.png")
    if err != nil {
        panic(err.Error())
    }
    defer outfile.Close()
    png.Encode(outfile, hatted)
  }
}
