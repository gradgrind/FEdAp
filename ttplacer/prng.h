#ifndef PRNG_H
#define PRNG_H

#include <cstdint>
#include <random>

class Random
{
    struct prng_32_short_s
    {
        uint32_t a;
        uint32_t b;
        uint32_t increment;
    };

    prng_32_short_s s;

public:
    Random(uint32_t seed = 0)
    {
        if (seed == 0) { // Use a non-deterministic seed
            std::random_device rd;
            seed = rd();
        }
        s = {.a = seed, .b = 0, .increment = 0};
    }

    uint32_t random()
    {
        s.a = ((s.a << 14) | (s.a >> 18)) ^ s.b;
        s.increment += 1111111111;
        s.b = ((s.b << 21) | (s.b >> 11)) + s.increment;
        return s.a + 1111111111;
    }
};

#endif // PRNG_H
