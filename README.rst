Golang implementation of the OpenSMTPD filter protocol
======================================================

**This library is in active development and even basic APIs might change any
second.**

This library tries to provide an easy way to implement OpenSMTPD filters and
reporters. Filters allow you to talk back to the OpenSMTPD daemon and change
protocol responses and even emails. Reporters give your program information
which allows you to record statistics, for example.

The basic principle of the library is to abstract away most of the ongoing
communication and common patterns from your program and allow you to focus on
purely the "business logic" of your filter. You do so by creating a class that
implements any of the interfaces listed below and then passing a pointer to an
instance to ``opensmtpd.NewFilter`` and subsequently ``opensmtpd.Run``.

.. code-block:: go

    import (
        "os"
        "log"
        "github.com/jdelic/opensmtpd-filters-go"
    )

    type FilterExample struct {}

    // implement the required opensmtpd.Filter interface
    func (ex *FilterExample) GetName() string {
        return "My example filter"
    }

    func (ex *FilterExample) LinkConnect(fw opensmtpd.FilterWrapper,
        ev opensmtpd.FilterEvent) {
        log.Println("link-connect received")
    }

    func main() {
        log.SetOutput(os.Stderr)  // stderr is forwarded to log by opensmtpd
        myFilter := opensmtpd.NewFilter(&FilterExample{})
        opensmtpd.Run(myFilter)
    }


Provided interfaces
===================

Every filter must implement the ``opensmtpd.Filter`` interface. That interface
is really just a placeholder for now. However, without it, Go makes it
impossible to enforce that the filter instance must be passed by reference and
as all filter methods should use a pointer receiver you would experience
problems otherwise as ``opensmtpd.Register`` would be unable to figure out
the interfaces your filter actually implements.

Reporters
---------

See `opensmtpd-filters-go/report_api_interfaces.go <reporters_>`__.

Filters
-------

See `opensmtpd-filters-go/filter_api_interfaces.go <filters_>`__.

EventResponder
--------------

This interface is handed back from ``FilterEvent.Responder()`` and handles the
communicating with the different versions of OpenSMTPD's filter API.

See `opensmtpd-filters-go/eventresponder.go <eventresponders_>`__.


.. _filters: https://github.com/jdelic/opensmtpd-filters-go/blob/master/filter_api_interfaces.go
.. _reporters: https://github.com/jdelic/opensmtpd-filters-go/blob/master/report_api_interfaces.go
.. _eventresponders: https://github.com/jdelic/opensmtpd-filters-go/blob/master/eventresponder.go
