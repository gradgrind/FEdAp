#include "fltk_minion.h"
#include <FL/Fl_Double_Window.H>
#include <FL/Fl_Flex.H>
#include <FL/Fl_Grid.H>
using namespace std;

// "Group" widgets

void new_window(
    string_view name, string_view parent, mmap data)
{
    int w = 800;
    get_minion_int(data, "WIDTH", w);
    int h = 600;
    get_minion_int(data, "HEIGHT", h);
    auto widg = new Fl_Double_Window(w, h);
    //TODO: widg->end(), or null current group, or ...?
    new WidgetData("Group:Window:Double", name, widg);
}

void new_flex(
    string_view name, string_view parent, mmap data)
{
    string orientation;
    get_minion_string(data, "ORIENTATION", orientation);
    auto widg = new Fl_Flex((orientation == "HORIZONTAL") ? Fl_Flex::ROW : Fl_Flex::COLUMN);
    //TODO: widg->end(), or null current group, or ...?
    new WidgetData("Group:Flex", name, widg);
}

void new_grid(
    string_view name, string_view parent, mmap data)
{
    auto widg = new Fl_Grid(0, 0, 0, 0);
    //TODO: widg->end(), or null current group, or ...?
    new WidgetData("Group:Grid", name, widg);
}

void _parm_widget_name(
    const mmap data, string &name)
{
    if (!get_minion_string(data, "NAME", name)) {
        throw fmt::format("Function '{}':\n no 'NAME' field", minion::dump_map_items(data, -1));
    }
}

void _parm_set_parent(
    const mmap data, Fl_Widget *widg)
{
    string parent;
    if (get_minion_string(data, "PARENT", parent) && !parent.empty()) {
        static_cast<Fl_Group *>(get_widget(parent))->add(widg);
    }
}

void _new_grid(
    mmap data)
{
    string name;
    _parm_widget_name(data, name);
    auto widg = new Fl_Grid(0, 0, 0, 0);
    Fl_Group::current(0);
    new WidgetData("Group:Grid", name, widg);
    _parm_set_parent(data, widg);
}
