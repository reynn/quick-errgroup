# Quick Error Group

Similar to [errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup) except when an error occurs it cancels the rest of the currently running routines and exits. Useful when there is a case when one routine can fail and interrupt proper application flow but you don't/can't wait for all of the routines to fail.
