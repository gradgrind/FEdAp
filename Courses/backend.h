#ifndef BACKEND_H
#define BACKEND_H

#include "mainwindow.h"
#include "messages.h"

#include <QJsonObject>
#include <QProcess>

class BackEnd : public QObject
{
    Q_OBJECT

    void received_invalid(QString);

    MainWindow* mainwindow;
    QProcess* process;
    QByteArray linebuffer;
    QJsonObject current_operation;

public:
    BackEnd(MainWindow*);

    void call_backend(const QJsonObject);
    void cancel_current();
    void quit(bool force = false);

    WaitingDialog* waiting_dialog;

private slots:
    void handleBackendOutput();
    void handleBackendError();
};

extern BackEnd* backend;

//TODO?
class SlotHandler : public QObject
{
    Q_OBJECT

public:
    SlotHandler()
        : QObject()
    {}

public slots:

    void comboboxSelectionChanged(
        int i)
    {
        auto o = sender();
        qDebug() << "$" << i << o->objectName();
    }
};

#endif // BACKEND_H
