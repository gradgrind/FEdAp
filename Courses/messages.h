#ifndef MESSAGES_H
#define MESSAGES_H

#include <QBasicTimer>
#include <QLineEdit>
#include <QMessageBox>
#include <QPlainTextEdit>
#include <QPushButton>

bool IgnoreError(QString text, QString info = "", QString details = "");

int TidyOnExit(QString info = "", QString details = "");

class WaitingDialog : public QDialog
{
    Q_OBJECT

    QWidget* dialog;
    QLineEdit* title_view;
    QLineEdit* progress;
    QPlainTextEdit* output;
    QPushButton* ok_button;
    QPushButton* cancel_button;
    bool autoclose;
    QBasicTimer timeout;

protected:
    virtual void closeEvent(QCloseEvent*) override;
    virtual void timerEvent(QTimerEvent*) override;

public:
    WaitingDialog(QWidget* parent);

    void start(QString title);
    void done();
    void add_text(QString text);
    void set_progress(QString text);
    void operation_cancelled();
    void force_open();

private slots:
    void handle_ok();
    void handle_cancel();
};

#endif // MESSAGES_H
