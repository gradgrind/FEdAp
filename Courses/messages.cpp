#include "messages.h"
#include <QApplication>
#include "backend.h"
#include "widget.h"

bool IgnoreError(
    QString text, QString info, QString details)
{
    QMessageBox msgBox;
    msgBox.setText(text);
    if (!info.isEmpty()) {
        msgBox.setInformativeText(info);
    }
    if (!details.isEmpty()) {
        msgBox.setDetailedText(details);
    }
    msgBox.setIcon(QMessageBox::Critical);
    msgBox.setStandardButtons(QMessageBox::Ignore | QMessageBox::Abort);
    //msgBox.setDefaultButton(QMessageBox::Abort);
    if (msgBox.exec() == QMessageBox::Ignore) {
        return true;
    }
    QApplication::quit();
    return false;
}

int TidyOnExit(
    QString info, QString details)
{
    QMessageBox msgBox;
    msgBox.setText("Quitting – save the data?");
    if (!info.isEmpty()) {
        msgBox.setInformativeText(info);
    }
    if (!details.isEmpty()) {
        msgBox.setDetailedText(details);
    }
    msgBox.setIcon(QMessageBox::Critical);
    msgBox.setStandardButtons(QMessageBox::Yes | QMessageBox::No
                              | QMessageBox::Cancel);
    msgBox.setDefaultButton(QMessageBox::Yes);
    auto result = msgBox.exec();
    if (result == QMessageBox::Cancel) {
        return 0;
    }
    if (result == QMessageBox::No) {
        return -1;
    }
    return 1;
}

WaitingDialog::WaitingDialog(
    QWidget *parent)
    : QDialog(parent)
{
    dialog = loadUiFile("background_monitor.ui", this);
    title_view = dialog->findChild<QLineEdit *>("title_view");
    progress = dialog->findChild<QLineEdit *>("progress");
    output = dialog->findChild<QPlainTextEdit *>("output");

    auto bb = dialog->findChild<QDialogButtonBox *>("buttonBox");
    ok_button = bb->button(QDialogButtonBox::Ok);
    cancel_button = bb->button(QDialogButtonBox::Cancel);

    connect(ok_button, &QPushButton::clicked, this, &WaitingDialog::handle_ok);
    connect(cancel_button,
            &QPushButton::clicked,
            this,
            &WaitingDialog::handle_cancel);
}

void WaitingDialog::start(
    QString title)
{
    title_view->setText(title);
    ok_button->setEnabled(false);
    cancel_button->setEnabled(true);
    progress->clear();
    output->clear();
    autoclose = true;
    timeout.start(200, this);
}

void WaitingDialog::timerEvent(
    QTimerEvent *event)
{
    timeout.stop();
    open();
}

void WaitingDialog::closeEvent(
    QCloseEvent *event)
{
    event->ignore();
}

void WaitingDialog::operation_cancelled()
{
    add_text("OPERATION_CANCELLED");
    done();
}

void WaitingDialog::done()
{
    timeout.stop();
    if (autoclose) {
        hide();
    } else {
        ok_button->setEnabled(true);
        cancel_button->setEnabled(false);
    }
}

void WaitingDialog::add_text(
    QString text)
{
    autoclose = false;
    output->appendPlainText(text + "\n");
}

void WaitingDialog::set_progress(
    QString text)
{
    progress->setText(text);
}

void WaitingDialog::handle_ok()
{
    hide();
}

void WaitingDialog::handle_cancel()
{
    autoclose = false;
    backend->cancel_current();
}
