package main;

import (
    "github.com/elazarl/goproxy"
    "github.com/elazarl/goproxy/ext/image"
    "github.com/eldstal/proxhyss"
    "net/http"
    "image"
    "fmt"
    "log"
)

func main() {

  hats := proxhyss.InitHats("hats")
  
  proxy := goproxy.NewProxyHttpServer()
  proxy.OnResponse().Do(goproxy_image.HandleImage(func(img image.Image, ctx *goproxy.ProxyCtx) image.Image {

    //fmt.Printf("%v\n", ctx.Req.URL)
    nimg,_ := hats.ApplyHats(img);

    return nimg
  }))
  proxy.Verbose = false

  fmt.Printf("Set browser's SOCKS proxy to localhost:8080\n")
  log.Fatal(http.ListenAndServe("localhost:8080", proxy))
}
