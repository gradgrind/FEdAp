#include "fltk_json.h"
#include <FL/Fl_Double_Window.H>

using namespace std;

// +++ DoubleWindow
DoubleWindow::DoubleWindow(
    string_view name, int width, int height)
    : Widget(name, new Fl_Double_Window(width, height))
{}

const string_view DoubleWindow::widget_type()
{
    return string_view{wtype};
}

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

// --- DoubleWindow

// +++ Flex

Flex::Flex(
    string_view name, bool horizontal)
    : Widget(name, new Fl_Flex(horizontal ? Fl_Flex::ROW : Fl_Flex::COLUMN))
{}

const string_view Flex::widget_type()
{
    return string_view{wtype};
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

// --- Flex

//TODO: It should be possible to get the class names for the classes
// which are constructed ...
void AddFunctions()
{
    auto w1 = newDoubleWindow;
    add_function("Window:Double", w1);
    add_function("Flex", newFlex);
}
