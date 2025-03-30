#include "fltk_json.h"

Widget::Widget(
    std::string_view _name, Fl_Widget *_widget)
    : name{_name}
    , widget{_widget}
{
    //TODO: Would unnamed widgets be useful? I could skip the
    // use of the map (also on delete).

    if (WidgetMap.contains(name)) {
        throw fmt::format("Widget name already exists: %s", name);
    }
    //WidgetMap.insert({name, this}); ... or ...
    WidgetMap.emplace(name, this);
}

Widget::~Widget()
{
    WidgetMap.erase(name);
}

Widget *Widget::get(
    std::string_view name)
{
    return WidgetMap.at(name);
}

Fl_Widget *Widget::get_flwidget(
    std::string_view name)
{
    return WidgetMap.at(name)->widget;
}
