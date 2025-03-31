#include "fltk_json.h"
#include <FL/Fl_Double_Window.H>
#include <iostream>

using namespace std;

void newWindow(
    string_view name, json data)
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
    new WidgetData("Window:Double", name, widg);
}

void newFlex(
    string_view name, json data)
{
    auto hz = data.contains("HORIZONTAL") && data["HORIZONTAL"];
    auto widg = new Fl_Flex(hz ? Fl_Flex::ROW : Fl_Flex::COLUMN);
    //TODO: widg->end(), or null current group, or ...?
    new WidgetData("Flex", name, widg);
}
