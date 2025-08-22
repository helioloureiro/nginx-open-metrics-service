% nginx-openmetrics(1) | User Commands
%
% "August 22 2025"

[comment]: # The lines above form a Pandoc metadata block. They must be
[comment]: # the first ones in the file.
[comment]: # See https://pandoc.org/MANUAL.html#metadata-blocks for details.

[comment]: # pandoc -s -f markdown -t man package.md -o package.1
[comment]: # 
[comment]: # A manual page package.1 will be generated. You may view the
[comment]: # manual page with: nroff -man package.1 | less. A typical entry
[comment]: # in a Makefile or Makefile.am is:
[comment]: # 
[comment]: # package.1: package.md
[comment]: #         pandoc --standalone --from=markdown --to=man $< --output=$@
[comment]: # 
[comment]: # The pandoc binary is found in the pandoc package. Please remember
[comment]: # that if you create the nroff version in one of the debian/rules
[comment]: # file targets, such as build, you will need to include pandoc in
[comment]: # your Build-Depends control field.

[comment]: # lowdown is a low dependency, lightweight alternative to
[comment]: # pandoc as a markdown to manpage translator. Use with:
[comment]: # 
[comment]: # package.1: package.md
[comment]: #         lowdown -s -Tman -o $@ $<
[comment]: # 
[comment]: # And add lowdown to the Build-Depends control field.

# NAME

nginx-openmetrics - a service to expose nginx statistics as open metrics

# SYNOPSIS

**nginx-openmetrics** [**\-\-loglevel=**<level>] [**\-\-port=**<port>]
                [**\-\-service=**<service>]

**nginx-openmetrics** [{**\-\-help**} | {**\-\-version**}]

# DESCRIPTION

This manual page documents briefly the **nginx-openmetrics** and **service**, 
**loglevel** and **port** commands.

**nginx-openmetrics** fetches the data from nginx server where the statistics
are enabled.  It updates every **15 seconds**.  Data is exposed using open metrics 
as counters and gauges.

The following metrics are displayed:

 - **nginx_active_connections**: the number of active connections (gauge).
 - **nginx_reading_connections**: the number of active reading connections (gauge).
 - **nginx_waiting_connections**: the number of waiting connections (gauge).
 - **nginx_writing_connections**: the number of active writing connections (gauge).
 - **nginx_server_accepts_total**: the total number of server accepted connections (counter).
 - **nginx_server_handled_total**: the total number of server handled connections (counter).
 - **nginx_server_requests_total**: the total number of server requests (counter).


# OPTIONS

The program follows the usual with long options starting with two dashes ('-'). 
A summary of options is included below. For a complete description, 
see the **info**(1) files.

**\-\-service=**<service>
:   The endpoint where nginx is serving the statistics.
    Something like "http://localhost:8080/api".
    Always include "http/https" on the uri.

**\-\-port=**<port>
:   The port to listen to.  Default is 9090.

**\-\-loglevel=**<level>
:   The logging messages.  Default is "info".
    You can select from the values: info, debug, warn, error and fatal.

**\-\-help**
:   Show summary of options.

**\-\-version**
:   Show version of program.



# DIAGNOSTICS

The following diagnostics may be issued on stderr:

Bad configuration file. Exiting.
:   The configuration file seems to contain a broken configuration
    line. Use the **\-\-loglevel=**__debug__ option, to get more info.

**nginx-openmetrics** provides some return codes, that can be used in scripts:

    Code Diagnostic
    0 Program exited successfully.
    1 The configuration file seems to be broken.

# BUGS

The program is currently limited to work without labels.

The upstream BTS can be found at:
https://github.com/helioloureiro/nginx-open-metrics-service/issues.

# SEE ALSO

**nginx**(8)


# AUTHOR

Helio Loureiro <helio@loureiro.eng.br>

# COPYRIGHT

Copyright © 2025 Helio Loureiro

This manual page was written for the Ubuntu system (and may be used by
others).

Permission is granted to copy, distribute and/or modify this document under
the terms of the GNU General Public License, Version 2 published by the 
Free Software Foundation.

On Debian systems, the complete text of the GNU General Public License
can be found in /usr/share/common-licenses/GPL.

[comment]: #  Local Variables:
[comment]: #  mode: markdown
[comment]: #  End:
