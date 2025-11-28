#include "one.h"
#include "two.h"
#include "one_private.h"

#ifndef BUILDING_LIBONE
#error "BUILDING_LIBONE should be defined"
#endif

#ifdef BUILDING_LIBTWO
#error "BUILDING_LIBTWO should not be defined"
#endif

int one() {
    return ONE_INT;
}

int one_plus_one() {
    return two();
}
