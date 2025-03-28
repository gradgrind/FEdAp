cmake_minimum_required(VERSION 3.16)

project(courses_gui LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(FLTK 1.4 CONFIG REQUIRED)

add_executable(courses_gui main.cpp
    courses/courses_gui.cpp courses/courses_gui.h
    courses/editform.h courses/editform.cpp

)

include(GNUInstallDirs)
install(TARGETS courses_gui
    LIBRARY DESTINATION ${CMAKE_INSTALL_LIBDIR}
    RUNTIME DESTINATION ${CMAKE_INSTALL_BINDIR}
)

target_link_libraries(courses_gui PRIVATE fltk::fltk)
