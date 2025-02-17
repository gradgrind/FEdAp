#ifndef BACKEND_H
#define BACKEND_H

#include "mainwindow.h"

#include <QJsonObject>
#include <QProcess>

class BackEnd : public QObject
{
    Q_OBJECT

    void received_invalid(QString);

    MainWindow* mainwindow;
    QProcess* process;
    QByteArray linebuffer;

public:
    BackEnd(MainWindow*);

    void call_backend(const QJsonObject);

private slots:
    void handleBackendOutput();
    void handleBackendError();
};

extern BackEnd* backend;

#endif // BACKEND_H
