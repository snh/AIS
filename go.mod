module github.com/snh/AIS

go 1.17

require (
	github.com/andmarios/aislib v0.0.0-20190131232958-3a9a58899c39
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/tormol/AIS v0.0.0-00010101000000-000000000000
)

replace github.com/tormol/AIS => github.com/snh/AIS type-24
