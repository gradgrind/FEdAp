#include "minion.h"
#include "iofile.h"
#include <chrono>
#include <fmt/format.h>
#include <iostream>
using namespace std;
using namespace std::chrono;

// *** Reading to custom object. This version is (still) using only
// string as the basic data type.

namespace minion {

MinionMap read_minion(
    string_view minion_string)
{
    minion::Minion mp(minion_string);
    if (mp.error_message.empty()) {
        return mp.top_level;
    }
    throw minion::MinionException(mp.error_message);
}

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

void dump_string(
    string &valstr, const string &s)
{
    valstr += '"';
    for (const Char &ch : s) {
        if (ch >= 32) {
            if (ch == '"') {
                valstr += "\\'";
            } else if (ch == '\\') {
                valstr += "\\/";
            } else if (ch != 127) {
                valstr += ch;
            } else {
                valstr += "\\{007F}";
            }
        } else if (ch == '\n') {
            valstr += "\\n";
        } else if (ch == '\t') {
            valstr += "\\t";
        } else {
            valstr += fmt::format("\\{{{:#04x}}}", ch);
        }
    }
    valstr += '"';
}

// Dump the value as MINION string.
// If level < 0, add no formatting/padding, otherwise format with
// indentation.
const int indent_depth = 2;
void dump(
    string &valstr, MinionValue item, int level = -1)
{
    if (holds_alternative<string>(item)) {
        auto s{get<string>(item)};
        dump_string(valstr, s);
    } else if (holds_alternative<MinionMap>(item)) {
        auto m{get<MinionMap>(item)};
        if (m.empty()) {
            valstr += "{}";
            return;
        }
        valstr += '{';
        valstr += dump_map_items(m, (level < 0) ? -1 : level + 1);
        if (level >= 0) {
            valstr += '\n' + string(indent_depth * level, ' ');
        }
        valstr += '}';
    } else if (holds_alternative<MinionList>(item)) {
        auto m{get<MinionList>(item)};
        if (m.empty()) {
            valstr += "[]";
            return;
        }
        valstr += '[';
        valstr += dump_list_items(m, (level < 0) ? -1 : level + 1);
        if (level >= 0) {
            valstr += '\n' + string(indent_depth * level, ' ');
        }
        valstr += ']';
    } else {
        valstr += "NO_VALUE";
    }
}

string dump_list_items(
    const MinionList m, int level)
{
    string padding;
    if (level >= 0) {
        padding += '\n' + string(indent_depth * level, ' ');
    }
    string valstr;
    for (const auto &item : m) {
        valstr += padding;
        dump(valstr, item, level);
    }
    return valstr;
}

string dump_map_items(
    const MinionMap m, int level)
{
    string padding;
    string keysep{':'};
    if (level >= 0) {
        padding += '\n' + string(indent_depth * level, ' ');
        keysep += ' ';
    }
    string valstr;
    for (const auto &item : m) {
        valstr += padding;
        dump_string(valstr, item.key);
        valstr += keysep;
        dump(valstr, item.value, level);
    }
    return valstr;
}

// This is, of course, rather inefficient for maps which are not very short.
// Making a map out of this would make the MinionValues a bit larger and lose
// the ordering, unless a more complicated map structure is used.
MinionValue MinionMap::get(
    string_view key)
{
    for (const auto &mmp : *this) {
        if (mmp.key == key) return mmp.value;
    }
    return MinionValue{};
}

/*
// Generate a JSON string from the parsed object.
// If "compact" is false, an indented structure will be produced.
void Minion::to_json(
    string &json_string, bool compact)
{
    if (top_level.size() == 0) {
        cerr << "JSON object: no content" << endl;
    }
    if (compact) {
        json_string = top_level.jdump();
    } else {
        json_string = top_level.jdump(2);
    }
}
*/

Minion::Minion(
    const string_view source)
    : minion_string{source}
    , source_size{source.size()}
    , iter_i{0}
    , line_i{1}
{
    ch_pending = 0;
    top_level = MinionMap();
    get_map(top_level, 0);
}

MinionValue Minion::macro_replace(
    MinionValue item)
{
    if (holds_alternative<string>(item)) {
        string s{get<string>(item)};
        if (s.starts_with('&')) {
            try {
                return macros.at(s);
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
Char Minion::read_ch(
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

void Minion::unread_ch(
    Char ch)
{
    if (ch_pending != 0) {
        throw minion::MinionException("Bug in unread_ch");
    }
    ch_pending = ch;
    if (ch == '\n') {
        --line_i;
    }
}

/* Read the next "item" from the input.
 *
 * Return a MinionValue, which may be a string, an "array" (list) or an
 * "object" (map). If no value could be read (end of input) or there was an
 * error during reading, a null value will be returned (m.index() == 0).
 * If there was an error, an error message will be added for it.
 */
Char Minion::get_item(
    MinionValue &m)
{
    string udstring{};
    Char ch;
    while (true) {
        ch = read_ch(false);
        if (!udstring.empty()) {
            // An undelimited string item has already been started
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
            m = MinionValue{udstring};
            //cout << "§2 " << udstring << endl;
            //cout << " :: " << j << endl;
            return ' ';
        }
        // Look for start of next item
        if (ch == 0) {
            m = MinionValue{}; // end of input, no next item
            return 0;
        }
        if (ch == ' ' || ch == '\n') {
            continue; // continue seeking start of item
        }
        if (ch == u'#') {
            // Start comment
            ch = read_ch(false);
            if (ch == u'[') {
                // Extended comment: read to "]#"
                //int comment_line = line_i;
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
            continue; // continue seeking item
        }
        // Delimited string
        if (ch == u'"') {
            get_string(m);
            return ' ';
        }
        // list
        if (ch == u'[') {
            m = MinionList{};
            get_list(m);
            if (m.index() == 0) {
                // I don't think this is sensibly recoverable
                throw minion::MinionException(error_message);
                //throw "Invalid list/array";
            }
            return ' ';
        }
        // map
        if (ch == u'{') {
            MinionMap mm;
            if (!get_map(mm, '}')) {
                // I don't think this is sensibly recoverable
                throw minion::MinionException(error_message);
                //throw "Invalid map";
            }
            m = MinionValue{mm};
            return ' ';
        }
        // further structural symbols
        if (ch == u']' || ch == u'}' || ch == u':') {
            m = MinionValue{};
            return ch; // no item, but significant terminator
        }
        //cout << "§0 " << int(ch) << endl;
        udstring += ch; // start undelimited string
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
 * Return the string as a MinionValue.
 * If an error was encountered, an error message will be added.
 */
void Minion::get_string(
    MinionValue &m)
{
    string dstring;
    Char ch;
    //int start_line = line_i;
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
                        ustr += '?'; // ensure the unicode hex is invalid ...
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
    m = MinionValue{dstring};
}

/* Read a "list" as a MinionValue (MinionList) from the input.
 *
 * It is entered after the initial '[' has been read, so the search for the
 * next item will begin the following character.
 *
 * Return the list as a json value (array type).
 * If an error was encountered, an error message will be added.
 */
void Minion::get_list(
    MinionValue &m)
{
    int start_line = line_i;
    int item_line;
    MinionValue item;
    MinionList l;
    while (true) {
        item_line = line_i;
        Char sep = get_item(item);
        if (item.index() == 0) {
            // No item found
            if (sep == ']') {
                m.emplace<MinionList>(l);
                return;
            }
            error_message.append(fmt::format(("Reading array starting in line {}."
                                              " In line {}: expected ']' or value\n"),
                                             start_line,
                                             item_line));
            m.emplace<0>();
            return;
        }
        l.emplace_back(macro_replace(item));
    }
}

bool Minion::get_map(
    MinionMap &m, Char terminator)
{
    int start_line = line_i;
    int item_line;
    Char ch;
    string key;
    MinionValue item;
    while (true) {
        // Read key
        item_line = line_i;
        Char sep = get_item(item);
        //cout << "§1 " << ((sep == 0) ? 0 : sep) << endl;
        //cout << " :: " << item << endl;
        if (item.index() == 0) {
            // No valid key found
            if (sep == terminator) {
                return true;
            }
            error_message.append(fmt::format(("Reading map starting in line {}."
                                              " Item at line {}: expected key string\n"),
                                             start_line,
                                             item_line));
            return false;
        }
        if (!holds_alternative<string>(item)) {
            //cout << item << endl;
            string itemstr;
            dump(itemstr, item, -1);
            error_message.append(fmt::format(("Reading map starting in line {}."
                                              " Item at line {}: expected key string,\n"
                                              "Found: {}\n"),
                                             start_line,
                                             item_line,
                                             itemstr));
            return false;
        }
        key = get<string>(item);
        for (const auto &mmp : m) {
            if (mmp.key == key) {
                error_message.append(fmt::format(("Reading map starting in line {}."
                                                  " Key \"{}\" repeated at line {}\n"),
                                                 start_line,
                                                 key,
                                                 item_line));
                return false;
            }
        }
        // Expect ':'
        item_line = line_i;
        sep = get_item(item);
        if (item.index() != 0 || sep != ':') {
            error_message.append(fmt::format(("Reading map starting in line {}."
                                              " Item at line {}: expected ':'\n"),
                                             start_line,
                                             item_line));
            return false;
        }
        item_line = line_i;
        get_item(item);
        if (item.index() == 0) {
            error_message.append(fmt::format(("Reading map starting in line {}."
                                              " Item at line {}: expected value"
                                              " for key \"{}\"\n"),
                                             start_line - 1,
                                             item_line - 1,
                                             key));
            return false;
        }
        auto val = macro_replace(item);
        // Check for top-level &-keys
        if (terminator == 0 && key.starts_with('&'))
            macros.emplace(key, val);
        else
            m.emplace_back(MinionMapPair{key, val});
    } // end of loop
}

bool MinionMap::get_int(
    string_view key, int &value)
{
    auto s = get(key);
    if (holds_alternative<string>(s)) {
        // Read as integer
        value = stoi(get<string>(s));
        return true;
    }
    return false;
}

} // namespace minion

void test_minion(
    const string &filepath)
{
    string idata;
    if (!readfile(idata, filepath)) {
        cerr << "Error opening file: " << filepath << endl;
        return;
    }
    cout << "FILE: " << filepath << endl;

    // Use auto keyword to avoid typing long
    // type definitions to get the timepoint
    // at this instant use function now()
    auto start = high_resolution_clock::now();

    minion::Minion mp(idata);

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
        //mp.to_json(odata, false);
        //cout << " >>> " << dump_map_items(mp.top_level, 0) << endl;

        /*
        auto p = filepath.rfind(".");
        string f;
        if (p == string::npos) {
            f = filepath;
        } else {
            f = filepath.substr(0, p);
        }
        string f1 = f + ".json";
        if (!writefile(odata, f) {
            cerr << "Error opening file: " << f << endl;
            return;
        }
        cout << " --> " << f << endl;
        */
    } else {
        cout << "ERROR:\n" << mp.error_message << endl;
        return;
    }
}

void testminion()
{
    test_minion("_data/test0.minion");
    test_minion("_data/test1.minion");
    test_minion("_data/test2.minion");

    string idata;
    string fp{"_data/test2e.minion"};
    if (readfile(idata, fp)) {
        cout << "Reading " << fp << endl;
        try {
            minion::read_minion(idata);
        } catch (minion::MinionException &e) {
            cerr << e.what() << endl;
        }
    } else {
        cerr << "Error opening file: " << fp << endl;
    }

    /*
    // Now look for leakages ...
    string xdata;
    string fpx{"_data/test2.minion"};
    minion::MinionMap mmap;
    minion::MinionValue mval;
    if (readfile(xdata, fpx)) {
        for (int i = 0; i < 40000; ++i) {
            mmap = minion::read_minion(xdata);
            mval = mmap.get("EXTRA_FIELD_WIDTHS");
        }
    } else {
        cerr << "Error opening file: " << fpx << endl;
    }
    */
}

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
 * These macro keys are not included in the resulting data map.
 * 
 * Note that keys beginning with '&' at lower levels will neither themselves
 * be replaced nor used to define replacement values.
*/
