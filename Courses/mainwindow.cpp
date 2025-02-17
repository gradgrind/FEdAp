#include "mainwindow.h"
#include <QFile>
#include <QUiLoader>
#include <QVBoxLayout>
#include "backend.h"
#include "messages.h"

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

void MainWindow::backend_started()
{
    backend_running = true;
}

void MainWindow::backend_finished()
{
    backend_running = false;
    close();
}

void MainWindow::received_input(
    QJsonObject jobj)
{
    qDebug() << "Received:" << jobj;
    auto jquit = jobj.value("QUIT");
    qDebug() << "Quit?" << jquit;
    if (jquit != QJsonValue::Undefined) {
        int q = TidyOnExit();
        if (q != 0) {
            QJsonObject jsave{{"QUIT", q}};
            backend->call_backend(jsave);
        }
    }
}

void MainWindow::closeEvent(
    QCloseEvent *event)
{
    if (backend_running) {
        event->ignore();
        QJsonObject quitcmd{{"QUIT", 0}};
        backend->call_backend(quitcmd);
    }
}
