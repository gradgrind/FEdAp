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

QDialog *loadDialogFile(
    QString uifile, QWidget *parent)
{
    QFile file("ui/" + uifile);
    file.open(QIODevice::ReadOnly);

    QUiLoader loader;
    return (QDialog *) loader.load(&file, parent);
}
