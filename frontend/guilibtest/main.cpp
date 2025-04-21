#include "iofile.h"
#include "layout.h" //TODO: less ...
#include "widget_base.h"
#include "widgets.h"
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

        //tmp_run(guidata);

        auto dolist0 = guidata.get("GUI");
        if (holds_alternative<minion::MinionList>(dolist0)) {
            auto dolist = get<minion::MinionList>(dolist0);
            for (const auto &cmd : dolist) {
                do_GUI(get<minion::MinionMap>(cmd));
            }
            //auto tt = NEW_MyTable(minion::MinionMap{});
            //auto fl = widget_info.at("l_MainWindow");
            //static_cast<Fl_Group *>(fl.widget)->add(tt);
            //Fl::run();
        } else {
            cerr << "Input data not a GUI command list" << endl;
        }

        //Fl_Widget *win = new Fl_Double_Window(800, 500, "Table Flex");
        //Fl_Group::current(0);
        auto win = widget_info.at("MainWindow").widget;

        //Fl_Widget *l = new Fl_Flex(0, 0, win->w(), win->h());
        //Fl_Group::current(0);
        //static_cast<Fl_Group *>(win)->add(l);
        auto l = widget_info.at("l_MainWindow").widget;

        cout << "??? " << l->parent() << " ? " << win << endl;
        cout << "  l: " << l->w() << " ? " << l->h() << endl;
        cout << "  win: " << win->w() << " ? " << win->h() << endl;

        win->color(0xffffff00);
        static_cast<Fl_Flex *>(l)->margin(5);
        Fl_Widget *b = new Fl_Box(FL_BORDER_BOX, 0, 0, 0, 0, "box");
        b->color(0x00ff0000);
        static_cast<Fl_Group *>(l)->add(b);
        static_cast<Fl_Flex *>(l)->fixed(b, 100);
        Fl_Widget *table = NEW_MyTable(minion::MinionMap{});
        static_cast<Fl_Group *>(l)->add(table);
        static_cast<Fl_Group *>(win)->resizable(l);
        win->show();
        Fl::run();

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
