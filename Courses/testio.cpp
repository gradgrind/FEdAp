#include "testio.h"
#include <QJsonDocument>
#include "backend.h"
#include "messages.h"

TestIo::TestIo()
    : MainWindow("testio.ui")
{
    sender = findChild<QLineEdit *>("sender");
    sent = findChild<QPlainTextEdit *>("sent");
    receiver = findChild<QPlainTextEdit *>("receiver");

    connect(sender, &QLineEdit::returnPressed, this, &TestIo::sendJson);
}

void TestIo::sendJson()
{
    auto text = sender->text();
    sender->clear();

    QJsonParseError jerr;
    QJsonDocument jin = QJsonDocument::fromJson(text.toUtf8(), &jerr);

    if (jerr.error == QJsonParseError::NoError) {
        if (jin.isObject()) {
            backend->call_backend(jin.object());
            sent->appendPlainText("----------------------\n");
            sent->appendPlainText(jin.toJson());
            return;
        }
        qDebug() << "CALLBACK ERROR, not object:" << text;
        //IgnoreError("INVALID INPUT", "CALLBACK ERROR, not object\n:: " + text);
        return;
    }
    // else: JSON parse failed
    auto s = text.trimmed();
    if (s.size() == 0) {
        //TODO: Ignore this?
        qDebug() << "CALLBACK EMPTY";
    } else {
        qDebug() << "CALLBACK ERROR:" << jerr.errorString() << "//" << text;
        //IgnoreError("CALLBACK ERROR: " + jerr.errorString() + "\n:: " + text);
    }
}

void TestIo::received_input(
    QJsonObject jobj)
{
    qDebug() << "TestIo received:" << jobj;

    MainWindow::received_input(jobj);
    receiver->appendPlainText("----------------------\n");
    receiver->appendPlainText(QJsonDocument(jobj).toJson());
}
