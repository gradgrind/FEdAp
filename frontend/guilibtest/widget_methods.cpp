#include "widget_methods.h"
#include "magic_enum/magic_enum.hpp"

//TODO: move somewhere more appropriate ...
int get_json_int(
    json data, string_view key)
{
    try {
        return data[key];
    } catch (...) {
        throw fmt::format("Integer parameter '{}' expected", key);
    }
}

string get_json_string(
    json data, string_view key)
{
    try {
        return data[key];
    } catch (...) {
        throw fmt::format("String parameter '{}' expected", key);
    }
}

Fl_Color _get_color(
    json data)
{
    auto s = get_json_string(data, "COLOR");
    if (s.length() == 6) {
        return static_cast<Fl_Color>(stoi(s, nullptr, 16)) * 0x100;
    }
    throw fmt::format("Invalid color: '{}'", s);
}

Fl_Boxtype _get_boxtype(
    json data)
{
    auto s = get_json_string(data, "BOXTYPE");
    try {
        return magic_enum::enum_cast<Fl_Boxtype>(s).value();
    } catch (...) {
        throw fmt::format("Invalid box type: '{}'", s);
    }
}

// ---

void widget_size(
    string_view name, json data)
{
    auto w = get_widget(name);
    w->size(get_json_int(data, "WIDTH"), get_json_int(data, "HEIGHT"));
}

void widget_set_box(
    string_view name, json data)
{
    auto w = get_widget(name);
    w->box(_get_boxtype(data));
}

void widget_set_color(
    string_view name, json data)
{
    auto w = get_widget(name);
    w->color(_get_color(data));
}
