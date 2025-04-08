#include "minion.h"
#include <chrono>
#include <fmt/format.h>
#include <fstream>
#include <iostream>
using namespace std;
using namespace std::chrono;

// *** Reading to nlohmann json object
// The minion reader seems quicker for small inputs, but slower for
// larger ones, even though the json is a lot larger than the minion.
//  Unoptimized, the difference is not so great (more for small inputs),
// but optimized nlohmann is about twice as fast for large inputs.
// By using strings as input, the time for I/O is factored out.

void readfile(
    string &data, const string &filepath)
{
    std::ifstream file(filepath);

    if (file) {
        data.assign((istreambuf_iterator<char>(file)), istreambuf_iterator<char>());
    } else {
        cerr << "Error opening file: " << filepath << endl;
    }
}

void writefile(
    const string &data, const string &filepath)
{
    std::ofstream file(filepath);
    if (file) {
        file << data;
        file.close();
    } else {
        cerr << "Error opening file: " << filepath << endl;
    }
}

void testminion1(
    const string &filepath)
{
    string idata{};
    readfile(idata, filepath);

    cout << "FILE: " << filepath << endl;

    // Use auto keyword to avoid typing long
    // type definitions to get the timepoint
    // at this instant use function now()
    auto start = high_resolution_clock::now();

    Minion::MinionParser mp(idata);

    // After function call
    auto stop = high_resolution_clock::now();

    // Subtract stop and start timepoints and
    // cast it to required unit. Predefined units
    // are nanoseconds, microseconds, milliseconds,
    // seconds, minutes, hours. Use duration_cast()
    // function.
    auto duration = duration_cast<microseconds>(stop - start);

    // To get the value of duration use the count()
    // member function on the duration object
    cout << "TIME: " << duration.count() << " microseconds" << endl;

    string odata;
    if (mp.error_message.empty()) {
        mp.to_json(odata, false);

        auto p = filepath.rfind(".");
        string f;
        if (p == string::npos) {
            f = filepath;
        } else {
            f = filepath.substr(0, p);
        }
        string f1 = f + ".json";
        writefile(odata, f);
        cout << " --> " << f << endl;
    } else {
        cout << "ERROR:\n" << mp.error_message << endl;
        return;
    }

    // Compare parsing with nlohmann
    start = high_resolution_clock::now();
    json data = json::parse(odata);
    stop = high_resolution_clock::now();
    duration = duration_cast<microseconds>(stop - start);
    cout << "TIME json: " << duration.count() << " microseconds" << endl;
}

void testminion()
{
    testminion1("_data/test0.minion");
    testminion1("_data/test1.minion");
    testminion1("_data/test2.minion");
}

