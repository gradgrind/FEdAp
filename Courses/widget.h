#ifndef WIDGET_H
#define WIDGET_H

#include <QWidget>

QWidget *loadUiFile(QString uifile, QWidget *parent = nullptr);
QDialog *loadDialogFile(QString uifile, QWidget *parent = nullptr);

#endif // WIDGET_H
