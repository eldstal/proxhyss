package proxhyss;

import (
  "github.com/fogleman/gg"
  "github.com/nfnt/resize"
  "gocv.io/x/gocv"

  "fmt"
  "image"
  "image/color"
  "image/draw"

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

func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func detectFaces_dnn(img image.Image) []image.Rectangle {

	ret := make([]image.Rectangle, 0, 1)

	proto := "params/deploy.prototxt"
	model := "params/res10_300x300_ssd_iter_140000_fp16.caffemodel"

	mat,_ := gocv.ImageToMatRGB(img)

	// open DNN classifier
	net := gocv.ReadNet(proto, model)
	if net.Empty() {
		fmt.Printf("Error reading network model from : %v %v\n", proto, model)
		return ret
	}
	defer net.Close()

	//green := color.RGBA{0, 255, 0, 0}

	W := float32(mat.Cols())
	H := float32(mat.Rows())

	// convert image Mat to 300x300 blob that the detector can analyze
	blob := gocv.BlobFromImage(mat, 1.0,
                             image.Pt(300,300),
                             gocv.NewScalar(104, 177, 123, 0),
														 false, false)
	defer blob.Close()

	// feed the blob into the classifier
	net.SetInput(blob, "")

	// run a forward pass through the network
	detBlob := net.Forward("")
	defer detBlob.Close()

	// extract the detections.
	// for each object detected, there will be 7 float features:
	// objid, classid, confidence, left, top, right, bottom.
	detections := gocv.GetBlobChannel(detBlob, 0, 0)
	defer detections.Close()

	var confidence_threshold float32 = 0.5
	n_good_detections := 0

	for r := 0; r < detections.Rows(); r++ {
		confidence := detections.GetFloatAt(r, 2)
		if confidence < confidence_threshold {
			continue
		}
		n_good_detections += 1
	}

	ret = make([]image.Rectangle, n_good_detections, n_good_detections)
	index := 0

	for r := 0; r < detections.Rows(); r++ {
		confidence := detections.GetFloatAt(r, 2)
		if confidence < confidence_threshold {
			continue
		}

		left := detections.GetFloatAt(r, 3) * W
		top := detections.GetFloatAt(r, 4) * H
		right := detections.GetFloatAt(r, 5) * W
		bottom := detections.GetFloatAt(r, 6) * H

		// scale to video size:
		left = min(max(0, left), W-1)
		right = min(max(0, right), W-1)
		bottom = min(max(0, bottom), H-1)
		top = min(max(0, top), H-1)

		// draw it
		ret[index] = image.Rect(int(left), int(top), int(right), int(bottom))
		index += 0
	}

	return ret

}


func detectFaces_legacy(img image.Image) []image.Rectangle {

  xmlFile := "params/haarcascade_frontalface_default.xml"


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

func (self HatDB) ApplyHatsColorMatched(img image.Image, color_match_hats bool, color_match_palette *color.Palette) (image.Image, int) {

    // This DNN detector seems to work better than the haarcascade one
    rects := detectFaces_dnn(img);

    //old_rects := detectFaces_legacy(img);
    //fmt.Printf("old: %v, new: %v faces\n", len(old_rects), len(rects))


    draw_context := gg.NewContextForImage(img)
    draw_context.SetLineWidth(3)
    draw_context.SetRGB(0,1,0)

    for _, head := range rects {
      //fmt.Printf("%v\n", head)

      // Select a hat and scale it to the head's size
      head_width := head.Max.X - head.Min.X
      pix := self.getScaledHat(uint(head_width * 3))

      // Force the hat to use the colors in the color match palette
      if (color_match_hats) {
        matched_hat := image.NewPaletted(pix.Bounds(), *color_match_palette)
        draw.Draw(matched_hat, matched_hat.Bounds(), pix, image.ZP, draw.Over)

        new_hat := image.NewRGBA(pix.Bounds())
        draw.DrawMask(new_hat, new_hat.Bounds(), matched_hat, image.ZP, pix, image.ZP, draw.Over)
        pix = new_hat
      }

      // Fit the hat to the head
      draw_context.DrawImageAnchored(pix, (head.Min.X + head.Max.X)/2, head.Min.Y, 0.5, 0.5)


    }


  return draw_context.Image(), len(rects)
}


func (self HatDB) ApplyHats(img image.Image) (image.Image, int) {
  return self.ApplyHatsColorMatched(img, false, nil)
}
