#ifndef EDITFORM_H
#define EDITFORM_H

#include <FL/Fl_Box.H>
#include <FL/Fl_Flex.H>
#include <FL/Fl_Output.H>

class EditForm : public Fl_Flex
{
    int label_width;

public:
    EditForm();
    Fl_Output *add_value(const char *label);
    void add_separator();
    void do_layout();
};

class EditFormSeparator : public Fl_Box
{
    void draw() FL_OVERRIDE;
    int *p_label_width;

public:
    EditFormSeparator(int *label_width);
};

#endif // EDITFORM_H
