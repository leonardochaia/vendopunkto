# VendoPunkto Design Notes

## Architecture

VendoPunkto is a three tier application, where an HTTP server is provided for
multiple frontend clients to consume. Persistance is done with PostgreSQL,
although technically other DB implementations could be supported.

VendoPunkto server hosts a public and internal API on different ports so that
it's easier for operation security. You can port-forward the public port and
keep the internal port secure on your LAN.

The internal API may be consumed by plugins or integrating apps, whereas the public
API is mainly used to display invoices by integrating apps, or to send an invoice
publicly with a customer, by say, email or social media.

### Plugin Architecture

VendoPunkto-to-plugin communication is done via HTTP.
This means that plugins are language-agnostic and easier to develop and build.
There're no fixed dependencies for a plugin, one could implement the HTTP API
"interface" on any programming language that supports listening for HTTP.

A `plugin` package is provided by VendoPunkto to facilitate the plugin creations;
currently it's only developed for Golang.

Communication with plugins is done one way: VendoPunkto-to-plugin only. This
simplifies things and allows for better resilience.

`wallet` plugins have two main responsibilities: generating addresses and
confirming payments.

VendoPunkto will poll `wallet` plugins regularly asking for latest incoming
payments. The main advantage of using a polling system is resilience:
If the host becomes unavailable for some time, when it starts up again it will
request all incoming transactions after the last known block height to all
wallet plugins. This means that if VendoPunkto is offline, payments will still
be received and processed once it goes back online.

If the plugin becomes unavailable, VendoPunkto will keep polling. When the plugin
becomes online again, VendoPunkto will request all transfers after the last known
block height, which should be all payments since the plugin went offline,
ensuring no payments get lost.

This also ensures that it is the host who is in control of the throughput
it can handle. The host may adjust the polling interval depending on the current 
load. This may also help for scaling in the future.

Polling is a tradeoff. Under high load, it's more efficient for the host to poll
the plugin, than the plugin calling the host every time a new payment is received.

## Public API vs Internal API

VendoPunkto's public API is intended to be exposed to the internet when the
operator requires publishing payments. When using VendoPunkto on-site, nothing
is exposed to the internet, no port-forwards are required.

If you'd like to share invoices publicly, however, you will need to port-forward
the public API's port.

The internal API must not be exposed to the public internet, ideally should not
be exposed to LAN either. In small setups, only host level networking may be needed.
