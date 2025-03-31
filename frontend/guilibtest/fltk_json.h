#ifndef FLTK_JSON_H
#define FLTK_JSON_H

#include <FL/Fl_Flex.H>
#include <fmt/format.h>
#include <functional>
#include <string>
#include <unordered_map>

#include <json.hpp>
using json = nlohmann::json;

using namespace std;

// The user_data in Fl_Widget is used to store additional fields,
// primarily the widget's name and type. As this feature is now no
// longer available to the programmer, another field (user_data) is
// provided within the WidgetData subclass of Fl_Callback_User_Data
// which is referred to by the original user data field.

extern unordered_map<string_view, Fl_Widget *> widget_map;

class WidgetData : public Fl_Callback_User_Data
{
    // Widget name, used for look-up, etc.
    string wname;
    // Widget type, which can be used to access a type's member
    // functions, also the name of the type.
    int wtype;
    // Substitute for Fl_Widget's user_data
    void *user_data;
    bool auto_delete_user_data;

public:
    WidgetData(std::string_view type, std::string_view name, Fl_Widget *widget);
    ~WidgetData() override;

    void add_widget(Fl_Widget *w);
    void remove_widget(std::string_view name);

    string_view widget_name();
    int widget_type();
    string_view widget_type_name();
};

void newFlex(string_view name, json data);

// ***

// Maybe the user_data in Fl_Widget can do what I need, which is
// to provide extra data fields. I could use a Fl_Callback_User_Data
// subclass, which has a virtual destructor.

// The reason I was using classes with multiple inheritance was
// primarily so that deletion of a widget (especially as a result
// of deleting its group) could allow removal of the Widget entry
// in the widget table, and releasing of its space.
// Without this there can be dangling pointers.
// Another advantage is that the fltk widget has effectively access
// to the full Widget data.
// HOWEVER, I couldn't get that to work properly ... I guess it
// might have something to do with the structure of the class
// representations, casting wasn't doing what I might have expected.

// Could something like this work?

/* This class handles mapping widget names to their actual widgets.
 * Note that the keys are string_view, so the underlying string is
 * expected to be held in the widget.
 */

class WidgetMap : public unordered_map<string_view, Fl_Widget *>
{
public:
    void add(
        string_view name, Fl_Widget *widget)
    {
        // Allow unnamed widgets. These are not placed in the map.
        if (name.empty())
            return;

        if (contains(name)) {
            throw fmt::format("Widget name already exists: %s", name);
        }
        //insert({name, widget}); ... or ...
        emplace(name, widget);
    }

    void remove(
        string_view name)
    {
        // Allow unnamed widgets. These are not placed in the map.
        if (name.empty())
            return;
        erase(name);
    }

    Fl_Widget *get(
        string_view name)
    {
        return at(name);
    }
};


class _Flex : public Fl_Flex
{
    std::string wname;

public:
    //TODO: This won't be generally available. As it stands, it
    // is necessary to cast the object to (in this case) _Flex*.
    constexpr static const std::string_view wtype{"Flex"};

    static _Flex *make(string_view name, json data);

    _Flex(string_view name, bool horizontal);

    ~_Flex();
};

// ---

class Gui
{
public:
    inline static unordered_map<string_view, function<void(string_view, json)>>
        FunctionMap;
};

class Widget
{
protected:
    std::string wname;
    inline static unordered_map<string_view, void *> WidgetMap;

    Widget(std::string_view name);
    virtual ~Widget();

public:
    virtual const string_view widget_type() = 0;
    //?? const string_view widget_name();
    static void *get(std::string_view name);
};

//TODO???: Change the constructors to take "standard" arguments, so that
// no extra constructor-functions are necessary? With parent as
// widget argument? Or "" ... using the automatic add-to-group, or
// blocking that feature (set current group to 0). Other parameters
// as JSON?

class DoubleWindow : public Widget
{
    inline static const std::string wtype{"Window:Double"};

public:
    DoubleWindow(string_view name, int width = 800, int height = 600);
    const string_view widget_type();
};

void newDoubleWindow(string_view name, json data);

class Flex : public Fl_Flex, Widget
{
    inline static const std::string wtype{"Flex"};

public:
    static Flex *make(string_view name, json data);
    Flex(string_view name, bool horizontal);
    const string_view widget_type() override;
};

// In this case, the widget (name) is the new name
void newFlex(string_view name, json data);

#endif // FLTK_JSON_H
