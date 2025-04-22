#include "widget_methods.h"
#include <FL/Fl_Box.H>
#include <FL/Fl_Flex.H>
#define MAGIC_ENUM_RANGE_MIN 0
#define MAGIC_ENUM_RANGE_MAX 255
#include "magic_enum/magic_enum.hpp"
#include <fmt/format.h>
using namespace std;

Fl_Color get_colour(
    string& colour)
{
    if (colour.length() == 6) {
        return static_cast<Fl_Color>(stoi(colour, nullptr, 16)) * 0x100;
    }
    throw fmt::format("Invalid colour: '{}'", colour);
}

Fl_Boxtype get_boxtype(
    string& boxtype)
{
    return magic_enum::enum_cast<Fl_Boxtype>(boxtype).value();
}

int int_param(
    minion::MinionList m, int i)
{
    return stoi(get<string>(m.at(i)));
}

Fl_Color colour_param(
    minion::MinionList m, int i)
{
    return get_colour(get<string>(m.at(i)));
}

void left_label(
    Fl_Widget *w, minion::MinionList m)
{
    // In an Hlayout an empty box will be inserted as the the widget
    // itself allocates no space for the label.
    auto lbl = get<string>(m.at(1));
    w->copy_label(lbl.c_str());
    auto pw = dynamic_cast<Fl_Flex *>(w->parent());
    if (pw && pw->horizontal()) {
        auto padbox = new Fl_Box(0, 0, 0, 0);
        pw->insert(*padbox, w);
        int wl, hl;
        w->measure_label(wl, hl);
        pw->fixed(padbox, wl + int_param(m, 2));
    }
}

/*TODO: move somewhere more appropriate ... DEPRECATED? see minion.h
int get_minion_int(
    mmap data, string_view key)
{
    auto s = data.get(key);
    if (holds_alternative<string>(s)) {
        // Read as integer
        return stoi(get<string>(s));
    }
    throw fmt::format("Integer parameter '{}' expected", key);
}

bool get_minion_int(
    mmap data, string_view key, int& value)
{
    auto s = data.get(key);
    if (holds_alternative<string>(s)) {
        // Read as integer
        value = stoi(get<string>(s));
        return true;
    }
    return false;
}

string get_minion_string(
    mmap data, string_view key)
{
    auto s = data.get(key);
    if (holds_alternative<string>(s)) {
        return get<string>(s);
    }
    throw fmt::format("String parameter '{}' expected", key);
}

// This version sets the referenced string argument if the key exists.
// Return true if successful.
bool get_minion_string(
    mmap data, string_view key, string& value)
{
    auto s = data.get(key);
    if (holds_alternative<string>(s)) {
        value = get<string>(s);
        return true;
    }
    return true;
}
*/
