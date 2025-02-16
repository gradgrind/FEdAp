#include "backend.h"

#include <QApplication>
#include "messages.h"
//#include <iostream>

// The input reader runs in a separate thread, continually reading
// lines from the back-end process.
// Currently a line is expected to be a complete JSON object, but if
// multiline objects should be desired, it could be achieved by
// accumulating lines until some terminator is encountered, maybe
// Ctrl-G ("\a", BEL)?
// The JSON object is then parsed and made available by means of a
// "received_input" signal. Invalid input will cause a "received_invalid"
// signal to be emitted.
// This continues until the program is closed.

//TODO: Maybe I don't need a thread at all! Use signals for reading input.
ReadInput::ReadInput()
    : QThread()
{}

void ReadInput::run()
{
    //std::string instring;
    QByteArray inbytes;
    QJsonParseError jerr;
    forever {
        inbytes = process->readLine(); //TODO: Does this block???
        // Actually, to be able to read both stdout and stderr I would
        // need to use the signals, I guess – it looks like a switch
        // needs setting to select the read channel.

        QJsonDocument jin = QJsonDocument::fromJson(inbytes, &jerr);
        if (jerr.error == QJsonParseError::NoError) {
            if (jin.isObject()) {
                auto jobj = jin.object();
                emit received_input(jobj);
                continue;
            }
            emit received_invalid(
                QString("CALLBACK ERROR, not object\n:: " + inbytes));
            continue;
        }
        auto s = QString(inbytes);
        auto st = s.trimmed();
        if (st.size() == 0) {
            //TODO: Ignore this?
            qDebug() << "CALLBACK EMPTY";
        } else {
            emit received_invalid(QString("CALLBACK ERROR: ")
                                  + jerr.errorString() + "\n:: " + s);
        }
    }
}

// Manage the back-end process, communicating with its stdio.
BackEnd::BackEnd()
    : QObject()
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

    connect(qApp, &QApplication::lastWindowClosed, this, &BackEnd::closing);

    process->start(program, arguments);
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
                    received_input(jobj);
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

void BackEnd::received_input(
    QJsonObject jobj)
{}

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
    workerThread.process->write(jbytes);
}
