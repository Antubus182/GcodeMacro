# MVP
**28-05-2024**
Minimal Viable Product is produced
In it's simplest form does the program read in a predefinded json file and sends the command over serial to the printer

# WIP
**18-01-2023**
Project is currently under development and not yet functional

# GcodeMacro
A simple program that reads Gcode commands from a json file and sends them to marlin firmware

# Description
This program can be used to send a simple sequence of Gcodes to for example a 3D printer for experimental/development purposes.
The program reads a macro.csv file in the root directory and sequentially sends them over serial to a device running marlin firmware. 
This can be usefull if you, for example want to repurpose a printer to function as a motion platform and want to quickly run a series of gcodes as a simple macro