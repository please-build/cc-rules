#include <stdio.h>

#include "one.h"

#ifdef BUILDING_LIBONE
#error "BUILDING_LIBONE should not be defined"
#endif

#ifdef BUILDING_LIBTWO
#error "BUILDING_LIBTWO should not be defined"
#endif

int main() {
    printf("one=%d one_plus_one=%d\n", one(), one_plus_one());
    return 0;
}
