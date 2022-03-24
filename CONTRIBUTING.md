# Holes
Each hole needs to have its section in the `/config/holes.toml` file.
Each section needs to define the following fields:
* `category` - One of 'Art', 'Computing', 'Gaming', 'Mathematics', 'Sequence', 'Transform'.
* `preamble` - HTML description of the hole. May include go templating marks. The root `.` is the data from the `data` field.

In addition, each section may define the following fields:
* `data` - JSON describing the data to be copied upon pressing "Copy as JSON". Relevant mainly for holes in the "Transform" category.
* `links` - List of tables with keys `name` and `url`.
* `experiment` - ID of the issue that suggested this hole or `-1` if the issue doesn't exist. Only defined for experimental holes.
* `variants` - List of names of the holes that are variants of this hole, including itself.

If the hole has no inputs, the expected output needs to be placed in `/hole/answers/hole-name.txt`. The filename must be the [URLized](https://github.com/code-golf/code-golf/blob/master/config/config.go#L13) version of the hole name in the `holes.toml`.

Otherwise, a case switch needs to be added to `/hole/play.go`. The value needs to match the URLized version of the hole name. When a case matches, a function defined in its own file in `/holes/` is called. This function must return randomized list of inputs and the corresponding output (the output takes a form of a single string).

## Tips for randomizing the tests
The point of multiple random tests is to avoid cheating solutions that only implement some part of the problem.
Therefore:
* There should be a fixed amount of tests.
* If the set of possible outputs (or some obvious groups of the outputs) is finite (and low), each case should be represented a random number of times but at least once.
* Pairs of inputs that are (in some sense) close but their corresponding outputs are different (or even far) should be included.
* The order of the tests should be randomized.
