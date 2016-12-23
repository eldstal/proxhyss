#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

  struct Face {
    int x1;
    int y1;
    int x2;
    int y2;
  };

  void TestFunc();

  // Image is 1 byte per channel, RGBA, row major.
  // Size is width*height*4 bytes.
  int FindFaces(uint8_t* pixels, int height, int width,
                struct Face* result, int maxresults);

#ifdef __cplusplus
}
#endif
