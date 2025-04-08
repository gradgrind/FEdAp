#ifndef MINION_H
#define MINION_H

#include <json.hpp>
#include <string>
using json = nlohmann::json;

namespace minion {
using Char = unsigned char;

bool unicode_utf8(std::string &utf8, const std::string &unicode);

class MinionParser
{
public:
    MinionParser(const std::string &source);
    void to_json(std::string &json_string, bool compact);

    json top_level;            // collect the top-level map here
    std::string error_message; // if not empty, explain failure

private:
    const std::string_view minion_string; // the source string
    const size_t source_size;
    int iter_i;
    int line_i;
    Char ch_pending;
    Char separator; // TODO: deprecated?

    Char read_ch(bool instring);
    void unread_ch(Char ch);
    Char get_item(json &j);
    void get_list(json &j);
    void get_map(json &j, Char terminator);
    void get_string(json &j);
    json macro_replace(json item);
};

// The basic minion types
enum value_type { M_STRING, M_MAP, M_LIST };
struct MinionValue
{
    value_type vtype;
    int index;
};

// The map class should preserve input order
struct MinionMapPair
{
    const std::string key;
    const MinionValue value;
};

struct MinionMap
{
    std::vector<MinionMapPair> data;
    std::map<const std::string *, const int> associate;
};

struct MinionList : public std::vector<MinionValue>
{};

class Minion
{
public:
    Minion(const std::string &source);
    void to_json(std::string &json_string, bool compact);

    json top_level;            // collect the top-level map here
    std::string error_message; // if not empty, explain failure

    MinionValue new_string(const std::string &s);
    MinionValue new_map();
    MinionValue new_list();

private:
    const std::string_view minion_string; // the source string
    std::vector<std::string> strings;
    std::vector<MinionMap> maps;
    std::vector<MinionList> lists;

    const size_t source_size;
    int iter_i;
    int line_i;
    Char ch_pending;
    Char separator; // TODO: deprecated?

    Char read_ch(bool instring);
    void unread_ch(Char ch);
    Char get_item(json &j);
    void get_list(json &j);
    void get_map(json &j, Char terminator);
    void get_string(json &j);
    json macro_replace(json item);
};

} // END namespace minion

void testminion();
void readfile(std::string &data, const std::string &filepath);
void writefile(const std::string &data, const std::string &filepath);

#endif // MINION_H
