package main;

import (
  "github.com/eldstal/proxhyss"
  "gocv.io/x/gocv"

  "fmt"

  "os"
  "image"
)

func main() {

  dir := "/tmp/pepparkak"
  if len(os.Args) > 1 {
    dir = os.Args[1]
  }

  err := os.Mkdir(dir, os.ModeDir | 0x755)
  if err != nil && !os.IsExist(err) {
    fmt.Println("Unable to create target directory.")
    return
  }

  hats := proxhyss.InitHats("hats")

  err = proxhyss.GifSetup()
  if err != nil {
    fmt.Println("Unable to create target directory.")
    return
  }

  tag := "dance"

  frames := make([]image.Image, 0, 1)
  new_frames := make([]image.Image, 0, 1)

  // Try new gifs until we find something to hat
  for {
    frames = proxhyss.GifGet(tag)
    new_frames = make([]image.Image, len(frames), len(frames))

    frames_with_hats := 0
    frames_without_hats := 0

    for i, frame := range(frames) {
      hatted,n := hats.ApplyHats(frame)
      new_frames[i] = hatted
      if (n == 0) {
        // No non-hat frames, please
        frames_without_hats += 1
      } else {
        frames_with_hats += 1
      }
    }

    // Only accept gifs where a majority of frames have hats
    if (frames_with_hats > frames_without_hats) { break }
    if (len(frames) == 0) { break }
  }

  window := gocv.NewWindow("Proxhyss")
	defer window.Close()

  // Export frames for flutting
  for _, hatted := range(new_frames) {
    mat,_ := gocv.ImageToMatRGBA(hatted)
    window.IMShow(mat)
    window.WaitKey(40)
  }
}
