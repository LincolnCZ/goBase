#ifndef __S2S_C__
#define __S2S_C__

#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

struct Buffer
{
    void *buffer;
    int size;
};

extern int initialize(const char *myName, const char *s2sKey, int myType);
// subscribe: int subscribe(PSubFilter filters);
extern int subscribe(struct Buffer input);

// pollNotify: PNotifyResult pollNotify();
extern struct Buffer pollNotify();

extern int setMine(const char *binData, int size);
extern int delMine();

// getMine: S2sMeta getMine();
extern struct Buffer getMine();

#ifdef __cplusplus
}
#endif

#endif
