#ifndef COURSES_GUI_H
#define COURSES_GUI_H

#include <FL/Fl_Box.H>
#include <FL/Fl_Flex.H>
#include <FL/Fl_Grid.H>
#include <map>
#include <string>

class CoursesGui : public Fl_Grid
{
public:
    CoursesGui();

    std::map<std::string, Fl_Widget *> Widgets;
};

#endif // COURSES_GUI_H
