## CI

This folder contains the source for a Docker image suitable for building
the lambda functions in a CI environment.

It is the standard Go toolkit image with `make` and `zip` added to it.

If you mount or git clone the source tree into it, you can easily
compile the lambas by running `make deployables`.
