# Serial buffer flood
**move the commands to the macro command**
Gcode M810 til M819 can be used to store a macro example: M815 G0 X0 Y0|G0 Z10|M300 S440 P50
commands are seperated by the pipe | 
calling M815 without arguments runs the macro
**wait for completion of execution before sending new commands**
Using M118 to have the printer send a serial command, we can detect if a set of instructions has been executed before sending in the next batch