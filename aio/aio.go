package aio

// package aio sounds a little strange in golang. It seems that all golang app should be concurrent.
// but for some complicated cases, such as MMO games, a single logic thread is more practical-- avoid
// concurrency control in logic which makes logic simple.
// so gnet provides aio package to support async calls.
