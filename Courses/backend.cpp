#include "backend.h"

#include <QApplication>
#include <QJsonDocument>
#include "messages.h"

BackEnd *backend;

// The reader for responses from the back-end runs continuously in the
// background. It emits a signal when some data has been received.
// Currently a line (terminated by '\n') is expected to be a complete
// JSON object, but it should be possible to accept multiline objects,
// for example by using an alternative terminator (Ctrl-G / "\a" / BEL?).
// The JSON object is then parsed and made available by means of a
// "received_input" signal. Invalid input will cause a "received_invalid"
// signal to be emitted.
// This continues until the program is closed.

// Manage the back-end process, communicating with its stdio.
BackEnd::BackEnd(
    MainWindow *window1)
    : QObject()
    , mainwindow{window1}
{
    //TODO: This is temporary stuff for testing, for production use
    // it needs some work ...
    QString program = "./backend";
    QStringList arguments;
    //arguments << "source_file";
    //ODOT

    //QProcess::setReadChannel(QProcess::ProcessChannel channel)
    // -- QProcess::StandardOutput  QProcess::StandardError
    process = new QProcess(this);

    connect(process,
            &QProcess::readyReadStandardOutput,
            this,
            &BackEnd::handleBackendOutput);

    connect(process,
            &QProcess::readyReadStandardError,
            this,
            &BackEnd::handleBackendError);

    connect(process, &QProcess::finished, this, &BackEnd::finished);

    //connect(qApp, &QApplication::lastWindowClosed, this, &BackEnd::closing);

    process->start(program, arguments);

    backend = this;
}

void BackEnd::handleBackendOutput()
{
    QJsonParseError jerr;
    char ch;
    while (process->getChar(&ch)) {
        if (ch == '\n') {
            QJsonDocument jin = QJsonDocument::fromJson(linebuffer, &jerr);
            if (jerr.error == QJsonParseError::NoError) {
                if (jin.isObject()) {
                    linebuffer.clear();
                    auto jobj = jin.object();
                    mainwindow->received_input(jobj);
                    continue;
                }
                received_invalid(
                    QString("CALLBACK ERROR, not object\n:: " + linebuffer));
                linebuffer.clear();
                continue;
            }
            // else: JSON parse failed
            auto s = QString(linebuffer);
            auto st = s.trimmed();
            linebuffer.clear();
            if (st.size() == 0) {
                //TODO: Ignore this?
                qDebug() << "CALLBACK EMPTY";
            } else {
                received_invalid(QString("CALLBACK ERROR: ")
                                 + jerr.errorString() + "\n:: " + s);
            }
            continue;
        }

        // Not newline, add to linebuffer
        linebuffer.append(ch);
    }
}

void BackEnd::received_invalid(
    QString text)
{
    IgnoreError("INVALID INPUT", text);
}

void BackEnd::handleBackendError()
{
    auto bytes = process->readAllStandardError();
    IgnoreError("BACKEND ERROR", QString(bytes));
}

void BackEnd::finished()
{
    qDebug() << "Backend finished";
    //TODO: close window? or otherwise wind up ...
}

void BackEnd::call_backend(
    const QJsonObject data)
{
    QJsonDocument jdoc(data);
    QByteArray jbytes = jdoc.toJson(QJsonDocument::Compact) + '\n';
    qDebug() << "Sending:" << QString(jbytes);
    process->write(jbytes);
}
