#ifndef BACKEND_H
#define BACKEND_H

#include <QJsonDocument>
#include <QJsonObject>
#include <QProcess>
#include <QThread>

class ReadInput : public QThread
{
    Q_OBJECT
    void run() override;

public:
    ReadInput();

    QProcess* process;

signals:
    void received_input(const QJsonObject);
    void received_invalid(const QString);
};

class BackEnd : public QObject
{
    Q_OBJECT
    ReadInput workerThread;

public:
    BackEnd();

    void call_backend(const QJsonObject);

private slots:
    void handleResult(const QJsonObject);
    void handleError(const QString);
    void closing();
    void finished();
};

#endif // BACKEND_H
