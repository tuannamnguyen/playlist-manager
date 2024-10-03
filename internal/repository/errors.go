package repository

import "fmt"

type rowScanError struct {
	err error
}

func (r *rowScanError) Error() string {
	return fmt.Sprintf("db row scan: %s", r.err.Error())
}

type queryError struct {
	err error
}

func (q *queryError) Error() string {
	return fmt.Sprintf("db query: %s", q.err.Error())
}

type structScanError struct {
	err error
}

func (q *structScanError) Error() string {
	return fmt.Sprintf("db query row: %s", q.err.Error())
}

type execError struct {
	err error
}

func (e *execError) Error() string {
	return fmt.Sprintf("db exec: %s", e.err.Error())
}

type selectError struct {
	err error
}

func (s *selectError) Error() string {
	return fmt.Sprintf("db select: %s", s.err.Error())
}

type beginTransactionError struct {
	err error
}

func (b *beginTransactionError) Error() string {
	return fmt.Sprintf("begin db transaction: %s", b.err.Error())
}

type transactionCommitError struct {
	err error
}

func (t *transactionCommitError) Error() string {
	return fmt.Sprintf("commit db transaction: %s", t.err.Error())
}

type prepareInQueryError struct {
	err error
}

func (p *prepareInQueryError) Error() string {
	return fmt.Sprintf("prepare IN query: %s", p.err.Error())
}

type requestMarshalError struct {
	err error
}

func (r *requestMarshalError) Error() string {
	return fmt.Sprintf("request marshall: %s", r.err.Error())
}

type responseDecodeError struct {
	err error
}

func (r *responseDecodeError) Error() string {
	return fmt.Sprintf("response decode: %s", r.err.Error())
}

type putObjectError struct {
	err error
}

func (p *putObjectError) Error() string {
	return fmt.Sprintf("put object in storage: %s", p.err.Error())
}
