# Creating holes

## Open a GitHub Issue

Create an issue and describe the hole

## Add hole description to `holes.toml`

Each hole needs to have its section in the `/config/data/holes.toml` file.
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
- `multiset-item-delimiter` - If set, indicates that the output should be understood as a collection of items that can appear in any order and that the items should be delimited by the provided token.
- `output-delimiter` - If set, treats output as fixed order collection of multiple outputs and treats each one individually.

Example:

```toml
['Hole Name']
category = 'hole category'
experiment = github-issue-number
links = [
    { name = 'Some useful link',    url = '//someusefullink' },
    { name = 'Another useful link', url = '//www.anotherusefullink' },
]
synopsis = 'Short overview of the hole'
preamble = '''
<p>Hole description and instructions</p>
'''
```

## Hole inputs and expected answers

Depending on the new hole design, the answers could be hardcoded (99 bottles, 12 days, etc...) or computed on the fly with a function (sudoku, intersection, etc...).

### Holes with hardcoded solutions

If the hole has no inputs, the expected output needs to be placed in `/config/hole-answers/hole-name.txt`. The filename must be the [URLized](https://github.com/code-golf/code-golf/blob/master/config/config.go#L13) version of the hole name in the `holes.toml`.

### Holes with computed solutions

For computed solutions, a function is defined and registered in its own file in `/holes/`.
This function must return data for at least one test run. Data for each run takes form of a list of inputs and a single string containing the expected output.

`/hole/hole-name.go`

```go
var _ = answerFunc("hole-id", func() []Answer {
    // Implement hole and create tests

    return outputTests(shuffle(tests))
}
```

### Custom judges
In addition to simply comparing the user output to the expected output (either from the static file or the generator) by equality, holes can employ custom judges. Given a run (list of inputs and a user output) a judge determines whether the output is correct or not. If it is, the judge must return the user output. Otherwise, it should return an output that is correct and is as close to the user output as possible. To write a custom judge, you may use the `perOutputJudge` or `oneOfPerOutputJudge` helper functions. A judge is registered in a similar way to the answer function:

```go
var _ = judge("hole-id", func(run Run) string {
    // Return a valid answer similar to `run.Stdout`

    return run.Answer.Answer;
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
