#include "iofile.h"
#include "layout.h" //TODO: less ...
#include <FL/Fl_Box.H>
#include <FL/Fl_Double_Window.H>
#include <FL/Fl_Flex.H>
#include <FL/fl_ask.H>
#include <FL/fl_draw.H>
#include <chrono>
#include <fmt/format.h>
#include <iostream>
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

//TODO: Window might get a special callback ...
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

    /* testing stuff
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
    */

    minion::MinionMap guidata;

    try {
        string gui;
        string fp{"data/gui.minion"};
        if (readfile(gui, fp)) {
            cout << "Reading " << fp << endl;
            try {
                guidata = minion::read_minion(gui);
            } catch (minion::MinionException &e) {
                cerr << e.what() << endl;
                return 99;
            }
        } else {
            cerr << "Error opening file: " << fp << endl;
            return 99;
        }
        //cout << minion::dump_map_items(guidata, 0) << endl;

        tmp_run(guidata);

        return 0;

        /*
        win->resizable(f1a);

        win->callback(main_callback);
        win->end();
        win->show();
        return Fl::run();
        */
    } catch (const std::exception &ex) {
        cout << "EXCEPTION: " << ex.what() << endl;
    } catch (const std::string &ex) {
        cout << "ERROR: " << ex << endl;
    }
}
