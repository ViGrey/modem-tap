# modem-tap

An audio Bell103 300 baud modem wiretap synthesizer

**_Modem-Tap is created by Vi Grey (https://vigrey.com) <vi@vigrey.com> and is licensed under the BSD 2-Clause License.  Read LICENSE for more license text._**

#### Description:
**THIS PROJECT IS A PERSONAL PROJECT.  I CANNOT PROMISE UPDATES FOR ANY FEATURES OR FIXES OTHER THAN FEATURES AND FIXES I WISH TO ADD.  QUITE A BIT OF THE ERROR HANDLING IS SIMPLY JUST ME KNOWING WHAT I SHOULD AND SHOULD NOT DO WITH THIS PROGRAM TO KEEP IT FROM CRASHING AND/OR RUNNING INTO ERRORS.**

Modem-Tap is a Bell103 300 baud modem noise synthesizer that emulates a 300bps dial-up connection and produces the incoming and outgoing connection data tones.

This project is called Modem-Tap because it creates the sounds that would be heard if the phone line handling transferring the 300 baud dial-up connection was wiretapped.

Connections through Modem-Tap are throttled to 300bps and uses 8-N-1 ascii encoding.

#### Platforms:
- GNU/Linux

#### Build Dependencies:
- gb
- Go >= 1.8
- Portaudio-Dev >= 19

#### Dependencies:
- telnet

#### Install:

    $ make
    $ sudo make install

#### Uninstall:
    $ sudo make uninstall

#### Usage:
    $ modem-tap [ OPTIONS ]... [ LISTENING PORT ]

    Options:

      -h, --help      Print help (this message) and exit
      -q, --quiet     Does not play connection sounds from the speaker
      -w, --wav       Saves connection sounds to a WAV file

    LISTENING PORT is 2600 by default

#### Startup:
    Upon startup, Modem-Tap will ask for a server address.  It will then ask for a server port.  After you supply both the address and the port, Modem-Tap will then listen at localhost on the LISTENING PORT you specified, which is 2600 if no port is specified.  You can then telnet into localhost:2600 and Modem-Tap will make a telnet connection out to the server address on the server port.  This connection will create the Bell 103 modulation frequency sounds for the incoming and outgoing internet traffic, which can be played through speakers and/or recorded to a WAV file.

    After the connection is closed, if Modem-Tap has not been closed with ctrl-c, it should ask for another server address and then another server port.

#### WAV Files:
    WAV files created by Modem-Tap will have the filename syntax of serveraddress-serverport-YYYYMMDDhhmmss.wav, for instance, a connection to vigrey.com on port 80 could produce the file name vigrey.com-80-20170717011252.wav.  This file will be a single channel 44100Hz 16-bit PCM WAV file.

    A single WAV file is created for each connection to a server from Modem-Tap.  Multiple connections can be created on a single instance of running Modem-Tap.
