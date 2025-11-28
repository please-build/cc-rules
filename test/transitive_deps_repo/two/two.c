#include "two.h"
#include "two_private.h"

#ifndef BUILDING_LIBTWO
#error "BUILDING_LIBTWO should be defined"
#endif

int two() {
    return TWO_INT;
}
