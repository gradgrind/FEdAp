/* This is a script to generate a multiple page document where each class,
 * teacher or room has a page of its own for its weekly timetable.
 *
 * The basic idea is to use a Typst-table to manages the headers, lines and
 * background colouring. The tiles of the timetable activities are overlaid
 * on this using the table cell boundary coordinates for orientation.
 *
 * A space (block) is left at the top of each page for a page title. This
 * is in the space left free by the TITLE_HEIGHT value.
 * The rest of the page will be used for the table, adjusting the cell size
 * to fit.
 *
 * Two variants of table are supported:
 *  - plain (default): A simple table is used with each lesson period having
 *                     the same length.
 *  - with breaks:     The height of the lesson periods is proportional to
 *                     their times in minutes. Also the breaks between lessons
 *                     will be shown, proportional to their times. This can
 *                     only work if the lesson periods (hours) are supplied
 *                     with correctly formatted start and end times (hh:mm,
 *                     24-hour clock).
 *                     To select this variant, the option Typst.WithBreaks
 *                     must be true.
 * If lesson period times are supplied, the parameter Typst.WithTimes must be
 * true (default: false) for them to be shown in the period headers.
 */
 
//#import "base.typ"

// Schriftart:
#set text(font: ("Nunito","DejaVu Sans"))
// Diese muss auf dem Waldorf 365 PDF-Server installiert sein. Wenden Sie sich an den Support.

#let PAGE_HEIGHT = 210mm
#let PAGE_WIDTH = 297mm
#let TITLE_HEIGHT = 10mm
#let PAGE_BORDER = (top:10mm, bottom: 10mm, left: 15mm, right: 15mm)
#let H_HEADER_HEIGHT = 13mm
#let V_HEADER_WIDTH = 30mm
#let HEADING_TO_PLAN_DISTANCE = 4mm
#let PLAN_TO_FOOTER_DISTANCE = 3mm

#let CELL_BORDER =0.5pt
#let CELL_INSETS =3pt
#let BIG_SIZE = 18pt
#let NORMAL_SIZE = 16pt
#let PLAIN_SIZE = 12pt
#let SMALL_SIZE = 10pt

#let FRAME_COLOUR = "#909090"
#let HEADER_COLOUR = "#f0f0f0"
#let EMPTY_COLOUR = "#ffffff"
#let BREAK_COLOUR = "#ffffff"

#let JOINSTR = ","
#let CLASS_GROUP_JOIN = "."

// Field placement fallbacks
#let boxText = (
    Class: (
        C: "SUBJECT",
        TL: "TEACHER",
        TR: "GROUP",
        //BL: "",
        BR: "ROOM",
    ),
    Teacher: (
        C: "GROUP",
        TL: "SUBJECT",
        TR: "TEACHER",
        //BL: "",
        BR: "ROOM",
    ),
    Room: (
        C: "GROUP",
        TL: "SUBJECT",
        //TR: "",
        //BL: "",
        BR: "TEACHER",
    ),
)

// Überschriften
// Folgende Zeichen werden ersetzt
// %0 = Kürzel des Lehrer / Raums ...
// %1 = Nachname des Lehrers, Name der Klasse oder des Raums
// %2 = Vorname des Lehrers
#let pageHeadings = (
    Class: "Klasse %0",
    Teacher: "%2 %1",
    Room: "Raumplan %1",
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// Load and parse data
#let xdata = json(sys.inputs.ifile)
#let typstMap = xdata.at("Typst", default: (:))
#let pagesFromTypstMap = typstMap.at("Pages", default: "")

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

#set page(height: PAGE_HEIGHT, width: PAGE_WIDTH,
  margin: PAGE_BORDER,
)

#let PLAN_AREA_WIDTH = (PAGE_WIDTH - PAGE_BORDER.left - PAGE_BORDER.right)

#let DAYS = ()
#for ddata in xdata.Info.Days {
    //TODO: Which field to use
    DAYS.push(ddata.Name)
}

#let HOURS = ()
#let TIMES = ()
#let WITHTIMES = typstMap.at("WithTimes", default: false)
#let WITHBREAKS = typstMap.at("WithBreaks", default: false)

#for hdata in xdata.Info.Hours {
    //TODO: Which field to use
  let hour = hdata.Short
  let time1 = hdata.at("Start")
  let time2 = hdata.at("End")
  if WITHTIMES {
    HOURS.push(hour + "|" + time1 + " – " + time2)
  } else {
    HOURS.push(hour)
  }
  if WITHBREAKS {
    // The period start and end times must be present and correct.
    let (h1, m1) = time1.split(":")
    let (h2, m2) = time2.split(":")
    // Convert the times to minutes only.
    TIMES.push((int(h1) * 60 + int(m1), int(h2) * 60 + int(m2)))
  }
}

