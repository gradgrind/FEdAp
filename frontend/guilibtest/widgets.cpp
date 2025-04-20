#include "minion.h"
#include <FL/Fl_Box.H>
using namespace std;
using mmap = minion::MinionMap;

// *** Non-layout widgets ***

Fl_Widget *NEW_Box(
    mmap param)
{
    return new Fl_Box(0, 0, 0, 0);
}
