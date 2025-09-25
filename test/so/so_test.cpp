#include <UnitTest++/UnitTest++.h>

#include "dolphin.hpp"

TEST(DepartureMessage) {
    Dolphin d;
    CHECK_EQUAL("So long, and thanks for all the fish", d.depart());
}
