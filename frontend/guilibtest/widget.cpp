#include "fltk_json.h"
#include <iostream>

Widget::Widget(
    std::string_view name)
    : wname{name}
{
    // Allow unnamed widgets. These are not placed in the map.
    if (wname.empty())
        return;

    if (WidgetMap.contains(wname)) {
        throw fmt::format("Widget name already exists: %s", wname);
    }

    //TODO: I think this may be wrong ... maybe this needs to be
    // done in the main constructor?
    cout << "Widget.this " << this << endl;
    //WidgetMap.insert({name, this}); ... or ...
    //WidgetMap.emplace(name, this);
}

Widget::~Widget()
{
    // Allow unnamed widgets. These are not placed in the map.
    if (wname.empty())
        return;
    WidgetMap.erase(wname);
}

void *Widget::get(
    std::string_view name)
{
    return WidgetMap.at(name);
}

/* seems not to work ...
const string_view Widget::widget_name()
{
    return wname;
}
*/
