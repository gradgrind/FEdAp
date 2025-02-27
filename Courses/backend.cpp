#include "backend.h"

#include <QApplication>
#include <QJsonDocument>
#include "messages.h"

/* Messages – commands/requests – are sent to the back-end process as JSON
 * objects. The key "DO" specifies the operation, its value is a string.
 * The parameters can take any form.
 * Also the result is a JSON object, the value with key "DONE" specifying
 * how it should be handled. Again, the parameters can take any form. The
 * back-end can, however, also send reports before completion of the
 * operation. In this case, the "DONE" value should be "" and the value
 * with key "REPORT" indicates what is to be done with the report.
 * At present the supported possibilities are:
 *   "PROGRESS" – the text field with key "TEXT" will be used to set the
 *   progress widget,
 *   "REPORT" –  the text field with key "TEXT" will be added to the
 *   report widget,
 *   "QUIT_UNSAVED?" – the back-end has unsaved data, request confirmation
 *   that the changes should be discarded.
 *   "BACKEND_BUSY" – the operation available via the "DATA" key was
 *   received by the backend while it was processing another command.
 * The whole process is event-driven (using signals and slots), so that
 * the GUI doesn't hang. If operations take longer than a brief time-out
 * interval, a modal pop-up will appear allowing cancellation and showing
 * feedback from the back-end process.
 * In general, only one operation may be active at any time. This is
 * managed in the BackEnd class, which will (normally) not allow an
 * operation to start until a previous one has completed. There may,
 * however, be a small number of operations which can be started while
 * another operation is running. This can be useful for cancelling a
 * long-running operation, for example.
 */

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

    waiting_dialog = new WaitingDialog(window1);

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

    connect(process,
            &QProcess::started,
            mainwindow,
            &MainWindow::backend_started);
    connect(process,
            &QProcess::finished,
            mainwindow,
            &MainWindow::backend_finished);

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
                    auto lb = linebuffer;
                    linebuffer.clear();
                    auto jobj = jin.object();
                    auto doneval = jobj.value("DONE").toString();
                    if (doneval == "") {
                        auto rp = jobj.value("REPORT").toString();
                        if (rp == "PROGRESS") {
                            waiting_dialog->set_progress(
                                jobj.value("TEXT").toString());

                            //TODO: Translate the message tags?

                        } else if (rp == "Error" || rp == "Warning"
                                   || rp == "Notice") {
                            waiting_dialog->force_open();
                            waiting_dialog->add_text(
                                "*" + rp + "* " + jobj.value("TEXT").toString());
                        } else if (rp == "Info") {
                            waiting_dialog->add_text(
                                "*" + rp + "* " + jobj.value("TEXT").toString());
                        } else if (rp == "Bug") {
                            QMessageBox ::critical(mainwindow,
                                                   "BUG",
                                                   jobj.value("TEXT").toString());
                        } else if (rp == "QUIT_UNSAVED?") {
                            if (QMessageBox ::warning(
                                    mainwindow,
                                    "WARNING",
                                    "LOSE_CHANGES?",
                                    QMessageBox::StandardButtons(
                                        QMessageBox::Ok | QMessageBox::Cancel))
                                == QMessageBox::Ok) {
                                quit(true);
                            }
                        } else if (rp == "BACKEND_BUSY") {
                            QMessageBox::critical(waiting_dialog,
                                                  "BACKEND_BUSY",
                                                  lb);
                        } else {
                            QMessageBox::critical(waiting_dialog,
                                                  "BACKEND_ERROR",
                                                  lb);
                        }
                        continue;
                    } else {
                        waiting_dialog->done();
                        current_operation = QJsonObject();
                    }
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
    qDebug() << "INVALID INPUT:" << text;
    //IgnoreError("INVALID INPUT", text);
}

void BackEnd::handleBackendError()
{
    auto bytes = process->readAllStandardError();
    qDebug() << "BACKEND ERROR:" << QString(bytes);
    //IgnoreError("BACKEND ERROR", QString(bytes));
}

// Send a normal command to the back-end, start the dialog with timer
// to delay the appearance of the progress window, which might not
// appear at all if the the operation is quick enough.
void BackEnd::call_backend(
    const QJsonObject data)
{
    if (!current_operation.empty()) {
        //TODO:
        QMessageBox::critical(waiting_dialog,
                              "BACKEND_OPERATION",
                              "STILL_RUNNING");
        return;
    }
    QJsonDocument jdoc(data);
    QByteArray jbytes = jdoc.toJson(QJsonDocument::Compact) + '\n';
    qDebug() << "Sending:" << QString(jbytes);

    // Start dialog (with timer)
    waiting_dialog->start(data.value("DO").toString());
    process->write(jbytes);
}

void BackEnd::cancel_current()
{
    qDebug() << "TODO: BackEnd::cancel_current";
    QJsonObject jobj{{"DO", "CANCEL"}};
    QJsonDocument jdoc(jobj);
    QByteArray jbytes = jdoc.toJson(QJsonDocument::Compact) + '\n';
    qDebug() << "Sending:" << QString(jbytes);
    process->write(jbytes);
}

void BackEnd::quit(
    bool force)
{
    qDebug() << "TODO: BackEnd::quit";
    QJsonObject jobj{{"DO", "QUIT"}, {"FORCE", force}};
    QJsonDocument jdoc(jobj);
    QByteArray jbytes = jdoc.toJson(QJsonDocument::Compact) + '\n';
    qDebug() << "Sending:" << QString(jbytes);
    process->write(jbytes);
}
