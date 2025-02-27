#include "messages.h"
#include <QApplication>
#include "backend.h"
#include "widget.h"

/* The dialogs here cover the reporting of back-end communication.
 * 
 * The primary one is [WaitingDialog], which pops up when a back-end
 * operation takes more than a certain time, indicating that the
 * processing is taking place and possibly providing feedback, while
 * blocking further command input.
 * If there is no feedback and the operation completes quickly, this
 * dialog will not appear. This dialog has a Cancel button, so that
 * the user can stop an operation before it completes, assuming the
 * operation is sufficently responsive – it must be listening for
 * interruptions. When a non-terminating operation is not responsive
 * the application must be forcibly closed.
 * TODO: It may be good to offer a Force-Closure operation, via a
 * pop-up which appears after a delay, but this should probably
 * disappear automatically if the operation does end after all ...
 * 
 * The [QuitUnsaved] dialog should appear when the application is
 * closed (e.g. with the "X" corner-button) and there is unsaved data.
 * It is a just a warning, the choices are to cancel the exit or to
 * exit losing the changes.
 * 
 * Normal error reports concerning the data will appear in the
 * [WaitingDialog] (and are possibly logged additionally in a file).
 * Abnormal error reports are probably connected with programm bugs
 * and appear in their own pop-up dialogs – [CriticalError]. Before
 * opening this dialog, any pending [WaitingDialog] should be suppressed,
 * to ensure that the bug report is the primary window. After accepting
 * such a report, the application should probably be forcibly closed.
*/

//TODO: Should there really be the possibility of ignoring the error and
// continuing? If so, what to do about any suppressed [WaitingDialog]?
//TODO: Handle [WaitingDialog] suppression. Maybe making all the
// messaging dialogs part of a single class would help?
bool CriticalError(
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

void WaitingDialog::done()
{
    if (autoclose) {
        timeout.stop();
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

void WaitingDialog::force_open()
{
    autoclose = false;
}
