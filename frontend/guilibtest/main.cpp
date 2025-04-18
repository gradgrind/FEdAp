#include "widget_methods.h"
#include <FL/Fl_Box.H>
#include <FL/Fl_Double_Window.H>
#include <FL/Fl_Flex.H>
#include <FL/fl_ask.H>
#include <FL/fl_draw.H>

#include <fmt/format.h>

#include <iostream>

#include "minion.h"

#include <chrono>

using namespace std;

chrono::time_point<std::chrono::steady_clock> t0;

void timetest_cb(
    void *a)
{
    const auto end = std::chrono::steady_clock::now();
    //cout << "Tick " << std::chrono::duration_cast<std::chrono::microseconds>(end - t0) << " ≈ "
    //     << (end - t0) / 1ms << "ms" << endl;
    t0 = end;
    //Fl::add_timeout(0.1, timetest_cb); // retrigger timeout
    // or use repeat_timeout for more regular intervals
    // (probably hardly any difference)
    Fl::repeat_timeout(2.0, timetest_cb);
}

void timetest()
{
    using namespace std::literals;
    t0 = std::chrono::steady_clock::now();
    Fl::add_timeout(1.0, timetest_cb);
}

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
    //timetest();
    //testminion();
    //return 0;

    string _bt1{"FL_NO_BOX"};
    cout << _bt1 << ": " << magic_enum::enum_cast<Fl_Boxtype>(_bt1).value() << endl;

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

    // *** Here is the real start ...
    try {
        int w0 = 1000;
        int h0 = 700;
        gui_new("W0", "Window", "", read_minion(fmt::format("WIDTH:{} HEIGHT:{}", w0, h0)));
        auto win = static_cast<Fl_Double_Window *>(get_widget("W0"));
        win->color(FL_WHITE);

        //* Test creation and deletion
        gui_new("F0", "Flex", "", {});
        auto f0a = get_widget("F0");
        delete f0a;
        //*/

        gui_new("F1", "Flex", "", {});
        cout << "??? " << minion::dump_list_items(list_widgets(), -1) << endl;

        auto f1a = static_cast<Fl_Flex *>(get_widget("F1"));
        auto ud = static_cast<WidgetData *>(f1a->user_data());
        cout << "? " << ud->widget_name() << " @ " << ud->widget_type() << " ~ "
             << ud->widget_type_name() << endl;
        //f1a->size(w0, h0);
        widget_set_size("F1", read_minion(fmt::format("WIDTH:{} HEIGHT:{}", w0, h0)));
        //f1a->box(FL_BORDER_BOX);
        widget_set_box("F1", read_minion("BOXTYPE:FL_BORDER_BOX"));
        //f1a->color(0x00ff0000);
        widget_set_color("F1", read_minion("COLOR:00FF00"));

        gui_new("B1", "Box", "", read_minion("LABEL:B1 BOXTYPE:FL_ENGRAVED_BOX"));
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
    } catch (const std::exception &ex) {
        cout << "EXCEPTION: " << ex.what() << endl;
    } catch (const std::string &ex) {
        cout << "ERROR: " << ex << endl;
    }
}
