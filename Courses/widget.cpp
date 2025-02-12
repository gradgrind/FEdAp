#include "widget.h"
#include <QtUiTools>

QWidget *loadUiFile(
    QWidget *parent)
{
    QFile file("courses.ui");
    file.open(QIODevice::ReadOnly);

    QUiLoader loader;
    return loader.load(&file, parent);
}

Widget::Widget(
    QWidget *parent)
    : QWidget(parent)
{}

Widget::~Widget() {}
