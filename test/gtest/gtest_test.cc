// Trivial test to illustrate gtest working.

#include "gtest/gtest.h"

#include "test/lib1.h"
#include "test/lib2.h"

namespace plz {

TEST(GTest, Number1) {
    EXPECT_EQ(107, get_number_1());
}

TEST(GTest, Number2) {
    EXPECT_EQ(215, get_number_2());
}

}
