#include "widget_base.h"
#include "functions.h"
#include "widgets.h"
#include <FL/Fl_Group.H>
#include <fmt/format.h>
using namespace std;
using namespace minion;

std::unordered_map<std::string, Widget> widget_info{};

void do_methods(
    Widget wd, MinionMap m)
{
    auto dolist = m.get("DO");
    if (holds_alternative<MinionList>(dolist)) {
        MinionList do_list = get<MinionList>(dolist);
        for (const auto &cmd : do_list) {
            MinionList m = get<MinionList>(cmd);
            string_view c = get<string>(m.at(0));
            wd.handle_method(wd, c, m);
        }
    } else if (dolist.index() != 0) {
        string s;
        minion::dump(s, dolist, 0);
        throw string{"Invalid DO list: "} + s;
    }
}

void do_NEW(
    string_view wtype, MinionMap m)
{
    //cout << "do_NEW " << wtype << ":" << minion::dump_map_items(m, 1) << endl;
    string name;
    string wwtype;
    w_method_handler h;
    Fl_Widget *w;
    if (m.get_string("NAME", name)) {
        if (wtype == "Window") {
            w = NEW_Window(m);
            h = w_group_method;
            wwtype = "Group:Window";
        } else if (wtype == "Vlayout") {
            w = NEW_Vlayout(m);
            h = w_flex_method;
            wwtype = "Group:Flex:V";
        } else if (wtype == "Hlayout") {
            w = NEW_Hlayout(m);
            h = w_flex_method;
            wwtype = "Group:Flex:H";
        } else if (wtype == "Grid") {
            w = NEW_Grid(m);
            h = w_grid_method;
            wwtype = "Group:Grid";
            // *** End of layouts, start of other widgets
        } else if (wtype == "Box") {
            w = NEW_Box(m);
            h = w_widget_method;
            wwtype = "Box";
        } else if (wtype == "Choice") {
            w = NEW_Choice(m);
            h = w_choice_method;
            wwtype = "Choice";
        } else if (wtype == "Output") {
            w = NEW_Output(m);
            h = w_input_method;
            wwtype = "Input:Output";
        } else if (wtype == "RowTable") {
            w = NEW_RowTable(m);
            h = w_rowtable_method;
            wwtype = "Table:Row";
        } else if (wtype == "MyTable") {
            //TODO-- ... just for testing
            w = NEW_MyTable(m);
            h = w_widget_method;
            wwtype = "TestTable";
        } else {
            throw fmt::format("Unknown widget type: {}", wtype);
        }
        string parent;
        if (m.get_string("PARENT", parent) && !parent.empty()) {
            auto pinfo = widget_info.at(parent);
            static_cast<Fl_Group *>(pinfo.widget)->add(w);
        }
        Widget wd{w, wwtype, h};
        widget_info.emplace(name, wd);
        // Handle methods
        do_methods(wd, m);
        return;
    }
    throw fmt::format("Bad NEW command: {}", minion::dump_map_items(m, -1));
}

void do_GUI(
    MinionMap m)
{
    string ws;
    if (m.get_string("NEW", ws)) {
        do_NEW(ws, m);
    } else if (m.get_string("WIDGET", ws)) {
        // Handle methods
        auto wd = widget_info.at(ws);
        do_methods(wd, m);
    } else if (m.get_string("FUNCTION", ws)) {
        auto f = function_map.at(ws);
        f(m);
    } else {
        throw fmt::format("Invalid GUI parameters: {}", dump_map_items(m, -1));
    }
}
