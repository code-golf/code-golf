# Creating holes

## Open a GitHub Issue

Create an issue and describe the hole

## Add hole description to `holes.toml`

Each hole needs to have its section in the `/config/holes.toml` file.
Each section needs to define the following fields:

- `category` - One of 'Art', 'Computing', 'Gaming', 'Mathematics', 'Sequence', 'Transform'.
- `preamble` - HTML description of the hole. May include go templating marks. The root `.` is the data from the `data` field.
- `synopsis` - Short overview of the hole.

In addition, each section may define the following fields:

- `data` - JSON describing the data to be copied upon pressing "Copy as JSON". Relevant mainly for holes in the "Transform" category.
- `links` - List of tables with keys `name` and `url`.
- `experiment` - ID of the issue that suggested this hole. All experimental holes need to link to an issue so that the community can vote and suggest improvements.
- `variants` - List of names of the holes that are variants of this hole, including itself.
- `case-fold` - A flag indicating that the output should be checked case insensitively.
- `item-delimiter` - If set, indicates that the output should be understood as a collection of items that can appear in any order and that the items should be delimited by the provided token.
- `multiset-delimiter` - If set, treats output as fixed order collection of multiple multisets, each separated by `item-delimiter`. Otherwise, output is treated as a single multiset.

Example:

```toml
['Hole Name']
category = 'hole category'
experiment = github-issue-number
links = [
    { name = 'Some useful link', url = '//someusefullink' },
    { name = 'Another useful link',    url = '//www.anotherusefullink' },
]
preamble = '''
<p>Hole description and instructions</p>
'''
```

## Hole inputs and expected answers

Depending on the new hole design, the answers could be hardcoded (99 bottles, 12 days, etc...) or computed on the fly with a function (sudoku, intersection, etc...).

### Holes with hardcoded solutions

If the hole has no inputs, the expected output needs to be placed in `/hole/answers/hole-name.txt`. The filename must be the [URLized](https://github.com/code-golf/code-golf/blob/master/config/config.go#L13) version of the hole name in the `holes.toml`.

### Holes with computed solutions

For computed solutions, a case switch needs to be added to `/hole/play.go`. The value needs to match the URLized version of the hole name. When a case matches, a function defined in its own file in `/holes/` is called.
This function must return data for at least one test run. Data for each run takes form of a list of inputs and a single string containing the expected output.

`/hole/play.go`

```go
	switch holeID {
	case "hole-name":
		scores = holeName()
```

`/hole/hole-name.go`

```go
func holeName() []Run {
	// Create Args and Answer

	return []Run{{Args: args, Answer: out}}
}
```

### Tips for randomizing the tests

The point of multiple random tests is to avoid cheating solutions that only implement some part of the problem.
Therefore:

- There should be a fixed amount of tests.
- If the set of possible outputs (or some obvious groups of the outputs) is finite (and low), each case should be represented.
- The set of all possible inputs must not be included in a single test - use the `outputMultirunTests` helper function to split the list of tests into two test runs.
- Pairs of inputs that are (in some sense) close but their corresponding outputs are different (or even far) should be included.
- The order of the tests should be randomized.

## Final steps

- Test the hole locally
- Commit and open a Pull Request
- Announce the hole on Discord and ask for feedback
