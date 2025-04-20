#ifndef FUNCTIONS_H
#define FUNCTIONS_H

#include "minion.h"
#include <functional>

using function_handler = std::function<void(minion::MinionMap)>;
extern std::unordered_map<std::string, function_handler> function_map;

#endif // FUNCTIONS_H
