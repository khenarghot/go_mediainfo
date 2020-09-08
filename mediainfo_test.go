package mediainfo_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/khenarghot/go_mediainfo"
	"github.com/stretchr/testify/assert"
)

const (
	ogg       = "testdata/test.ogg"
	mp3       = "testdata/test.mp3"
	nonExists = "testdata/non_exists.ogg"
)

func TestOpenWithOgg(t *testing.T) {
	assert := assert.New(t)
	mi := mediainfo.NewMediaInfo()
	assert.NoError(mi.OpenFile(ogg))
}

func TestOpenWithMp3(t *testing.T) {
	assert := assert.New(t)
	mi := mediainfo.NewMediaInfo()
	assert.NoError(mi.OpenFile(mp3))
}

func TestOpenWithUnexistsFile(t *testing.T) {
	assert := assert.New(t)
	mi := mediainfo.NewMediaInfo()
	assert.Error(mi.OpenFile(nonExists))
}

func TestOpenMemoryWithOgg(t *testing.T) {
	assert := assert.New(t)
	mi := mediainfo.NewMediaInfo()
	f, _ := os.Open(ogg)
	bytes, _ := ioutil.ReadAll(f)

	assert.NoError(mi.OpenMemory(bytes))
}

func TestOpenMemoryWithMp3(t *testing.T) {
	assert := assert.New(t)
	mi := mediainfo.NewMediaInfo()
	f, _ := os.Open(mp3)
	bytes, _ := ioutil.ReadAll(f)

	assert.NoError(mi.OpenMemory(bytes))
}

func TestOpenMemoryWithEmptyArray(t *testing.T) {
	assert := assert.New(t)
	mi := mediainfo.NewMediaInfo()
	assert.Error(mi.OpenMemory([]byte{}))

}

func TestInformWithOgg(t *testing.T) {
	mi := mediainfo.NewMediaInfo()
	mi.OpenFile(ogg)

	if len(mi.Inform()) < 10 {
		t.Fail()
	}
}

func TestInformWithMp3(t *testing.T) {
	mi := mediainfo.NewMediaInfo()
	mi.OpenFile(mp3)

	if len(mi.Inform()) < 10 {
		t.Fail()
	}
}

func TestAvailableParametersWithOgg(t *testing.T) {
	mi := mediainfo.NewMediaInfo()
	mi.OpenFile(ogg)

	if len(mi.AvailableParameters()) < 10 {
		t.Fail()
	}
}

func TestAvailableParametersWithMp3(t *testing.T) {
	mi := mediainfo.NewMediaInfo()
	mi.OpenFile(mp3)

	if len(mi.AvailableParameters()) < 10 {
		t.Fail()
	}
}

func TestDurationWithOgg(t *testing.T) {
	assert := assert.New(t)
	mi := mediainfo.NewMediaInfo()
	mi.OpenFile(ogg)

	assert.Equal(3494, mi.Duration())
}

func TestDurationWithMp3(t *testing.T) {
	assert := assert.New(t)
	mi := mediainfo.NewMediaInfo()
	mi.OpenFile(mp3)

	assert.Equal(87771, mi.Duration())
}

func TestCodecWithOgg(t *testing.T) {
	assert := assert.New(t)
	mi := mediainfo.NewMediaInfo()
	mi.OpenFile(ogg)
	assert.Equal("Ogg", mi.Codec())
}

func TestCodecWithMp3(t *testing.T) {
	assert := assert.New(t)
	mi := mediainfo.NewMediaInfo()
	mi.OpenFile(mp3)
	assert.Equal("MPEG Audio", mi.Codec())
}

func TestFormatWithOgg(t *testing.T) {
	assert := assert.New(t)
	mi := mediainfo.NewMediaInfo()
	mi.OpenFile(ogg)
	assert.Equal("Ogg", mi.Format())
}

func TestFormatWithMp3(t *testing.T) {
	assert := assert.New(t)
	mi := mediainfo.NewMediaInfo()
	mi.OpenFile(mp3)
	assert.Equal("MPEG Audio", mi.Format())
}

//----------------------------------------------------------------------------------------------------------------------
func BenchmarkOpenAndDurationWithOgg(b *testing.B) {
	for n := 0; n < b.N; n++ {
		mi := mediainfo.NewMediaInfo()
		mi.OpenFile(ogg)

		mi.Duration()
	}
}

func BenchmarkOpenAndDurationWithMp3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		mi := mediainfo.NewMediaInfo()
		mi.OpenFile(mp3)

		mi.Duration()
	}
}

func BenchmarkOpenMemoryAndDurationWithOgg(b *testing.B) {
	for n := 0; n < b.N; n++ {
		mi := mediainfo.NewMediaInfo()
		f, _ := os.Open(ogg)
		bytes, _ := ioutil.ReadAll(f)

		mi.OpenMemory(bytes)
		mi.Duration()
	}
}

func BenchmarkOpenMemoryAndDurationWithMp3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		mi := mediainfo.NewMediaInfo()
		f, _ := os.Open(mp3)
		bytes, _ := ioutil.ReadAll(f)

		mi.OpenMemory(bytes)
		mi.Duration()
	}
}

//----------------------------------------------------------------------------------------------------------------------

func ExampleUsage() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	mi := mediainfo.NewMediaInfo()
	err = mi.OpenMemory(bytes)
	if err != nil {
		panic(err)
	}
	fmt.Println(mi.AvailableParameters()) // Print all supported params for Get
	fmt.Println(mi.Get("BitRate"))        // Print bitrate
}
