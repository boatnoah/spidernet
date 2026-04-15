package crawler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkExtractor(t *testing.T) {
	assert := assert.New(t)

	htmlString := `
		<button>
			<a href="https://www.boatnoah.com">
				press me 
			</a>
		</button>
		<button>
			<a href="http://www.boatnoah.com">
				press me 
			</a>
			<a href="http://www.boatnoah.com">
				this is a duplicate
			</a>
		</button>
		<button>
			<a href="https://github.com/boatnoah">
				press me 
			</a>
		</button>
		<button>
			<a href="https://www.google.com">
				press me 
			</a>
		</button>

		<ul>
		  <li><a href="https://wikipedia.org">Wikipedia</a></li>
		  <li><a href="https://www.w3schools.com">W3Schools</a></li>
		  <li><a href="mailto:m.bluth@example.com">Email</a></li>
		</ul>
		
		<div> 
			<a href="https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/a" target="_blank">
				mozilla
			</a>
		</div>
	`
	expected := []string{
		"https://www.boatnoah.com",
		"http://www.boatnoah.com",
		"https://github.com/boatnoah",
		"https://www.google.com",
		"https://wikipedia.org",
		"https://www.w3schools.com",
		"https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/a",
	}

	foundLinks, err := extractLinks([]byte(htmlString))
	if err != nil {
		assert.Error(err)
		return
	}

	assert.Len(foundLinks, len(expected))

	for i := range expected {
		assert.Equal(expected[i], foundLinks[i])
	}
}
