#ifndef IOFILE_H
#define IOFILE_H

#include <string>

bool readfile(std::string &data, const std::string &filepath);
bool writefile(const std::string &data, const std::string &filepath);

#endif // IOFILE_H
