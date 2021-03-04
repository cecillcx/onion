module main

go 1.13

require (
	controller v0.0.0
	onion v0.0.0
)

replace controller v0.0.0 => ./controller

replace onion v0.0.0 => ./onion
