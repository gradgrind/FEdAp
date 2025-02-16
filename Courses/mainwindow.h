#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QCloseEvent>
#include <QWidget>

class MainWindow : public QWidget
{
    Q_OBJECT
public:
    MainWindow();

    void closeEvent(QCloseEvent *) override;
};

#endif // MAINWINDOW_H
