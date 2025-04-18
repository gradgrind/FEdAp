#include "iofile.h"
#include <fstream>

using namespace std;

bool readfile(
    string &data, const string &filepath)
{
    std::ifstream file(filepath);
    if (file) {
        data.assign((istreambuf_iterator<char>(file)), istreambuf_iterator<char>());
        return true;
    }
    return false;
}

bool writefile(
    const string &data, const string &filepath)
{
    std::ofstream file(filepath);
    if (file) {
        file << data;
        file.close();
        return true;
    }
    return false;
}