namespace Minion {

// Convert a unicode code point (as hex string) to a UTF-8 string
bool unicode_utf8(
    string &utf8, const string &unicode)
{
    // Convert the unicode to an integer
    unsigned int code_point;
    stringstream ss;
    ss << hex << unicode;
    ss >> code_point;

    // Convert the code point to a UTF-8 string
    if (code_point <= 0x7F) {
        utf8 += static_cast<Char>(code_point);
    } else if (code_point <= 0x7FF) {
        utf8 += static_cast<Char>((code_point >> 6) | 0xC0);
        utf8 += static_cast<Char>((code_point & 0x3F) | 0x80);
    } else if (code_point <= 0xFFFF) {
        utf8 += static_cast<Char>((code_point >> 12) | 0xE0);
        utf8 += static_cast<Char>(((code_point >> 6) & 0x3F) | 0x80);
        utf8 += static_cast<Char>((code_point & 0x3F) | 0x80);
    } else if (code_point <= 0x10FFFF) {
        utf8 += static_cast<Char>((code_point >> 18) | 0xF0);
        utf8 += static_cast<Char>(((code_point >> 12) & 0x3F) | 0x80);
        utf8 += static_cast<Char>(((code_point >> 6) & 0x3F) | 0x80);
        utf8 += static_cast<Char>((code_point & 0x3F) | 0x80);
    } else {
        // Invalid input
        return false;
    }
    return true;
}

/* Generate a JSON string from the parsed object.
 * If "compact" is false, an indented structure will be produced.
*/
void MinionParser::to_json(
    string &json_string, bool compact)
{
    if (top_level.size() == 0) {
        cerr << "JSON object: no content" << endl;
    }
    if (compact) {
        json_string = top_level.dump();
    } else {
        json_string = top_level.dump(2);
    }
}

MinionParser::MinionParser(
    const string_view source)
    : minion_string{source}
    , source_size{source.size()}
    , iter_i{0}
    , line_i{1}
{
    ch_pending = 0;
    top_level = json::object();
    get_map(top_level, 0);
}

json MinionParser::macro_replace(
    json item)
{
    if (item.is_string()) {
        string s{item};
        if (s.starts_with('&')) {
            try {
                return top_level.at(s);
            } catch (...) {
                error_message.append(
                    fmt::format("Undefined macro ({}) used in line {}\n", s, line_i));
            }
        }
    }
    return item;
}

/* Read the next input character.
 *
 * Parameter instring is true if a delimited string is being read.
 *
 * Returns the next input character, if it is valid.
 * If the source is exhausted return a null char.
 * If an illegal character is read an error report is added and a space
 * character is returned.
 */
Char MinionParser::read_ch(
    bool instring)
{
    if (ch_pending != 0) {
        Char ch = ch_pending;
        ch_pending = 0;
        if (ch == '\n') {
            ++line_i;
        }
        return ch;
    }
    if (iter_i < source_size) {
        Char ch = minion_string.at(iter_i++);
        //cout << "[CH: " << ch << "]" << endl;
        if (ch == '\n') {
            ++line_i; // increment line counter
            // These are not acceptable within strings:
            if (!instring) {
                // Don't return ' ', because unread_ch must be able to
                // distinguish the two, in order to adjust line_i
                return ch;
            }
            error_message.append(
                fmt::format("Unexpected newline in delimited string, line {}\n", line_i - 1));
        } else if (ch == '\r' || ch == '\t') {
            // These are acceptable in the source, but not within strings.
            if (!instring) {
                return ' ';
            }
        } else if (ch >= 32 && ch != 127) {
            return ch;
        }
        error_message.append(fmt::format("Illegal character ({:#x}) in line {}\n", ch, line_i - 1));
        return ' ';
    }
    return 0;
}

void MinionParser::unread_ch(
    Char ch)
{
    if (ch_pending != 0) {
        throw "Bug";
    }
    ch_pending = ch;
    if (ch == '\n') {
        --line_i;
    }
}

/* Read the next "item" from the input.
 *
 * Return a json value, which may be a string, an "array" (list) or
 * an "object" (map). If no value could be read (end of input) or there
 * was an error during reading, a null value will be returned.
 * If there was an error, an error message will be added for it.
 */
Char MinionParser::get_item(
    json &j)
{
    string udstring{};
    Char ch;
    separator = '*';
    while (true) {
        ch = read_ch(false);
        if (!udstring.empty()) {
            // An item has already been started
            while (true) {
                // Test for an item-terminating character
                if (ch == 0 || ch == ' ' || ch == '\n' || ch == '#' || ch == '"' || ch == '['
                    || ch == '{' || ch == ']' || ch == '}' || ch == ':') {
                    unread_ch(ch);
                    break;
                }
                udstring += ch;
                ch = read_ch(false);
            }
            j = json(udstring);
            //cout << "§2 " << udstring << endl;
            //cout << " :: " << j << endl;
            return ' ';
        }
        // Look for start of next item
        if (ch == 0) {
            j = json();
            return 0;
            separator = 0;
            break; // End of input => no further items
        }
        if (ch == ' ' || ch == '\n') {
            continue;
        }
        if (ch == u'#') {
            // Start comment
            ch = read_ch(false);
            if (ch == u'[') {
                // Extended comment: read to "]#"
                int comment_line = line_i;
                ch = read_ch(false);
                while (true) {
                    if (ch == u']') {
                        ch = read_ch(false);
                        if (ch == u'#') {
                            break;
                        }
                        continue;
                    }
                    if (ch == 0) {
                        error_message.append(
                            fmt::format("Unterminated comment ('#[ ...') in line {}\n", line_i - 1));
                        break;
                    }
                    // Comment loop ... read next character
                    ch = read_ch(false);
                }
                // End of extended comment
            } else {
                // "Normal" comment: read to end of line
                while (true) {
                    if (ch == '\n' || ch == 0) {
                        break;
                    }
                    ch = read_ch(false);
                }
            }
            continue;
        }
        // Delimited string
        if (ch == u'"') {
            get_string(j);
            return ' ';
        }
        // list
        if (ch == u'[') {
            j = json::array();
            get_list(j);
            if (j.is_null()) {
                // I don't think this is sensibly recoverable
                throw "Invalid list/array";
            }
            return ' ';
        }
        // map
        if (ch == u'{') {
            j = json::object();
            get_map(j, '}');
            if (j.is_null()) {
                // I don't think this is sensibly recoverable
                throw "Invalid map";
            }
            return ' ';
        }
        // further structural symbols
        if (ch == u']' || ch == u'}' || ch == u':') {
            j = json();
            return ch;
        }
        //cout << "§0 " << int(ch) << endl;
        udstring += ch;
    } // End of item-seeking loop
    cout << "BUG" << endl;
}

/* Read a delimited string (terminated by '"') from the input.
 *
 * It is entered after the initial '"' has been read, so the next character
 * will be the first of the string.
 *
 * Escapes, introduced by '\', are possible – see MINION specification.
 *
 * Return the string as a json value.
 * If an error was encountered, an error message will be added.
 */
void MinionParser::get_string(
    json &j)
{
    string dstring;
    Char ch;
    int start_line = line_i;
    while (true) {
        ch = read_ch(true);
        if (ch == 0) {
            error_message.append(
                fmt::format("Unterminated delimited string in line {}\n", line_i - 1));
            break;
        }
        if (ch == '"') {
            break; // end of string
        }
        if (ch == '\\') {
            // Deal with escapes:
            // "\n" ; "\t" ; "\/" ; "\'" ; "\{xxxx}" ; "\[ ... ]\"
            ch = read_ch(true);
            if (ch == u'n') {
                dstring += '\n';
                continue;
            }
            if (ch == u't') {
                dstring += '\n';
                continue;
            }
            if (ch == u'/') {
                dstring += '\\';
                continue;
            }
            if (ch == u'\'') {
                dstring += '"';
                continue;
            }
            if (ch == u'{') {
                // unicode character
                string ustr;
                while (true) {
                    // For the moment accept string characters.
                    ch = read_ch(true);
                    if (ch == '}') {
                        break;
                    }
                    if (ch == 0) {
                        error_message.append(
                            fmt::format("Unterminated unicode point in string in line {}\n",
                                        line_i - 1));
                        break;
                    }
                    if (ustr.size() > 5) {
                        ustr += '?'; // ensure the string is invalid ...
                        break;
                    }
                    ustr += ch;
                }
                if (!unicode_utf8(dstring, ustr)) {
                    error_message.append(
                        fmt::format("Invalid unicode point ({}) in string in line {}\n",
                                    ustr,
                                    line_i - 1));
                }
                continue;
            }
            if (ch == u'[') {
                // embedded comment: read to "]\"
                int comment_line = line_i;
                ch = read_ch(false);
                while (true) {
                    if (ch == u']') {
                        ch = read_ch(false);
                        if (ch == u'\\') {
                            break;
                        }
                        continue;
                    }
                    if (ch == 0) {
                        error_message.append(
                            fmt::format("Unterminated string comment ('\[ ...') in line {}\n",
                                        line_i - 1));
                        break;
                    }
                    // Comment loop ... read next character
                    ch = read_ch(false);
                }
                continue;
            }
        }
        // Add to string
        dstring += ch;
        // Loop ... read next character
    } // end of main loop
    j = json(dstring);
}

/* Read a "list" as a JSON array from the input.
 *
 * It is entered after the initial '[' has been read, so the search for the
 * next item will begin the following character.
 *
 * Return the list as a json value (array type).
 * If an error was encountered, an error message will be added.
 */
void MinionParser::get_list(
    json &j)
{
    int start_line = line_i;
    int item_line;
    json item;
    while (true) {
        item_line = line_i;
        Char sep = get_item(item);
        if (item.is_null()) {
            // No item found
            if (sep == ']') {
                return;
            }
            error_message.append(fmt::format(("Reading array starting in line {}."
                                              " In line {}: expected ']' or value\n"),
                                             start_line - 1,
                                             item_line - 1));
            j = json();
            return;
        }
        j.push_back(macro_replace(item));
    }
}

/* Read a "map" as a JSON object from the input.
 *
 * It is entered after the initial '{' has been read, so the search for the
 * next item will begin with the following character. The parameter is
 * '}', except for the top-level map, which has a null terminator.
 *
 * If the map was read successfully, it will be returned. If not, a null
 * item will be returned.
 */
void MinionParser::get_map(
    json &j, Char terminator)
{
    int start_line = line_i;
    int item_line;
    Char ch;
    string key;
    json item;
    while (true) {
        // Read key
        item_line = line_i;
        Char sep = get_item(item);
        //cout << "§1 " << ((sep == 0) ? 0 : sep) << endl;
        //cout << " :: " << item << endl;
        if (item.is_null()) {
            // No valid key found
            if (sep == terminator) {
                return;
            }
            error_message.append(fmt::format(("Reading map starting in line {}."
                                              " Item at line {}: expected key string\n"),
                                             start_line - 1,
                                             item_line - 1));
            j = json::object();
            return;
        }
        if (!item.is_string()) {
            //cout << item << endl;
            error_message.append(fmt::format(("Reading map starting in line {}."
                                              " Item at line {}: expected key string,\n"
                                              "Found: {}\n"),
                                             start_line - 1,
                                             item_line - 1,
                                             item.dump()));
            j = json::object();
            return;
        }
        key = item;
        // Expect ':'
        item_line = line_i;
        sep = get_item(item);
        if (item.is_null() && sep == ':') {
            //TODO: OK, read value
        } else {
            error_message.append(fmt::format(("Reading map starting in line {}."
                                              " Item at line {}: expected ':'\n"),
                                             start_line - 1,
                                             item_line - 1));
            j = json::object();
            return;
        }
        item_line = line_i;
        get_item(item);
        if (item.is_null()) {
            error_message.append(fmt::format(("Reading map starting in line {}."
                                              " Item at line {}: expected value"
                                              " for key \"{}\"\n"),
                                             start_line - 1,
                                             item_line - 1,
                                             key));
            j = json::object();
            return;
        }
        if (j.contains(key)) {
            error_message.append(fmt::format(("Reading map starting in line {}."
                                              " Key \"{}\" repeated at line {}\n"),
                                             start_line - 1,
                                             key,
                                             item_line - 1));
            j = json::object();
            return;
        }
        j[key] = macro_replace(item);
    } // end of loop
}

} // namespace Minion

