/*
Photobomb conducts workflow tests triggered by requests to its web server.

Load-testing or stress-testing an API can be imagined as an air strike, which
consists of:

    1. a Mission that describes when and what ordnance will be deployed
    2. a Squadron of planes that can simultaneously deploy their Arsenals
    3. an Arsenal on each plane consisting of Bombs and Missiles

The default mission (1) causes the configured squadron to simultaneously deploy
their configured arsenals every 5 seconds. Some ordnance will hit its target
and result in reportable damage (response time) sooner than others, and the
output logging will reflect this, as the reports will typically arrive "out of
order" as the concurrent weapon deployments finish.

Bombs are used for API transactions where all inputs can be known before
runtime, such as getting a list of all Submission Batches for your account,
all Contributions within a Submission Batch with a known ID, or the data for a
specific Release with a known ID. Bombs assume the state of objects associated
with your account and are trusted to hit their target. If they "miss" with a
404 or other 4xx, this will be reported as an error.

Missiles are used for API transactions that depend on the current state of your
account, such as deleting the most-recently-created Submission Batch. Their
initialization is more complex than that for Bombs because they need to be
provided a function that can be lazy-evaluated when Missile.Fire() occurs.
Unlike a Bomb, Missiles need to be "guided" in this fashion to ensure they
will hit their target. As with Bombs, 4xx responses are reported as errors.
*/
package main
