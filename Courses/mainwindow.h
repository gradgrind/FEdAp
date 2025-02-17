#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QCloseEvent>
#include <QWidget>

class MainWindow : public QWidget
{
    Q_OBJECT

    bool backend_running = false;

public:
    MainWindow(QString);

    virtual void received_input(QJsonObject);

    void closeEvent(QCloseEvent *) override;

public slots:
    void backend_started();
    void backend_finished();
};

#endif // MAINWINDOW_H
