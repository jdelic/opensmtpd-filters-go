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
instance to ``opensmtpd.NewFilter``.

.. code-block:: go

    import (
        "os"
        "log"
        "github.com/jdelic/opensmtpd-filters-go"
    )
    
    type FilterExample struct {}
    
    func (ex *FilterExample) LinkConnect(fw opensmtpd.FilterWrapper, 
        verb string, sh SessionHolder, sessionId string, params []string) {
        log.Println("link-connect received")
    } 
    
    func main() {
        log.SetOutput(os.Stderr)  // stderr is forwarded to log by opensmtpd
        myFilter := opensmtpd.NewFilter(&FilterExample{})
        opensmtpd.Run(myFilter)
    }
    

Provided interfaces
===================

Reporters
---------