#let classList = xdata.Info.ClassNames
#let teacherList = xdata.Info.TeacherNames
#let roomList = xdata.Info.RoomNames
#let subjectList = xdata.Info.SubjectNames

#let listListToMap(listList) = {
    let map = (:)
    for list in listList {
        map.insert(list.at(0), list)
    }
    map
}

#let nameMap = (
    "Class": listListToMap(classList),
    "Teacher": listListToMap(teacherList),
    "Room": listListToMap(roomList),
    "Subject": listListToMap(subjectList),
)

#let legendItems(ilist) = {
    let items = ()
    for item in ilist {
        if item.len() == 3 {
            // Assume teacher
            items.push(box(item.at(0)+" = "+item.at(2)+" "+item.at(1)))
        } else {
            items.push(box(item.at(0)+" = "+item.at(1)))
        }
    }
    items.join(", ")
}

#let legendItemsLookup(elementType, shortList) = {
    let nmap = nameMap.at(elementType)
    let elist = shortList.map(it => nmap.at(it))
    legendItems(elist)
}

// Type of table ("Class", "Teacher" or "Room")
#let tableType = xdata.TableType

// Determine the field placements in the tiles
#let fieldPlacements = typstMap.at("FieldPlacement", default: (:))
#if fieldPlacements.len() == 0 {
    // fallback
    fieldPlacements = boxText.at(tableType, default: (:))
}

#let lastChange0 = typstMap.at("LastChange", default: "")
#let legend0 = typstMap.at("Legend", default: (:))

