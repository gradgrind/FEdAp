#include "iofile.h"
#include "minion.h"
#include <cstdio>
#include <cstdlib>
#include <time.h>

using namespace minion;

int main()
{
    auto fplist = {
        //
        "../../testdata/Demo1_tt.json",
        "../../testdata/x01_tt.json",
        //
    };

    std::string indata;

    struct timespec start, end; //, xtra;
    //struct timespec remaining, request = {1, 0}; // 1 sec.

    MValue m;

    //for (int i = 0; i < 10; ++i) {
    for (const auto& fp : fplist) {
        indata = readfile(fp);
        if (indata.empty()) {
            printf("File not found: %s\n", fp);
            exit(1);
        }

        //printf("Taking a nap...\n");
        //fflush(stdout);
        //nanosleep(&request, &remaining);

        clock_gettime(CLOCK_PROCESS_CPUTIME_ID, &start); // Initial timestamp

        m = Reader::read(indata);

        clock_gettime(CLOCK_PROCESS_CPUTIME_ID, &end); // Get current time

        //printf("Taking another nap...\n");
        //fflush(stdout);
        //nanosleep(&request, &remaining);

        //Writer w(m, -1);
        m = {}; // free memory

        //clock_gettime(CLOCK_PROCESS_CPUTIME_ID, &xtra); // Get current time

        double elapsed = end.tv_sec - start.tv_sec;
        elapsed += (end.tv_nsec - start.tv_nsec) / 1000000.0;
        printf("%0.2f milliseconds elapsed\n", elapsed);

        //elapsed = xtra.tv_sec - end.tv_sec;
        //elapsed += (xtra.tv_nsec - end.tv_nsec) / 1000000.0;
        //printf("%0.2f milliseconds dumping\n", elapsed);
    }
    printf("  - - - - -\n");
    //}

    /*
    m = Reader::read(indata);
    if (m.type() == T_Error) {
        printf("PARSE ERROR: %s\n", m.error_message());
    } else if (!m.is_null()) {
        Writer writer(m, 0);
        const char* result = writer.dump_c();
        if (result)
            printf("\n -->\n%s\n", result);
        else
            printf("*** Dump failed\n");
    }
    */

    return 0;
}
