#include "editform.h"
#include "layout.h"
#include "widget_methods.h"
#include "widgets.h"
#include <FL/Fl_Flex.H>
#include <FL/fl_draw.H>
#include <iostream>
#include <ostream>
using namespace std;
using namespace minion;

//TODO: Do I really need the EditFormEntry struct as it is?
// How dynamic should the widget be? Should I be able to add entries
// later? Remove them?

// I could have the entries as a constructor parameter and use the single
// layout function if there is to be no change later.

//TODO: How to do the callbacks? Entry name + current value?
//TODO: How to write to the entries?
//  Entry_name + new_value (ENTRY)
//  Entry_name + [new_value ...] (LIST)

struct EditFormEntryData : public Fl_Callback_User_Data
{
    std::string w_name;

    EditFormEntryData(
        std::string_view n)
        : Fl_Callback_User_Data()
        , w_name{n}
    {}
};

void editform_method(
    Fl_Widget* w, string_view c, MinionList m)
{
    auto efw = static_cast<EditForm*>(w);
    int entry_n = efw->entries.size();

    //TODO ...

    if (c == "ENTRY") {
        string name = get<string>(m.at(1));
        string label = get<string>(m.at(2));
        auto e1 = new Fl_Output(0, 0, 0, efw->entry_height);
        efw->entry_map.emplace(name, e1);
        efw->add(e1);
        e1->copy_label(label.c_str());
        e1->align(FL_ALIGN_LEFT);
        e1->color(efw->entry_bg);

        //TODO
        e1->callback(
            [](Fl_Widget* w, void* ud) {
                cout << "Activated: " << static_cast<EditFormEntryData*>(ud)->w_name << ": "
                     << Fl::callback_reason() << endl;
            },
            new EditFormEntryData(name),
            true);

        e1->clear_visible_focus(); // no cursor, but text cannot be copied
        e1->horizontal_label_margin(5);

        //??
        EditFormEntry efe{
            .widget = e1,
            .name = name,
        };
        efw->entries.emplace_back(efe);

        efw->layout(entry_n + 1, 2);
        efw->widget(e1, entry_n, 1);
        efw->row_weight(entry_n, 0);
        if (entry_n == 0)
            efw->col_weight(0, 0);

        int wl, hl;
        e1->measure_label(wl, hl);
        if (wl > efw->label_width) {
            efw->label_width = wl;
            efw->col_gap(0, efw->label_width + 15);
        }

    } else if (c == "SEPARATOR") {
        auto e1 = new Fl_Box(FL_BORDER_FRAME, 0, 0, 0, 1, "");
        efw->add(e1);
        efw->entries.emplace_back(EditFormEntry{.widget = e1, .spanning = true});
        efw->layout(entry_n + 1, 2);
        efw->widget(e1, entry_n, 0, 1, 2);
        if (entry_n == 0)
            efw->col_weight(0, 0);
        efw->row_weight(entry_n, 0);

    } else if (c == "LIST") {
        string name = get<string>(m.at(1));
        string label = get<string>(m.at(2));
        //TODO
        auto e1 = new Fl_Select_Browser(0, 0, 0, 0);
        e1->end();
        efw->entry_map.emplace(name, e1);
        efw->add(e1);
        e1->copy_label(label.c_str());
        e1->align(FL_ALIGN_TOP_LEFT);
        e1->color(efw->entry_bg);

        //TODO
        e1->callback([](Fl_Widget* w, void* a) {
            cout << "Chosen: " << static_cast<Fl_Select_Browser*>(w)->value() - 1 << endl;
        });

        //TODO--
        e1->add("First item");
        e1->add("Second item");
        e1->add("Third item");

        //e1->vertical_label_margin(5);

        auto efe = EditFormEntry{.widget = e1,
                                 .name = name,
                                 .padabove = efw->list_title_space,
                                 .spanning = true,
                                 .growable = true};

        efw->entries.emplace_back(efe);

        efw->layout(entry_n + 1, 2);
        efw->widget(e1, entry_n, 0, 1, 2);
        if (entry_n == 0)
            efw->col_weight(0, 0);
        if (efe.padabove != 0) {
            if (entry_n == 0) {
                efw->margin(-1, efe.padabove);
            } else {
                efw->row_gap(entry_n - 1, efe.padabove);
            }
        }

    } else if (c == "SET") {
        string name = get<string>(m.at(1));
        auto e0 = efw->entry_map.at(name);
        if (auto e1 = dynamic_cast<Fl_Output*>(e0)) {
            auto v = get<string>(m.at(2)).c_str();
            e1->value(v);
            // Setting the tooltip allows overlong texts to be readable
            e1->tooltip(v);
        } else if (auto e1 = dynamic_cast<Fl_Select_Browser*>(e0)) {
            e1->clear();
            for (int i = 2; i < m.size(); ++i) {
                //TODO? add can have a second argument (void *), which can
                // refer to data ...
                e1->add(get<string>(m.at(i)).c_str());
            }
        }
        //TODO?
    } else {
        widget_method(w, c, m);
    }
}

EditForm::EditForm()
    : Fl_Grid(0, 0, 0, 0)
{
    Fl_Group::current(0); // disable "auto-grouping"
    box(FL_BORDER_FRAME);
    gap(10, 5);
    margin(5, 5, 5, 5);
}

Fl_Widget* container;
Fl_Widget* widget;
std::string name;
bool spanning;
bool growable;

void EditForm::add_separator()
{
    auto e1 = new Fl_Box(FL_BORDER_FRAME, 0, 0, 0, 1, "");
    entries.emplace_back(EditFormEntry{.widget = e1, .spanning = true});
}

void EditForm::add_list(
    const char* name, const char* label)
{
    auto e1 = new Fl_Select_Browser(0, 0, 0, 0, label);
    e1->end();
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

    entries.emplace_back(EditFormEntry{.widget = e1,
                                       .name = name,
                                       .padabove = 30,
                                       .spanning = true,
                                       .growable = true});
}

void EditForm::do_layout()
{
    int labwidth = 0; // for measuring the max. label width
    int wl, hl;
    int n_entries = entries.size();
    layout(n_entries, 2);
    for (int i = 0; i < n_entries; ++i) {
        auto e = entries[i];
        add(e.widget);
        if (e.spanning) {
            widget(e.widget, i, 0, 1, 2);
        } else {
            widget(e.widget, i, 1);
        }
        if (!e.growable) {
            row_weight(i, 0);
        }
        if (e.padabove != 0) {
            if (i == 0) {
                margin(-1, e.padabove);
            } else {
                row_gap(i - 1, e.padabove);
            }
        }

        wl = 0;
        e.widget->measure_label(wl, hl);
        if (wl > labwidth)
            labwidth = wl;
    }
    col_weight(0, 0);
    col_gap(0, labwidth + 15);
}

Fl_Widget* NEW_EditForm(
    MinionMap param)
{
    auto ef = new EditForm();
    string clr{"ffffc8"};
    param.get_string("ENTRY_BG", clr);
    ef->entry_bg = get_colour(clr);
    //ef->col_weight(0, 0);
    return ef;
}
