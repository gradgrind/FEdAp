#include "widget_methods.h"
using namespace std;

//TODO: move somewhere more appropriate ...
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

Fl_Color _get_color(
    mmap data)
{
    auto s = get_minion_string(data, "COLOR");
    if (s.length() == 6) {
        return static_cast<Fl_Color>(stoi(s, nullptr, 16)) * 0x100;
    }
    throw fmt::format("Invalid color: '{}'", s);
}

Fl_Boxtype _get_boxtype(
    mmap data)
{
    string s;
    if (get_minion_string(data, "BOXTYPE", s)) {
        return magic_enum::enum_cast<Fl_Boxtype>(s).value();
    } else {
        throw fmt::format("Invalid box type: '{}'", s);
    }
}

// ---

//OLD--
void widget_set_size(
    string_view name, mmap data)
{
    auto w = get_widget(name);
    w->size(get_minion_int(data, "WIDTH"), get_minion_int(data, "HEIGHT"));
}

void widget_set_box(
    string_view name, mmap data)
{
    auto w = get_widget(name);
    w->box(_get_boxtype(data));
}

void widget_set_color(
    string_view name, mmap data)
{
    auto w = get_widget(name);
    w->color(_get_color(data));
}

//NEW++
void _widget_set_size(
    Fl_Widget* w, mmap data)
{
    w->size(get_minion_int(data, "WIDTH"), get_minion_int(data, "HEIGHT"));
}

void _widget_set_box(
    Fl_Widget* w, mmap data)
{
    w->box(_get_boxtype(data));
}

void _widget_set_color(
    Fl_Widget* w, mmap data)
{
    w->color(_get_color(data));
}

member_map widget_methods{{"set_size", _widget_set_size},
                          {"set_box", _widget_set_box},
                          {"set_color", _widget_set_color}};
// Each subclass has its own map, constructed by COPYING the parent class
// map and adding its own methods
