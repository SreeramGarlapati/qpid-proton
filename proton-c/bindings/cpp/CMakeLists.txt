#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

include(${CMAKE_SOURCE_DIR}/cpp.cmake) # Compiler checks

include_directories(
  "${CMAKE_SOURCE_DIR}/proton-c/include"
  "${CMAKE_CURRENT_SOURCE_DIR}/include"
  "${CMAKE_CURRENT_SOURCE_DIR}/src")

set(qpid-proton-cpp-source
  src/acceptor.cpp
  src/blocking_connection.cpp
  src/blocking_connection_impl.cpp
  src/blocking_fetcher.cpp
  src/blocking_link.cpp
  src/blocking_receiver.cpp
  src/blocking_sender.cpp
  src/connection.cpp
  src/connection_options.cpp
  src/connector.cpp
  src/container.cpp
  src/container_impl.cpp
  src/contexts.cpp
  src/data.cpp
  src/decoder.cpp
  src/delivery.cpp
  src/duration.cpp
  src/encoder.cpp
  src/endpoint.cpp
  src/engine.cpp
  src/error.cpp
  src/event.cpp
  src/handler.cpp
  src/link.cpp
  src/message.cpp
  src/messaging_adapter.cpp
  src/messaging_event.cpp
  src/messaging_handler.cpp
  src/object.cpp
  src/proton_bits.cpp
  src/proton_event.cpp
  src/proton_handler.cpp
  src/reactor.cpp
  src/receiver.cpp
  src/reconnect_timer.cpp
  src/request_response.cpp
  src/sasl.cpp
  src/sender.cpp
  src/session.cpp
  src/ssl.cpp
  src/ssl_domain.cpp
  src/task.cpp
  src/terminus.cpp
  src/transport.cpp
  src/types.cpp
  src/url.cpp
  src/uuid.cpp
  src/value.cpp
  )


set_source_files_properties (
  ${qpid-proton-cpp-source}
  PROPERTIES
  COMPILE_FLAGS "${CXX_WARNING_FLAGS}"
  )

add_library(qpid-proton-cpp SHARED ${qpid-proton-cpp-source})

target_link_libraries (qpid-proton-cpp ${PLATFORM_LIBS} qpid-proton)

set_target_properties (
  qpid-proton-cpp
  PROPERTIES
  LINKER_LANGUAGE CXX
  VERSION   "${PN_LIB_SOMAJOR}.${PN_LIB_SOMINOR}"
  SOVERSION "${PN_LIB_SOMAJOR}"
  LINK_FLAGS "${CATCH_UNDEFINED}"
  )

## Install

install(TARGETS qpid-proton-cpp
  EXPORT  proton
  ARCHIVE DESTINATION ${LIB_INSTALL_DIR}
  LIBRARY DESTINATION ${LIB_INSTALL_DIR})

# Install windows qpid-proton-cpp pdb files
if (MSVC)
  install(FILES ${CMAKE_CURRENT_BINARY_DIR}/Debug/qpid-proton${CMAKE_DEBUG_POSTFIX}.pdb
    DESTINATION bin
    CONFIGURATIONS Debug
    OPTIONAL)
  install(FILES ${CMAKE_CURRENT_BINARY_DIR}/RelWithDebInfo/qpid-proton.pdb
    DESTINATION bin
    CONFIGURATIONS RelWithDebInfo
    OPTIONAL)
endif (MSVC)

# Install header files
file(GLOB headers "include/proton/*.hpp")
install (FILES ${headers} DESTINATION ${INCLUDE_INSTALL_DIR}/proton)

add_subdirectory(docs)
add_subdirectory(${CMAKE_SOURCE_DIR}/tests/tools/apps/cpp ${CMAKE_BINARY_DIR}/tests/tools/apps/cpp)

# Pkg config file
configure_file(
  ${CMAKE_CURRENT_SOURCE_DIR}/libqpid-proton-cpp.pc.in
  ${CMAKE_CURRENT_BINARY_DIR}/libqpid-proton-cpp.pc @ONLY)
install (FILES 
  ${CMAKE_CURRENT_BINARY_DIR}/libqpid-proton-cpp.pc
  DESTINATION ${LIB_INSTALL_DIR}/pkgconfig)

if (DEFINED CMAKE_IMPORT_LIBRARY_PREFIX)
set(PROTONCPPLIB ${CMAKE_IMPORT_LIBRARY_PREFIX}qpid-proton-cpp${CMAKE_IMPORT_LIBRARY_SUFFIX})
set(PROTONCPPLIBDEBUG ${CMAKE_IMPORT_LIBRARY_PREFIX}qpid-proton-cpp${CMAKE_DEBUG_POSTFIX}${CMAKE_IMPORT_LIBRARY_SUFFIX})
else ()
set(PROTONCPPLIB ${CMAKE_SHARED_LIBRARY_PREFIX}qpid-proton-cpp${CMAKE_SHARED_LIBRARY_SUFFIX})
set(PROTONCPPLIBDEBUG ${CMAKE_SHARED_LIBRARY_PREFIX}qpid-proton-cpp${CMAKE_DEBUG_POSTFIX}${CMAKE_SHARED_LIBRARY_SUFFIX})
endif ()

configure_file(
  ${CMAKE_CURRENT_SOURCE_DIR}/ProtonCppConfig.cmake.in
  ${CMAKE_CURRENT_BINARY_DIR}/ProtonCppConfig.cmake @ONLY)
configure_file(
  ${CMAKE_CURRENT_SOURCE_DIR}/ProtonCppConfigVersion.cmake.in
  ${CMAKE_CURRENT_BINARY_DIR}/ProtonCppConfigVersion.cmake @ONLY)
install (FILES
  ${CMAKE_CURRENT_BINARY_DIR}/ProtonCppConfig.cmake
  ${CMAKE_CURRENT_BINARY_DIR}/ProtonCppConfigVersion.cmake
  DESTINATION ${LIB_INSTALL_DIR}/cmake/ProtonCpp)

## Test
if (ENABLE_VALGRIND AND VALGRIND_EXE)
  set(memcheck-cmd ${VALGRIND_EXE} --error-exitcode=1 --quiet --leak-check=full --trace-children=yes)
endif ()

macro(add_cpp_test test)
  add_executable (${test} src/${test}.cpp)
  target_link_libraries (${test} qpid-proton qpid-proton-cpp)
  if (CMAKE_SYSTEM_NAME STREQUAL Windows)
    add_test (NAME cpp_${test}
      COMMAND ${env_py}
      "PATH=$<TARGET_FILE_DIR:qpid-proton>"
      $<TARGET_FILE:${test}> ${ARGN})
  else ()
    add_test (NAME cpp_${test} COMMAND ${memcheck-cmd} ${CMAKE_CURRENT_BINARY_DIR}/${test} ${ARGN})
  endif ()
endmacro(add_cpp_test)

add_cpp_test(interop_test ${CMAKE_SOURCE_DIR}/tests)
add_cpp_test(message_test)
add_cpp_test(encode_decode_test)
