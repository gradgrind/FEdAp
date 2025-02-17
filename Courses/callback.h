#ifndef CALLBACK_H
#define CALLBACK_H

#include <QJsonDocument>
#include <QJsonObject>

#include <QBasicTimer>
#include <QThread>

class ReadInputStream : public QThread
{
    Q_OBJECT

    void run() override;

signals:
    void received_input(const QJsonObject);
    void received_invalid(const QString);
};

class CallBackManager : public QObject
{
    Q_OBJECT
    ReadInputStream workerThread;

public:
    CallBackManager();

    static void call_backend(const QJsonObject);

private slots:
    void handleResult(const QJsonObject);
    void handleError(const QString);
    void closing();
};

#endif // CALLBACK_H
