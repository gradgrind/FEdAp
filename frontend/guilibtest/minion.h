#ifndef MINION_H
#define MINION_H

#include <json.hpp>
#include <string>
using json = nlohmann::json;

namespace Minion {
using Char = unsigned char;

class MinionParser
{
public:
    MinionParser(std::string_view source);
    void to_json(std::string &json_string, bool compact);

    json top_level;            // collect the top-level map here
    std::string error_message; // if not empty, explain failure

private:
    const std::string_view minion_string; // the source string
    const size_t source_size;
    int iter_i;
    int line_i;
    Char ch_pending;
    Char separator;

    Char read_ch(bool instring);
    void unread_ch(Char ch);
    json get_item();
    json get_list();
    json get_map(Char terminator);
    json get_string();
    json macro_replace(json item);
};

} // END namespace Minion

void testminion();
void readfile(std::string &data, const std::string &filepath);
void writefile(const std::string &data, const std::string &filepath);

#endif // MINION_H
