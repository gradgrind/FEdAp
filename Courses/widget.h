#ifndef WIDGET_H
#define WIDGET_H

#include <QWidget>

QWidget *loadUiFile(QString uifile, QWidget *parent = nullptr);

class Widget : public QWidget
{
    Q_OBJECT

public:
    Widget(QWidget *parent = nullptr);
    ~Widget();
};
#endif // WIDGET_H
