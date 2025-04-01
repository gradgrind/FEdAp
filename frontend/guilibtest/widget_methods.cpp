#include "widget_methods.h"

//TODO: move somewhere more appropriate ...
int get_json_int(
    json data, string_view key)
{
    try {
        int i = data[key];
        return i;
    } catch (...) {
        throw fmt::format("Integer parameter '{}' expected", key);
    }
}

Fl_Color get_RGB()
{
    return static_cast<Fl_Color>(stoi("FFffaa", nullptr, 16) * 256);
}

// ---

void widget_size(
    string_view name, json data)
{
    auto w = get_widget(name);
    w->size(get_json_int(data, "WIDTH"), get_json_int(data, "HEIGHT"));
}
