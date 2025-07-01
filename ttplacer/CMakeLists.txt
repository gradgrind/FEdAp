cmake_minimum_required(VERSION 3.16)

project(ttplacer
    VERSION 0.0.0
    LANGUAGES CXX)

if(NOT CMAKE_BUILD_TYPE)
  set(CMAKE_BUILD_TYPE Release)
endif()

set(CMAKE_CXX_FLAGS "-Wall -Wextra")
set(CMAKE_CXX_FLAGS_DEBUG "-g")
set(CMAKE_CXX_FLAGS_RELEASE "-O2")

add_executable(ttplacer
    main.cpp
    ../mugui/minion.cpp
    ../mugui/iofile.cpp
    baseconstraint.h baseconstraint.cpp
)
    
target_include_directories(ttplacer PUBLIC ../mugui)

set_property(TARGET ttplacer PROPERTY CXX_STANDARD 20)
