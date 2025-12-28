package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderWith1NormalFieldValPair(t *testing.T) {

	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("HOst: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	host,_:=headers.Get("hOsT")
	assert.Equal(t, "localhost:42069", host)
	assert.Equal(t, 25, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
func TestHeaderWithMultiNormalFieldValPair(t *testing.T) {

	// Test: Multi-field values
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n Foo: Bar \r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	host,_:=headers.Get("hosT")
	foo,_:=headers.Get("FOO")
	assert.Equal(t, "localhost:42069",host)
	assert.Equal(t, "Bar", foo)
	//assert.Equal(t, 23, n)
	assert.True(t, done)
}

func TestEmpty(t *testing.T) {

	// Test: Empty
	headers := NewHeaders()
	data := []byte("\r\n\r\n       Host: localhost:42069       \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.True(t, done)

}

func TestMalformed(t *testing.T) {
	// Test: Malformed
	headers := NewHeaders()
	data := []byte("Host     \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, // malformed requestLine
// min requestLine
// max requestLine
err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestMutlipleHeaders(t *testing.T) {

	// Test: Multi-field values
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n Host: Bar \r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	host ,_:= headers.Get("hosT")
	assert.Equal(t, "localhost:42069, Bar", host)
	//assert.Equal(t, 23, n)
	assert.True(t, done)
}





// malformed requestLine
// min requestLine
// max requestLine
