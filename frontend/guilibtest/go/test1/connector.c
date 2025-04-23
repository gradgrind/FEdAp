#include "connector.h"
#include <stdio.h>

// _cgo_export.h is auto-generated and has Go //export funcs
#include "_cgo_export.h"

char* c_to_go_callback(char* data) {
  printf("C callback got '%s'\n", data);
  return GoCallback(data);
}


void init(char* data0) {
  printf("C says: init '%s'\n", data0);
  
  // Call actual start-up function
  //start(data0);
  
  char* result = c_to_go_callback("C callback");
  printf("C callback returned '%s'\n", result);
}
