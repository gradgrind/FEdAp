#ifndef MINION_H
#define MINION_H

#include <map>
#include <string>
#include <variant>
#include <vector>

namespace minion {

class MinionException : public std::exception
{
private:
    std::string message;

public:
    // Constructor accepts the exception message
    MinionException(
        std::string_view msg)
        : message{std::string{"MinionException: "} + std::string{msg}}
    {}

    // Override the what() method to return the error message
    const char *what() noexcept { return message.c_str(); }
};

using Char = unsigned char;

bool unicode_utf8(std::string &utf8, const std::string &unicode);

// *** The basic minion types ***
// Use forward declarations to allow mutual references.

class MinionMap;
class MinionList;
using MinionValue = std::variant<
    std::monostate, std::string, MinionMap, MinionList>;

void dump(std::string &valstr, MinionValue item, int level = -1);
std::string dump_list_items(const MinionList m, int level);
std::string dump_map_items(const MinionMap m, int level);

// The map class should preserve input order, so it is implemented as a vector.
// For very small maps this might be completely adequate, but if multiple
// lookups to larger maps are required, a proper map should be built.
struct MinionMapPair;

class MinionMap : public std::vector<MinionMapPair>
{
public:
    MinionValue get(std::string_view key);
    bool get_int(std::string_view key, int &value);
    bool get_string(std::string_view key, std::string &value);
};

class MinionList : public std::vector<MinionValue>
{};

struct MinionMapPair
{
    std::string key;
    MinionValue value;
};


class Minion
{
public:
    Minion(const std::string_view source);
    //void to_json(std::string &json_string, bool compact);

    MinionMap top_level;       // collect the top-level map here
    std::string error_message; // if not empty, explain failure

private:
    const std::string_view minion_string; // the source string
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
    MinionValue macro_replace(MinionValue item);
};

MinionMap read_minion(std::string_view minion_string);

} // END namespace minion

void testminion();

#endif // MINION_H
