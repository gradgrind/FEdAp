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

    void received_input(QJsonObject);
    void received_invalid(QString);

public:
    BackEnd();

    void call_backend(const QJsonObject);

    QProcess* process;
    QByteArray linebuffer;

private slots:
    void handleBackendOutput();
    void handleBackendError();

    //void handleResult(const QJsonObject);
    //void handleError(const QString);
    void finished();
};

#endif // BACKEND_H
