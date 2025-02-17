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
