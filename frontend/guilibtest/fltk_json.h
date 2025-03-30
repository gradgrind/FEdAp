#ifndef FLTK_JSON_H
#define FLTK_JSON_H

#include <FL/Fl_Flex.H>
#include <fmt/format.h>
#include <functional>
#include <iostream>
#include <string>
#include <unordered_map>

#include <json.hpp>
using json = nlohmann::json;

using namespace std;

class Gui
{
public:
    inline static unordered_map<string_view, function<void(string_view, json)>>
        FunctionMap;
};

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

//TODO: Change the constructors to take "standard" arguments, so that
// no extra constructor-functions are necessary? With parent as
// widget argument? Or "" ... using the automatic add-to-group, or
// blocking that feature (set current group to 0). Other parameters
// as JSON?
class DoubleWindow : public Widget
{
public:
    DoubleWindow(string_view name, int width = 800, int height = 600);
    const string_view widget_type() { return string_view{"Window:Double"}; }
};

class Flex : public Widget
{
public:
    Flex(string_view name, bool horizontal = false);
    const string_view widget_type() { return string_view{"Flex"}; }
};

// In this case, the widget (name) is the new name
void newFlex(string_view name, json data);

#endif // FLTK_JSON_H
