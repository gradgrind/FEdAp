#include "widget.h"
#include <QtUiTools>

QWidget *loadUiFile(
    QString uifile, QWidget *parent)
{
    QFile file("ui/" + uifile);
    file.open(QIODevice::ReadOnly);

    QUiLoader loader;
    return loader.load(&file, parent);
}

Widget::Widget(
    QWidget *parent)
    : QWidget(parent)
{}

Widget::~Widget() {}
