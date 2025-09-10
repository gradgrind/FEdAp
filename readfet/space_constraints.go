package readfet

import "encoding/xml"

type spaceConstraints struct {
	XMLName                          xml.Name `xml:"Space_Constraints_List"`
	ConstraintBasicCompulsorySpace   basicSpaceConstraint
	ConstraintActivityPreferredRooms []roomChoice
	ConstraintActivityPreferredRoom  []placedRoom
	ConstraintRoomNotAvailableTimes  []roomNotAvailable
}

type basicSpaceConstraint struct {
	XMLName           xml.Name `xml:"ConstraintBasicCompulsorySpace"`
	Weight_Percentage int
	Active            bool
}

type placedRoom struct {
	XMLName              xml.Name `xml:"ConstraintActivityPreferredRoom"`
	Weight_Percentage    int
	Activity_Id          int
	Room                 string
	Number_of_Real_Rooms int      `xml:",omitempty"`
	Real_Room            []string `xml:",omitempty"`
	Permanently_Locked   bool     // false
	Active               bool     // true
}

type roomChoice struct {
	XMLName                   xml.Name `xml:"ConstraintActivityPreferredRooms"`
	Weight_Percentage         int
	Activity_Id               int
	Number_of_Preferred_Rooms int
	Preferred_Room            []string
	Active                    bool // true
}
