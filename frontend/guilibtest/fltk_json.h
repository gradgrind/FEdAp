#ifndef FLTK_JSON_H
#define FLTK_JSON_H

#include <FL/Fl_Flex.H>
#include <fmt/format.h>
#include <iostream>
#include <string>
#include <unordered_map>
using namespace std;

class Widget
{
    inline static unordered_map<string_view, Widget *> WidgetMap;

    std::string name;
    Fl_Widget *widget;

protected:
    Widget(
        std::string_view _name, Fl_Widget *_widget)
        : name{_name}
        , widget{_widget}
    {
        //TODO: Would unnamed widgets be useful? I could skip the
        // use of the map (also on delete).

        if (WidgetMap.contains(name)) {
            throw fmt::format("Widget name already exists: %s", name);
        }
        //WidgetMap.insert({name, this});
        WidgetMap.emplace(name, this);
        cout << "Widget " << this << endl;
    }

    ~Widget() { WidgetMap.erase(name); }

public:
    virtual const string_view widget_type() = 0;

    static Widget *get(
        std::string_view name)
    {
        return WidgetMap.at(name);
    }

    static Fl_Widget *get_flwidget(
        std::string_view name)
    {
        return WidgetMap.at(name)->widget;
    }
};

class Flex : public Widget
{
public:
    Flex(string name, bool horizontal = false);
    const string_view widget_type() { return string_view{"Flex"}; }
};

#endif // FLTK_JSON_H
