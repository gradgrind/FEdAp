#include "messages.h"
#include <QApplication>

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
