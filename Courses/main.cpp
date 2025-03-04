#include "backend.h"
//#include "callback.h"
//#include "mainwindow.h"
//#include "testio.h"
//#include "widget.h"

#include <QApplication>

#include <QFile>

#include <QJsonDocument>

//TODO--
const char *testcmd = R"(
{
    "DO": "LOAD_W365_JSON"
}
)";

class EvHandler : public QObject
{
    //QObject *wtarget;

public:
    /*
    EvHandler(
        QObject *target)
        : QObject()
    {
        wtarget = target;
        target->installEventFilter(this);
    }
*/
    bool eventFilter(
        QObject *object, QEvent *event)
    {
        //if (object == wtarget && event->type() == QEvent::MouseButtonDblClick) {
        if (event->type() == QEvent::MouseButtonDblClick) {
            qDebug() << "Double Click" << object->objectName();
            return true;
        }
        if (event->type() == QEvent::KeyPress) {
            QKeyEvent *keyEvent = static_cast<QKeyEvent *>(event);
            if (keyEvent->key() == Qt::Key_Return) {
                qDebug() << "Ate return key press on" << object->objectName();
                return true;
            }
        }

        // standard event processing
        return QObject::eventFilter(object, event);
    }
};

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

    auto lect = w.findChild<QLineEdit *>("choose_teachers");
    auto lecg = w.findChild<QLineEdit *>("choose_groups");
    auto evh = EvHandler();
    lect->installEventFilter(&evh);
    lecg->installEventFilter(&evh);
    BackEnd cbman(&w);
    //CallBackManager cbman;

    QJsonParseError jerr;
    QJsonDocument jcmd = QJsonDocument::fromJson(testcmd, &jerr);
    backend->call_backend(jcmd.object());

    w.show();
    //QWidget *w = loadUiFile("courses.ui");
    //w->show();
    return a.exec();
}
