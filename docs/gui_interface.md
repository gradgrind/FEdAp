# The GUI and its Interface to the Back-end

Separating the GUI from the data processing has several advantages, which in my view outweigh any disadvantages for this application.

By having a clearly specified interface, the front-end (GUI) and back-end (data processing) can be developed relatively independently of each other. They can also use completely different languages, as is here the case.

I decided I wanted a flexible interface which would not be tied too closely to the features and requirements of the two communicating sides. To this end I settled on a message-passing approach using JSON as the medium. There is no need for high-speed animations or high-volume data transfers, so the efficiency of the communication should not be critical.

Then there is the question of the level of the communication and the distribution of the logic between front-end and back-end. I decided to put as much of the logic as possible in the back-end, so that it should be quite straightforward to adopt a completely different front-end, even a browser-based one, if desired. Because GUI technology varies so widely it seemed sensible to keep the GUI code as simple as possible. This also means a sort of lower-common-denominator approach, but I think that is, on the whole, compatible with the aims of the application, even if "neater" ways of doing things might be possible by concentrating on one technology.

One aspect of GUIs that can be confusing, or at least unintuitive is the representation of the program state. If a GUI can cause the underlying data to change, it can be unclear whether the data shown in the GUI is a true mirror of the underlying data – at what moment do changes in the displayed data (e.g. by changing text in an editor or adding a new line in a table) become integrated in the data being processed? There is unfortunately no single, simple and satisfying answer to this.

The data could be stored in an SQL-database, or some non-SQL-database; it could be just a JSON or XML file. When the data is modified, the changes could be stored in the database immediately, or only when the user activates a "save" function. Changes made in the front-end could be passed directly to the back-end, even at the level of individual key-strokes, or they could be buffered and only sent when a "submit" function is activated (rather like the idea behind HTML forms).

I decided to start with a back-end written in Go, because of its simplicity, portability and compilation speed, in the hope that this would lead to an efficient and satisfying development experience, and also because I already had quite a lot of relevant code. For data storage I decided to start with fairly simple JSON files, with the ability to import from one or two other formats.

For the front-end I eventually settled on Qt. This is big and written in C++, but cross-platform and very capable. I reckoned that it would cover all my needs for widgets, etc., it has a built-in JSON library, process control, communication protocols, threads, and just about everything else. I also had some experience in using it. There is also Qt Creator to assist in the complicated build process, testing, etc. Maybe it would turn out later that something else would be better, but it seemed a good place to start. I also decided to try the form-submission approach to handling modifications, as far as possible.

## Front-end to Back-end

The front-end needs to communicate user actions to the back-end. This should basically be limited to a command-response type interaction, and limited to one action at a time. That means the GUI would essentially block further user input after a user action until the response from the back-end has been received. The front-end itself is event-driven and should not block, so unexpected user input should either be ignored or, perhaps better, reported as invalid (e.g. by means of a pop-up). If an action doesn't complete within a certain short period of time, a pop-up window should be shown in which report messages from the back-end can be displayed, concerning progress, errors or whatever. This can be a modal window, which would block interaction with the main GUI anyway.

In case the user should wish to abort a long-running back-end process, there should be a "cancel" button in the report window. This command would differ from normal commands, because it must be passed to the back-end even though the current process is still running.

The basic structure of a front-end command to the back-end would be {"DO":"command"}, where the command-string depends on the action initiated in the front-end. Any further information which may be necessary can be passed as key-value pairs.

## Back-end to Front-end

There are three basic types of message which may be passed: operation completion, report, GUI control.

### Operation completion

When processing of a command from the front-end is completed, this must be signalled back to the front-end by a simple {"DONE":true/false}, "true" if the operation completed successfully.

### Reports

There can be progress reports, information, warnings and errors (there may also be bug reports). The operative key here is "REPORT", its value depending on the type of message. A simple translation mechanism is supported to enable showing the messages in a supported language.

Normally the reports will be shown in the report window, which opens soon after an operation is initiated. However, if an operation is completed before this window appears it will not appear, unless a report above a certain level of urgency has been received. Also, the window will be closed automatically if the operation completes without such an important message. Errors and warnings will cause the window to appear and remain visible. There is also a more neutral type of information, "notice", which will likewise force the appearance; an "info" message will not.

### Controlling the GUI

As the bulk of the GUI control is to be handled by the back-end, quite low-level commands must be sent to the front-end. The basic structure of such a command is:

   {"GUI": "command","OBJECT": "object-name"}

Any parameters needed can be passed as further key-value pairs.

These commands are directed towards single GUI elements. They act as an isolation layer between the application logic and the details of the GUI-toolkit.