cmake_minimum_required(VERSION 3.16)

project(guilibtest LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(FLTK 1.4 CONFIG REQUIRED)

add_executable(guilibtest main.cpp

    layout.cpp
    widgetdata.cpp
    dispatcher.cpp
    widget_methods.cpp
    widget_methods.h
    widgets.cpp

    minion.h
    minion.cpp
    iofile.cpp
    iofile.h
    layout.h
    widgets.h
    functions.h
    functions.cpp
    rowtable.h rowtable.cpp
    editform.h
    editform.cpp





)

include(GNUInstallDirs)
install(TARGETS guilibtest
    LIBRARY DESTINATION ${CMAKE_INSTALL_LIBDIR}
    RUNTIME DESTINATION ${CMAKE_INSTALL_BINDIR}
)

target_link_libraries(guilibtest PRIVATE fltk::fltk)
target_link_libraries(guilibtest PRIVATE fmt)
target_include_directories(guilibtest PRIVATE ../_lib)
