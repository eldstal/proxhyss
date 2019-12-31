package main;

import (
  "github.com/eldstal/proxhyss"
  "gocv.io/x/gocv"

  "fmt"
  "flag"

  "os"
  "image"
  "image/png"
  "image/gif"
  "path/filepath"
)

func main() {

  var export_gif = flag.Bool("export_gif", false, "Also export an animated GIF")
  var search_tag = flag.String("tag", "nope", "Search term for new GIFs")
  var show = flag.Bool("show", false, "Also show the generated frames in a window")
  var out_dir = flag.String("dir", "/tmp/pepparkak", "Directory to put output files in. Will be created.")
  var help = flag.Bool("h", false, "Show usage")

  flag.Parse()

  if (*help) {
    flag.PrintDefaults()
    return
  }

  err := os.Mkdir(*out_dir, os.ModeDir | 0x755)
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

  var GIF *gif.GIF
  frames := make([]image.Image, 0, 1)
  new_frames := make([]image.Image, 0, 1)

  // Try new gifs until we find something to hat
  for {
    frames,GIF = proxhyss.GifGet(*search_tag)
    new_frames = make([]image.Image, len(frames), len(frames))

    frames_with_hats := 0
    frames_without_hats := 0

    fmt.Printf("Finding faces in %d frames...\n", len(frames))

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

    fmt.Printf("%d frames with faces. %d frames without.\n", frames_with_hats, frames_without_hats)

    // Only accept gifs where a majority of frames have hats
    if (frames_with_hats > frames_without_hats) { break }
    if (len(frames) == 0) { break }
  }


  if (*show) {
    window := gocv.NewWindow("Proxhyss")

    frame_number := 0
    for _, hatted := range(new_frames) {

      if (*show) {
        mat,_ := gocv.ImageToMatRGBA(hatted)
        window.IMShow(mat)
        window.WaitKey(GIF.Delay[frame_number] * 10)  // frames 
      }
    }

    window.Close()
  }

  // Export frames for flutting
  frame_number := 0
  for _, hatted := range(new_frames) {

    filename := fmt.Sprintf("frame_%03d.png", frame_number)
    path := filepath.Join(*out_dir, filename)
    fmt.Printf("Exporting: %s\n", path)

    frame_number += 1

    f,err := os.OpenFile(path, os.O_RDWR | os.O_CREATE, 0644)
    if (err != nil) {
      fmt.Printf("Failed to open output file: %v\n", err);
      return
    }
    defer f.Close()

    png.Encode(f, hatted)

  }

  // Repack the animated GIF and export that too, for fun
  if (*export_gif) {
    proxhyss.GifRepack(GIF, new_frames)

    path := filepath.Join(*out_dir, "hats.gif")
    fmt.Printf("Exporting: %s\n", path)
    f,err := os.OpenFile(path, os.O_RDWR | os.O_CREATE, 0644)
    if (err != nil) {
      fmt.Printf("Failed to open output file: %v\n", err);
      return
    }
    defer f.Close()

    gif.EncodeAll(f, GIF)
  }
}
