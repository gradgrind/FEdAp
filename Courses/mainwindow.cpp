#include "mainwindow.h"
#include <QFile>
#include <QUiLoader>
#include <QVBoxLayout>

#include <QJsonObject>

MainWindow::MainWindow(
    QString uifile)
    : QWidget()
{
    QUiLoader loader;

    QFile file(QString("ui/") + uifile);

    file.open(QFile::ReadOnly);
    QWidget *w = loader.load(&file, this);
    file.close();

    QVBoxLayout *layout = new QVBoxLayout;
    layout->addWidget(w);
    setLayout(layout);
}

void MainWindow::received_input(
    QJsonObject jobj)
{
    qDebug() << "Received:" << jobj;
}

void MainWindow::closeEvent(
    QCloseEvent *event)
{
    event->ignore();

    QJsonObject quitcmd{{"QUIT", 0}};
}
