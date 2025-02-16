#include "mainwindow.h"
#include <QFile>
#include <QUiLoader>
#include <QVBoxLayout>

MainWindow::MainWindow()
    : QWidget()
{
    QUiLoader loader;

    //TODO: This is probably not really the main window ...
    QFile file("ui/courses.ui");

    file.open(QFile::ReadOnly);
    QWidget *w = loader.load(&file, this);
    file.close();

    QVBoxLayout *layout = new QVBoxLayout;
    layout->addWidget(w);
    setLayout(layout);
}

void MainWindow::closeEvent(
    QCloseEvent *event)
{
    event->ignore();
}
