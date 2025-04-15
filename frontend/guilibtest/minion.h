#ifndef MINION_H
#define MINION_H

#include <json.hpp>
#include <string>
#include <variant>
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

// *** The basic minion types ***
// Use forward declarations to allow mutual references.

class MinionMap;
class MinionList;
using MinionValue = std::variant<
    std::monostate, std::string, MinionMap, MinionList>;

// The map class should preserve input order, so it is implemented as a vector.
// For very small maps this might be completely adequate, but if multiple
// lookups to larger maps are required, a proper map should be built.
struct MinionMapPair;

class MinionMap : public std::vector<MinionMapPair>
{
public:
    void add(std::string &key, MinionValue mval);
    MinionValue get(std::string & key);
};

class MinionList : public std::vector<MinionValue>
{
public:
    void add(MinionValue mval);
};

struct MinionMapPair
{
    const std::string key;
    const MinionValue value;
};


class Minion
{
public:
    Minion(const std::string &source);
    void to_json(std::string &json_string, bool compact);

    MinionMap top_level;            // collect the top-level map here
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
    std::map<std::string, MinionValue> macros;

    Char read_ch(bool instring);
    void unread_ch(Char ch);
    Char get_item(MinionValue &m);
    void get_list(MinionValue &m);
    bool get_map(MinionMap &m, Char terminator);

    void get_string(MinionValue &m);
    json macro_replace(MinionValue item);
};

} // END namespace minion

void testminion();
void readfile(std::string &data, const std::string &filepath);
void writefile(const std::string &data, const std::string &filepath);

#endif // MINION_H
