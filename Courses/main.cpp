#include "widget.h"

#include <QApplication>

int main(
    int argc, char *argv[])
{
    QApplication a(argc, argv);
    QWidget *w = loadUiFile("courses.ui");
    w->show();
    return a.exec();
}
