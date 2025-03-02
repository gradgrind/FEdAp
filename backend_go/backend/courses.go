package backend

func init() {
	commandMap["GET_COURSES"] = getCourses
}

func getCourses(cmd map[string]any, outmap map[string]any) string {
	outmap["Courses"] = DB.Courses
	outmap["SuperCourses"] = DB.SuperCourses
	outmap["SubCourses"] = DB.SubCourses
	return "OK"
}

//-- {"DO":"LOAD_W365_JSON"}
//-- {"DO":"GET_COURSES"}
