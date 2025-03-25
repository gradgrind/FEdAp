#include "courses/courses_gui.h"
#include <FL/Fl_Double_Window.H>
#include <iostream>

using namespace std;

int main()
{
    cout << "Hello World!" << endl;
    int w0 = 1000;
    int h0 = 700;
    auto win = new Fl_Double_Window(w0, h0);
    win->color(FL_WHITE);
    auto vbox = new Fl_Flex(0, 0, w0, h0);
    auto panel = new Fl_Box(FL_BORDER_FRAME, 0, 0, 0, 0, "Panel");
    vbox->fixed(panel, 100);
    auto cg = new CoursesGui();

    vbox->end();
    win->resizable(vbox);
    win->end();
    win->show();
    return Fl::run();
}
