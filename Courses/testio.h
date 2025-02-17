#ifndef TESTIO_H
#define TESTIO_H

#include <QLineEdit>
#include <QPlainTextEdit>
#include "mainwindow.h"

class TestIo : public MainWindow
{
    QLineEdit* sender;
    QPlainTextEdit* sent;
    QPlainTextEdit* receiver;

public:
    TestIo();

    virtual void received_input(QJsonObject) override;

private slots:
    void sendJson();
};

#endif // TESTIO_H
