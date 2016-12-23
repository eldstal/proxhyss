#include "detect.hpp"
#include <iostream>
#include <opencv2/core/core.hpp>

using std::cout;
using std::endl;


using namespace cv;


void TestFunc() {
  cout << "Yes, Go calls C++ code." << endl;
}

int FindFaces(uint8_t* pixels, int height, int width,
              struct Face* result, int maxresults) {

  Mat image(height, width, CV_8UC4);

  // TODO: Try with OpenCV
  // TODO: Try with dlib

  result[0].x1 = 12;
  result[0].x2 = 150;
  return maxresults;
}
