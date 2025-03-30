#include "fltk_json.h"

using namespace std;

Flex::Flex(
    string name, bool horizontal)
    : Widget(name, new Fl_Flex(horizontal ? Fl_Flex::ROW : Fl_Flex::COLUMN))
{
    cout << "Flex " << this << endl;
}
