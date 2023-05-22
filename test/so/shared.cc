// Simple shared object to test that cc_shared_object actually does something useful.

#include <string>

#include "test/embed/embedded_files.h"

std::string get_file1() {
  return plz::embedded_file1_contents();
}

std::string get_file3() {
  return plz::embedded_file3_contents();
}
