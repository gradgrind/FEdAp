#include "functions.h"
#include <FL/Fl.H>
#include <iostream>
using mmap = minion::MinionMap;

//TODO: This is just temporary – I still need to receive input somehow!
void f_RUN(
    mmap m)
{
    auto cc = Fl::run();
    std::cout << "Main Loop ended: " << cc << std::endl;
}

std::unordered_map<std::string, function_handler> function_map{{"RUN", f_RUN}};
