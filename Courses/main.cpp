#include "backend.h"
//#include "callback.h"
//#include "mainwindow.h"
#include "testio.h"
//#include "widget.h"

#include <QApplication>

#include <QFile>

#include <QJsonDocument>

//TODO--
const char *testcmd = R"(
{
    "FirstName": "John",
    "LastName": "Doe",
    "Age": 43,
    "Address": {
        "Street": "Downing Street 10",
        "City": "London",
        "Country": "Great Britain"
    },
    "Phone numbers": [
        "+44 1234567",
        "+44 2345678"
    ]
}
)";

int main(
    int argc, char *argv[])
{
    QApplication a(argc, argv);

    MainWindow w("courses_1.ui");
    //TestIo w;
    const QList<QPushButton *> allPButtons = w.findChildren<QPushButton *>();
    for (const auto pb : allPButtons) {
        qDebug() << pb->objectName();
    }

    BackEnd cbman(&w);
    //CallBackManager cbman;

    //QJsonParseError jerr;
    //QJsonDocument jcmd = QJsonDocument::fromJson(testcmd, &jerr);
    //backend->call_backend(jcmd.object());

    w.show();
    //QWidget *w = loadUiFile("courses.ui");
    //w->show();
    return a.exec();
}
