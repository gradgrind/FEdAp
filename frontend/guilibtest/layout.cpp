#include "fltk_json.h"
#include <FL/Fl_Double_Window.H>
#include <FL/Fl_Flex.H>
#include <FL/Fl_Grid.H>
using namespace std;

// "Group" widgets

void new_window(
    string_view name, string_view parent, json data)
{
    int w = 800;
    if (data.contains("WIDTH")) {
        w = data["WIDTH"];
    }
    int h = 600;
    if (data.contains("HEIGHT")) {
        h = data["HEIGHT"];
    }
    auto widg = new Fl_Double_Window(w, h);
    //TODO: widg->end(), or null current group, or ...?
    new WidgetData("Group:Window:Double", name, widg);
}

void new_flex(
    string_view name, string_view parent, json data)
{
    auto hz = data.contains("HORIZONTAL") && data["HORIZONTAL"];
    auto widg = new Fl_Flex(hz ? Fl_Flex::ROW : Fl_Flex::COLUMN);
    //TODO: widg->end(), or null current group, or ...?
    new WidgetData("Group:Flex", name, widg);
}

void new_grid(
    string_view name, string_view parent, json data)
{
    auto widg = new Fl_Grid(0, 0, 0, 0);
    //TODO: widg->end(), or null current group, or ...?
    new WidgetData("Group:Grid", name, widg);
}

void _parm_widget_name(
    const json data, string &name)
{
    if (!get_json_string(data, "NAME", name)) {
        throw fmt::format("Function '{}':\n no 'NAME' field", data);
    }
}

void _parm_set_parent(
    const json data, Fl_Widget *widg)
{
    string parent;
    if (get_json_string(data, "PARENT", parent) && !parent.empty()) {
        static_cast<Fl_Group *>(get_widget(parent))->add(widg);
    }
}

void _new_grid(
    json data)
{
    string name;
    _parm_widget_name(data, name);
    auto widg = new Fl_Grid(0, 0, 0, 0);
    Fl_Group::current(0);
    new WidgetData("Group:Grid", name, widg);
    _parm_set_parent(data, widg);
}
