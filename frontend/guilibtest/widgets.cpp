#include "fltk_json.h"
#include <FL/Fl_Box.H>

void new_box(
    string_view name, string_view parent, json data)
{    
    Fl_Box* widg;
    string label;
    if (get_json_string(data, "LABEL", label)) {
        widg = new Fl_Box(0, 0, 0, 0, label.c_str());
    } else {
        widg = new Fl_Box(0, 0, 0, 0);
    }
    new WidgetData("Box", name, widg);
    // set boxtype
    //widget_set_box(name, data); NO, it assumes the field is present
    string s;
    if (get_json_string(data, "BOXTYPE", s)) {
        widg->box(magic_enum::enum_cast<Fl_Boxtype>(s).value());
    }
}
