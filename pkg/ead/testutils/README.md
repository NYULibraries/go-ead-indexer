Keeping _testutils/_ in this package for now, but there's a fairly good chance
this will need to be a shared package at top level, so for now, only include
test utilities that can be shared with any other package.  All test utilities
specific to `ead` package should go in _ead_test.go_ and or _ead/util/_.

