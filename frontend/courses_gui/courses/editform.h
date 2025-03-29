#ifndef EDITFORM_H
#define EDITFORM_H

#include <FL/Fl_Box.H>
#include <FL/Fl_Grid.H>
#include <FL/Fl_Output.H>
#include <FL/Fl_Select_Browser.H>
#include <string>
#include <vector>

struct EditFormEntry
{
    Fl_Widget *widget;
    std::string name;
    int padabove;
    bool spanning;
    bool growable;
};

class EditForm : public Fl_Grid
{
public:
    std::vector<EditFormEntry> entries;

    EditForm();
    void add_value(const char *name, const char *label);
    void add_separator();
    void add_list(const char *name, const char *label);
    void do_layout();
};

#endif // EDITFORM_H
