#include "editform.h"
#include <FL/fl_draw.H>
#include <iostream>
#include <ostream>

EditForm::EditForm()
    : Fl_Flex(Fl_Flex::COLUMN)
{
    box(FL_BORDER_FRAME);
    gap(10);
    end();
}

Fl_Output* EditForm::add_value(
    const char* label)
{
    auto e1 = new Fl_Output(0, 0, 0, 0, label);
    e1->align(FL_ALIGN_LEFT);
    e1->color(fl_rgb_color(255, 255, 200));

    //TODO
    e1->callback([](Fl_Widget* w, void* a) {
        std::cout << "Activated: " << Fl::callback_reason() << std::endl;
    });

    //TODO--
    e1->value(label);
    // Setting the tooltip allows overlong texts to be readable
    e1->tooltip(label);

    e1->clear_visible_focus(); // no cursor, but text cannot be copied
    e1->horizontal_label_margin(5);
    add(e1);
    fixed(e1, 30);
    return e1;
}

void EditForm::add_separator()
{
    auto sep = new EditFormSeparator(&label_width);
    add(sep);
    fixed(sep, 10);
}

void EditForm::do_layout()
{
    int labwidth = 0;
    int wl, hl;
    for (int i = 0; i < children(); ++i) {
        auto cw = child(i);
        wl = 0;
        cw->measure_label(wl, hl);
        if (wl > labwidth)
            labwidth = wl;
    }
    label_width = labwidth;
    margin(labwidth + 15, 5, 5, 5);
}

EditFormSeparator::EditFormSeparator(
    int* label_width)
    : Fl_Box(0, 0, 0, 0)
    , p_label_width{label_width}
{}

void EditFormSeparator::draw()
{
    //    Fl_Box::draw();

    int bw = w();
    int bh = h();
    int bm = y() + bh / 2;
    int start = x() - *p_label_width + 10;
    int end = start + *p_label_width + bw - 20;
    //std::cout << bw << " " << bh << " " << x() << " " << y() << std::endl;
    //std::cout << " + " << *p_label_width << std::endl;
    fl_line(start, bm, end, bm);
}
