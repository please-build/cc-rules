#include <dlfcn.h>
#include <string>
#include <UnitTest++/UnitTest++.h>

#include "test/embed/embedded_files.h"

TEST(File1Matches) {
  void* shared = dlopen("test/so/shared.so", RTLD_LAZY | RTLD_LOCAL);
  CHECK(shared);
  auto get_file1 = (std::string(*)())dlsym(shared, "get_file1");
  CHECK(get_file1);
  CHECK_EQUAL(embedded_file1_contents(), get_file1());
}