/*
 * MINION: MINImal Object Notation, v.4
 * 
 * MINION is a simple data-transfer format taking some basic ideas
 * from JSON. It has features which make it suitable for easily readable
 * and writable configuration-files.
 * 
 * The only data type is the string. In addition there are containers:
 * lists and maps (associative arrays). Files must be encoded as utf-8.
 * Most of the ASCII control characters (0-31 and 127) are not allowed
 * and should be reported as errors. The permitted exceptions are '\n',
 * '\t', '\r' as layout/spacing characters, but not within strings.
 * There are other unicode characters which should probably be avoided,
 * but no checks are made as this is generally a difficult problem and it
 * is perhaps not clear where the line should be drawn.
 * 
 * The parsed data (not the source text!) is completely compatible with
 * parsed JSON, so a "to_json" method is provided for convenience.
 * 
 * A string may be enclosed in quotation marks (" ... "), but this is not
 * necessary if no "special" characters are included in the string.
 * 
 * Whitespaces are necessary as separators between items only when the
 * separation is not clear otherwise. They may, however, be added freely.
 * 
 * A plain comment continues to the end of the line. However, if the '#' is
 * directly followed by '[', the comment is terminated by "]#" and can
 * continue over line breaks.
 * 
 * The "special" characters are:
 *     "whitespace" characters (space, newline, etc.) – separators
 *     '#': start a comment
 *     ':': separates key from value in a map
 *     '{': start a map
 *     '}': end a map
 *     '[': start a list
 *     ']': end a list
 *     '"': string delimiter
 *     '\': string "escape" character (allowed in delimited string)
 *
 * map: { key:value key:value ... }
 *     "key" is a string.
 *     A "value" may be a string, a list or a map.
 *
 * list: [ value value ... ]
 *     A "value" may be a string, a list or a map.
 *
 * Certain characters are not directly possible in a string, they may be
 * included (only when the string is delimited by '"' characters) by means
 * of escape sequences:
 *     '"': "\'"
 *     '\': "\/"
 *     tab: "\t"
 *     newline: "\n"
 *     hexadecimal unicode character: "\{xxxx}" / "\{xxxxx}"
 *
 * In addition, it is possible to have an "embedded comment" in a
 * delimited string. This starts with "\[" and is ended by "]\".
 * As it may include newlines, this may be used to split a string
 * over several lines.
 *
 * The top level of a MINION text is a map – without the surrounding
 * braces ({ ... }).
 * 
 * There is also a very limited macro-like feature. Elements declared at the
 * top level which start with '&' may be referenced (which basically means
 * included) at any later point in a data structure by means of the macro
 * name, e.g.:
 *     &MACRO1: [A list of words]
 *       ...
 *     DEF1: { X: &MACRO1 }
 *
 * Note that keys beginning with '&' at lower levels will neither themselves
 * be replaced nor used to define replacement values.
*/
