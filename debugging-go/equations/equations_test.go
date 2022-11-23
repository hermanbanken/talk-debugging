package main

import (
	"context"
	"dummy/calculator/job"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEquation(t *testing.T) {
	j, err := parseEquation(context.Background(), []byte("1 + 2 * 3"))
	assert.NoError(t, err)
	assert.Equal(t, job.Job{
		Op: job.Add,
		A:  job.Value(1),
		B: job.Job{
			Op: job.Multiply,
			A:  job.Value(2),
			B:  job.Value(3),
		},
	}, j)

	j, err = parseEquation(context.Background(), []byte("1 + 2 * 3 + 4"))
	assert.NoError(t, err)
	assert.Equal(t, job.Job{
		Op: job.Add,
		A:  job.Value(1),
		B: job.Job{
			Op: job.Add,
			A: job.Job{
				Op: job.Multiply,
				A:  job.Value(2),
				B:  job.Value(3),
			},
			B: job.Value(4),
		},
	}, j)

	j, err = parseEquation(context.Background(), []byte("1 + 2 * 3 / 4"))
	assert.NoError(t, err)
	assert.Equal(t, job.Job{
		Op: job.Add,
		A:  job.Value(1),
		B: job.Job{
			Op: job.Multiply,
			A:  job.Value(2),
			B: job.Job{
				Op: job.Divide,
				A:  job.Value(3),
				B:  job.Value(4),
			},
		},
	}, j)

}
