# flow
Flow is a framework for Dagger for writing flexible CI pipelines in Go that have consistent behavior when ran locally or in a CI server.

Write your pipeline once, run it locally and produce the config for your CI provider from the same code.
# Status

This is still under development. Expect breaking changes and incomplete features.
# Why?

With Flow you can:

    - Run pipelines locally for testing using Dagger.
    - Generate configurations for existing CI providers.
    - Use tools like delve to debug your pipelines.
    - Use Go features to make complex pipelines easier to develop and maintain.

# Why not only use Dagger?

Flow internally uses Dagger to accomplish consistency. The promise of running the same thing locally that's ran in a CI service could not happen without it.

Flow adds a few features on top of Dagger, like:

    - Executing an anonymous function instead of a command.
    - Generating a CI configuration from your pipeline code.

# Running Locally / testing
