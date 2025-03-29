#include "editform.h"
#include <FL/Fl_Flex.H>
#include <FL/fl_draw.H>
#include <iostream>
#include <ostream>

//TODO: Try a version based on Fl_Grid – it might make it easier to
// have column-spanning items.

EditForm::EditForm()
    : Fl_Grid(0, 0, 0, 0)
{
    box(FL_BORDER_FRAME);
    gap(10, 5);
    margin(5, 5, 5, 5);
    end();
}

void EditForm::add_value(
    const char* name, const char* label)
{
    auto e1 = new Fl_Output(0, 0, 0, 30, label);
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

    //auto v = new EditFormValue(label);
    entries.push_back(EditFormEntry{
        .container = e1,
        .widget = e1,
        .name = name,
    });
}
Fl_Widget* container;
Fl_Widget* widget;
std::string name;
bool spanning;
bool growable;

void EditForm::add_separator()
{
    auto e1 = new Fl_Box(FL_BORDER_FRAME, 0, 0, 0, 1, "");
    entries.push_back(
        EditFormEntry{.container = e1, .widget = e1, .spanning = true});
}

void EditForm::add_list(
    const char* name, const char* label)
{
    auto b1 = new Fl_Flex(0, 0, 0, 120);
    auto p1 = new Fl_Box(0, 0, 0, 0);
    b1->fixed(p1, 20);
    auto e1 = new Fl_Select_Browser(0, 0, 0, 0, label);
    e1->end();
    b1->end();
    e1->align(FL_ALIGN_TOP_LEFT);
    e1->color(fl_rgb_color(255, 255, 200));

    //TODO
    e1->callback([](Fl_Widget* w, void* a) {
        std::cout << "Chosen: " << ((Fl_Select_Browser*) w)->value()
                  << std::endl;
    });

    //TODO--
    e1->add("First item");
    e1->add("Second item");
    e1->add("Third item");

    //e1->vertical_label_margin(5);

    entries.push_back(EditFormEntry{.container = b1,
                                    .widget = b1,
                                    .name = name,
                                    .spanning = true,
                                    .growable = true});
}

void EditForm::do_layout()
{
    int labwidth = 0;
    int wl, hl;
    int n_entries = entries.size();
    layout(n_entries, 2);
    for (int i = 0; i < n_entries; ++i) {
        auto e = entries[i];
        add(e.container);
        if (e.spanning) {
            widget(e.container, i, 0, 1, 2);
        } else {
            widget(e.container, i, 1);
        }
        if (!e.growable) {
            row_weight(i, 0);
        } else {
            row_gap(i - 1, 25);
        }

        wl = 0;

        //std::cout << "§ " << i << std::endl;
        //auto l = e.widget->label();
        //std::cout << "§+ " << l << std::endl;

        e.widget->measure_label(wl, hl);
        //std::cout << "§++ " << wl << std::endl;
        if (wl > labwidth)
            labwidth = wl;
    }
    //col_width(0, labwidth + 15);
    //col_width(0, 0);
    col_weight(0, 0);
    col_gap(0, labwidth + 15);
}

