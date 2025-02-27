#include "callback.h"
#include <QApplication>
#include "messages.h"
#include <iostream>

/* For reference:
    QFile _file;
    _file.open(stdin, QIODevice::ReadOnly);
    // or _file.open(stdin, QIODevice::ReadOnly | QIODevice::Text);
    QByteArray line;
    forever {
        line = _file.readLine();
        qDebug() << "GOT:" << line << QString(line);
        if (line[0] == '\n') {
            break;
        }
    } 
*/

// The input reader runs in a separate thread, continually reading
// lines from standard input.
// Currently a line is expected to be a complete JSON object, but if
// multiline objects should be desired, it could be achieved by
// accumulating lines until some terminator is encountered, maybe
// Ctrl-G ("\a", BEL)?
// The JSON object is then parsed and made available by means of a
// "received_input" signal. Invalid input will cause a "received_invalid"
// signal to be emitted.
// This continues until the program is closed.

void ReadInputStream::run()
{
    std::string instring;
    QJsonParseError jerr;
    forever {
        getline(std::cin, instring);

        QJsonDocument jin = QJsonDocument::fromJson(instring.data(), &jerr);
        if (jerr.error == QJsonParseError::NoError) {
            if (jin.isObject()) {
                auto jobj = jin.object();
                emit received_input(jobj);
                continue;
            }
            emit received_invalid(QString::fromStdString(
                "CALLBACK ERROR, not object\n:: " + instring));
            continue;
        }
        auto s = QString::fromStdString(instring);
        auto st = s.trimmed();
        if (st.size() == 0) {
            //TODO: Ignore this?
            qDebug() << "CALLBACK EMPTY";
        } else {
            emit received_invalid(QString::fromStdString("CALLBACK ERROR: ")
                                  + jerr.errorString() + "\n:: " + s);
        }
    }
}

CallBackManager::CallBackManager()
    : QObject()
{
    //workerThread = new WorkerThread(this);
    connect(&workerThread,
            &ReadInputStream::received_input,
            this,
            &CallBackManager::handleResult);
    connect(&workerThread,
            &ReadInputStream::received_invalid,
            this,
            &CallBackManager::handleError);

    connect(qApp,
            &QApplication::lastWindowClosed,
            this,
            &CallBackManager::closing);

    workerThread.start();
}

void CallBackManager::handleResult(
    const QJsonObject jobj)
{
    qDebug() << jobj;
}

void CallBackManager::handleError(
    const QString text)
{
    qDebug() << "CallBackManager::handleError" << text;
    //IgnoreError("INVALID INPUT", text);
}

void CallBackManager::closing()
{
    qDebug() << "Closing";
    workerThread.terminate();
}

void CallBackManager::call_backend(
    const QJsonObject data)
{
    QJsonDocument jdoc(data);
    std::cout << jdoc.toJson(QJsonDocument::Compact).constData() << '\n';
    std::cout.flush();
}
