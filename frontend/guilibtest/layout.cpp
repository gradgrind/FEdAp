#include "fltk_json.h"
#include <FL/Fl_Double_Window.H>

using namespace std;

DoubleWindow::DoubleWindow(
    string_view name, int width, int height)
    : Widget(name, new Fl_Double_Window(width, height))
{}

void newDoubleWindow(
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
    new DoubleWindow(name, w, h);
}

Flex::Flex(
    string_view name, bool horizontal)
    : Widget(name, new Fl_Flex(horizontal ? Fl_Flex::ROW : Fl_Flex::COLUMN))
{
    cout << "Flex " << this << endl;
}

void newFlex(
    string_view name, json data)
{
    auto hz = data.contains("HORIZONTAL") && data["HORIZONTAL"];
    new Flex(name, hz);
}

void add_function(
    string_view name, function<void(string_view, json)> f)
{
    if (Gui::FunctionMap.contains(name)) {
        throw fmt::format("Can't define function %s, it is already defined",
                          name);
    }
    Gui::FunctionMap.emplace(name, f);
}

//TODO: It should be possible to get the class names for the classes
// which are constructed ...
void AddFunctions()
{
    auto w1 = newDoubleWindow;
    add_function("Window:Double", w1);
    add_function("Flex", newFlex);
}