#let buildLegend(activities, legend) = {
    // Collect elements
    let slist = ()
    let tlist = ()
    let rlist = ()
    for a in activities {
        slist.push(a.at("Subject"))
        tlist += a.at("Teachers", default: ())
        rlist += a.at("Rooms", default: ())
    }
    // Add text items
    if legend.at("Subjects", default: false) {
        [*Fächer:* #legendItemsLookup("Subject", slist.dedup()) \ ]
    }
    if legend.at("Teachers", default: false) {
        [*Lehrkräfte:* #legendItemsLookup("Teacher", tlist.dedup()) \ ]
    }
    if legend.at("Rooms", default: false) {
        [*Räume:* #legendItemsLookup("Room", rlist.dedup()) \ ]
    }
}

// CREATE THE PAGES

#let npage = 0
#for p in xdata.Pages {
    if npage != 0 {
        pagebreak()
    }

	let legend = legend0 + pagesFromTypstMap.at(
        npage, default: (:)).at("Legend", default: (:))
    npage += 1
    
    let legendRemark = legend.at("Remark", default: "")
    let footnotesLegend = legend.at("Footnotes", default: ())

    context { 
    
        let legendBlock = text(size:9pt)[
            #if legendRemark != "" [*Hinweis:* #legendRemark \ ]
            #if footnotesLegend.len() != 0 [*Anmerkungen:* #legendItems(footnotesLegend) \ ]
            #buildLegend(p.Activities, legend)
        ]
   
        // All code using the footerHeight needs to be wrapped inside a context 
        // as the context is needed for the measure method.

        let footerHeight = measure(width : PLAN_AREA_WIDTH, legendBlock).height
    
        let PLAN_AREA_HEIGHT = (PAGE_HEIGHT - PAGE_BORDER.top
            - PAGE_BORDER.bottom - TITLE_HEIGHT - footerHeight
            - HEADING_TO_PLAN_DISTANCE - PLAN_TO_FOOTER_DISTANCE - 4mm)

        let vfactor = if WITHBREAKS {
            // Here it is a factor with which to multiply the minutes
            let tdelta = TIMES.at(-1).at(1) - TIMES.at(0).at(0)
            (PLAN_AREA_HEIGHT - H_HEADER_HEIGHT) / tdelta
        } else {
            // Here it is just the height of a period box
            (PLAN_AREA_HEIGHT - H_HEADER_HEIGHT) / HOURS.len()
        }

        // Build the row structure
        let table_content = ([],) + DAYS
        let isbreak = (false,)
        let hlines = ()
        let trows = (H_HEADER_HEIGHT,)
        let t = 0mm
        let m0 = -1
        let i = 0
        for h in HOURS {
            if WITHBREAKS {
            let (m1, m2) = TIMES.at(i)
            if m0 < 0 {
                m0 = m1
            }
            let t1 = (m1 - m0) * vfactor
            if t > 0mm {
                trows.push(t1 - t)
                table_content += ("",) + ([],) * DAYS.len()
                isbreak.push(true)
            }
            t = (m2 - m0) * vfactor
            hlines.push((t1 + H_HEADER_HEIGHT, t + H_HEADER_HEIGHT))
            trows.push(t - t1)
            } else {
                let period_height = vfactor
                let t0 = t
                t += period_height
                hlines.push((t0 + H_HEADER_HEIGHT, t + H_HEADER_HEIGHT))
                trows.push(period_height)
            }
            table_content += (h,) + ([],) * DAYS.len()
            isbreak.push(false)
            i += 1
        }

        // Build the vertical lines

        let vlines = (V_HEADER_WIDTH,)
        let colwidth = (PLAN_AREA_WIDTH - V_HEADER_WIDTH) / DAYS.len()
        let d0 = V_HEADER_WIDTH
        //COLWIDTH #colwidth
        for d in DAYS {
            d0 += colwidth
            vlines.push(d0)
        }
        let tcolumns = (V_HEADER_WIDTH,) + (colwidth,)*DAYS.len()

        show table.cell: it => {
            if it.y == 0 {
            set text(size: NORMAL_SIZE) //, weight: "bold")
            align(center + horizon, it.body.at("text", default: ""))
            } else if it.x == 0 {
            let txt = it.body.at("text", default: "")
            let t1t2 = txt.split("|")
            let tt = text(size: SMALL_SIZE, weight: "bold", t1t2.at(0))
            if t1t2.len() > 1 {
                tt += [\ ] + text(size: SMALL_SIZE, t1t2.at(1))
            }
            align(center + horizon, tt)
            } else {
            it
            }
        }

        // On text lines (top or bottom of cell) with two text items:
        // If one is smaller than 25% of the space, leave this and shrink the
        // other to 90% of the reamining space. Otherwise shrink both.

        let scale2inline(vpos, width, text1, text2) = {
            let t1 = text(size: SMALL_SIZE, text1)
            let t2 = text(size: SMALL_SIZE, text2)
            let w4 = width / 4
            context {
                let s1 = measure(t1)
                let s2 = measure(t2)
                if (s1.width + s2.width) > width * 0.9 {
                    if s1.width < w4 {
                        // shrink only text2
                        let w2 = width - s1.width
                        let scl = (w2 * 0.9) / s2.width
                        place(vpos + left, t1)
                        place(vpos + right, text(size: scl * SMALL_SIZE, text2))
                    } else if s2.width < w4 {
                        // shrink only text1
                        let w2 = width - s2.width
                        let scl = (w2 * 0.9) / s1.width
                        place(vpos + left, text(size: scl * SMALL_SIZE, text1))
                        place(vpos + right, t2)
                    } else {
                        // shrink both
                        let scl = (width * 0.9) / (s1.width + s2.width)
                        place(vpos + left, text(size: scl * SMALL_SIZE, text1))
                        place(vpos + right, text(size: scl * SMALL_SIZE, text2))
                    }
                } else {
                    place(vpos + left, t1)
                    place(vpos + right, t2)
                }
            }
        }

        let scaleinline(vpos, width, textc) = {
            let t = text(size: NORMAL_SIZE, weight: "bold", textc)
            context {
                let s = measure(t)
                if s.width > width * 0.9 {
                    let scl = (width * 0.9 / s.width)
                    let ts = text(size: scl * NORMAL_SIZE, weight: "bold", textc)
                    place(vpos + center, ts)
                } else {
                    place(vpos + center, t)
                }
            }
        }

        let cell_inset = CELL_BORDER
        let cell_width = colwidth 

        let ttcell(
            Day: 0,
            Hour: 0,
            Duration: 1,
            Offset: 0,
            Fraction: 1,
            Total: 1,
            Subject: "",
            Groups: (),
            Teachers: (),
            Rooms: (),
            Background: "",
            Footnote: "",
        ) = {
            // Prepare texts
            let clist = () // Class part only
            let glist = ()
            let cglist = ()
            for (c, g) in Groups {
                clist.push(c)
                let cg = c
                if g != "" {
                    cg += CLASS_GROUP_JOIN + g
                }
                if p.Short == c {
                    if g != "" {
                        glist.push(g)
                    }
                } else {
                    glist.push(cg)
                }
                cglist.push(cg)
            }
            clist = clist.dedup() // Remove duplicates
            let texts = (
                SUBJECT: Subject,
                CLASS: clist.join(JOINSTR),
                TEACHER: Teachers.join(JOINSTR),
                ROOM: Rooms.join(JOINSTR),
                GROUP: glist.join(JOINSTR),
                CLASS_WITH_GROUP: cglist.join(JOINSTR),
            )
            
            let footnoteIndicator = if Footnote == "" {""} else {" "+Footnote}
            
            let centre = texts.at(fieldPlacements.at("C", default: ""), default: "") 
            let tl = texts.at(fieldPlacements.at("TL", default: ""), default: "")
            let tr = texts.at(fieldPlacements.at("TR", default: ""), default: "") + footnoteIndicator
            let bl = texts.at(fieldPlacements.at("BL", default: ""), default: "")
            let br = texts.at(fieldPlacements.at("BR", default: ""), default: "")

            let cellBorderColour = FRAME_COLOUR
            if Background == "" {
                Background = "#FFFFFF"
                cellBorderColour = "#000000"
            }
            let bg = rgb(Background)
            
            // Used by fgWCAG2
            let rgblumin(c) = {
                c = c / 100%
                if c <= 0.04045 {
                    c/12.92
                } else {
                    calc.pow((c+0.055)/1.055, 2.4)
                }
            }

            // Decide on black or white for text, based on background colour (WCAG2).
            let fgWCAG2(colour) = {
                let (r,g,b,a) = rgb(colour).components()
                let l = 0.2126 * rgblumin(r) + 0.7152 * rgblumin(g) + 0.0722 * rgblumin(b)
                if l > 0.179 { black } else { white }
            }

            let fgW365(colour) = {
                let (r,g,b,a) = rgb(colour).components()
                let l = 0.2126 * calc.pow(r / 100%, 2.2)
                    + 0.7152 * calc.pow(g / 100%, 2.2)
                    + 0.0722 * calc.pow(b/ 100%, 2.2)
                if l > calc.pow(0.7, 2.2) { black } else { white }
            }
        
            // Get text colour
            // 1) This converts the background to grey-scale and uses a threshold:
            let bw = luma(bg)
            //let textcolour = if bw.components().at(0) < 55% { white } else { black }
            //let textcolour = fgW365(bg)
            let textcolour = if bw.components().at(0) < 75% { white } else { black }
            set text(textcolour)

            // Determine grid lines
            let (y0, y1) = hlines.at(Hour)
            let x0 = vlines.at(Day)
            if Duration > 1 {
                y1 = hlines.at(Hour + Duration - 1).at(1)
            }
            let wfrac = cell_width * Fraction / Total
            let xshift = cell_width * Offset / Total
            // Shrink excessively large components.
            let b = box(
                fill: bg,
                stroke: (paint: rgb(cellBorderColour), thickness:CELL_BORDER),
                inset: CELL_INSETS,
                height: y1 - y0 ,
                width: wfrac,
            )[
                #scale2inline(top, wfrac, tl, tr)
                #scaleinline(horizon, wfrac, centre)
                #scale2inline(bottom, wfrac, bl, br)
            ]
            place(top + left,
                dx: x0  + xshift,
                dy: y0 ,
                b
            )
        }

        let tbody = table(
            columns: tcolumns,
            rows: trows,
            gutter: 0pt,
            stroke: (paint: rgb(FRAME_COLOUR), thickness:CELL_BORDER),
            inset: 0pt,
            fill: (x, y) =>
                if y != 0 {
                    if isbreak.at(y) {
                        rgb(BREAK_COLOUR)
                    } else if x != 0 {
                        rgb(EMPTY_COLOUR)
                    }
                },
            ..table_content
        )

        show heading: it => text(weight: "bold", size: BIG_SIZE,
            bottom-edge: "descender",
            pad(left: 0mm, it))

        let phead = pageHeadings.at(tableType, default: "")
        let nameList = nameMap.at(tableType).at(p.Short)
        let title = phead.replace(regex("%([0-2])"),
            m => nameList.at(int(m.captures.first())))

        let lastChange = pagesFromTypstMap.at(
            npage - 1,default: (:)).at("LastChange", default: "")
        if lastChange == "" {
            lastChange = lastChange0
        }
        let st = typstMap.at("Subtitle", default: "")
        if st == "" {
            st = lastChange
        } else if lastChange != "" {
            st += " · " + lastChange
        }

        // HEADER
        block(height: TITLE_HEIGHT, above: 0mm,
            below: HEADING_TO_PLAN_DISTANCE, inset: 0mm)[
            #place(top)[= #title]
            #place(top)[#h(1fr) #text(17pt, xdata.Info.Institution)]
            #place(bottom)[#st]
        ]

        // BODY
        
        box([
            #tbody
            #for a in p.Activities {
                ttcell(..a)
            }
        ])

        // FOOTER
        
        legendBlock
        
    }
}
