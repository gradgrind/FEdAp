#include "fltk_json.h"
#include <FL/Fl_Box.H>
#include <FL/Fl_Double_Window.H>
#include <FL/fl_ask.H>
#include <FL/fl_draw.H>

#include <fmt/format.h>

#include <iostream>

using namespace std;

void main_callback(
    Fl_Widget *, void *)
{
    if (Fl::event() == FL_SHORTCUT && Fl::event_key() == FL_Escape)
        return; // ignore Escape

    //TODO: If changed data, ask about closing
    if (fl_choice("Are you sure you want to quit?", "continue", "quit", NULL))
        exit(0);
}

class A
{
public:
    A() { cout << "A's constructor called" << endl; }
    ~A() { cout << "A's destructor called" << endl; }
};

class B
{
public:
    B() { cout << "B's constructor called" << endl; }
    ~B() { cout << "B's destructor called" << endl; }
};

class C : public B, public A // Note the order
{
public:
    C() { cout << "C's constructor called" << endl; }
    ~C() { cout << "C's destructor called" << endl; }
};

int main()
{
    //FL_NORMAL_SIZE = 16;
    //Fl::set_font(FL_HELVETICA, "DejaVu Serif");
    auto s = fmt::format("? {:>10} {}", "This", 3);
    cout << s << endl;

    try {
        throw s;
    } catch (std::string &e) {
        std::cout << "CRITICAL: " << e << endl;
    }

    auto c = new C;
    delete c;

    int w0 = 1000;
    int h0 = 700;
    auto win = new Fl_Double_Window(0, 0);
    win->size(w0, h0);
    win->color(FL_WHITE);

    // Test creation and deletion
    newFlex("F0", json{});
    auto f0a = widget_map.at("F0");
    delete f0a;
    //--

    newFlex("F1", json{});
    for (auto const &[k, v] : widget_map) {
        cout << "??? " << k << " @ " << v << endl;
    }

    auto f1a = (Fl_Flex *) widget_map.at("F1");
    auto ud = (WidgetData *) f1a->user_data();
    cout << "? " << ud->widget_name() << " @ " << ud->widget_type() << " ~ "
         << ud->widget_type_name() << endl;
    f1a->size(w0, h0);
    f1a->box(FL_BORDER_BOX);
    f1a->color(FL_GREEN);
    auto vbox1 = new Fl_Box(FL_BORDER_BOX, 0, 0, 0, 0, "B1");
    vbox1->color(FL_RED);
    f1a->fixed(vbox1, 200);
    auto vbox2 = new Fl_Box(FL_BORDER_BOX, 0, 0, 0, 0, "B2");
    vbox2->color(FL_YELLOW);
    f1a->fixed(vbox2, 200);
    f1a->end();

    win->resizable(f1a);

    win->callback(main_callback);
    win->end();
    win->show();
    return Fl::run();
}
