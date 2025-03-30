#ifndef FLTK_JSON_H
#define FLTK_JSON_H

#include <FL/Fl_Flex.H>
#include <fmt/format.h>
#include <functional>
#include <string>
#include <unordered_map>

#include <json.hpp>
using json = nlohmann::json;

using namespace std;

// The reason I was using classes with multiple inheritance was
// primarily so that deletion of a widget (especially as a result
// of deleting its group) could allow removal of the Widget entry
// in the widget table, and releasing of its space.
// Without this there can be dangling pointers.
// Another advantage is that the fltk widget has effectively access
// to the full Widget data.

// Could something like this work?
class _Flex : public Fl_Flex
{
    inline static const std::string wtype{"Flex"};

    // WidgitData wdata;

public:
    _Flex(
        string_view name, bool horizontal)
        : Fl_Flex(horizontal ? Fl_Flex::ROW : Fl_Flex::COLUMN)
    {}

    ~_Flex() {}
};

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
    Widget(std::string_view _name, Fl_Widget *_widget);
    ~Widget();

public:
    virtual const string_view widget_type() = 0;

    static Widget *get(std::string_view name);
    static Fl_Widget *get_flwidget(std::string_view name);
};

//TODO???: Change the constructors to take "standard" arguments, so that
// no extra constructor-functions are necessary? With parent as
// widget argument? Or "" ... using the automatic add-to-group, or
// blocking that feature (set current group to 0). Other parameters
// as JSON?

class DoubleWindow : public Widget
{
    inline static const std::string wtype{"Window:Double"};

public:
    DoubleWindow(string_view name, int width = 800, int height = 600);
    const string_view widget_type();
};

void newDoubleWindow(string_view name, json data);

class Flex : public Widget
{
    inline static const std::string wtype{"Flex"};

public:
    Flex(string_view name, bool horizontal = false);
    const string_view widget_type();
};

// In this case, the widget (name) is the new name
void newFlex(string_view name, json data);

#endif // FLTK_JSON_H
