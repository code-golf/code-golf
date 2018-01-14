package terminal

import (
	"reflect"
	"testing"
)

var errorCases = []struct {
	name     string
	input    string
	expected string
}{
	{
		`sequence does not begin with 1337;File= or 1338;`,
		"foobar",
		`expected sequence to start with 1337;File=, 1338; or 1339;, got "foobar" instead`,
	}, {
		`sequence beginning error is cropped`,
		"123456789012345678901234567890123456789012345678901234567890",
		`expected sequence to start with 1337;File=, 1338; or 1339;, got "1234567890" instead`,
	}, {
		`1337: sequence does not have a content part and a arguments part`,
		"1337;File=foobar",
		`expected sequence to have one arguments part and one content part, got 1 part(s)`,
	}, {
		`1337: sequence has too many parts`,
		"1337;File=foobar:baz:foo",
		`expected sequence to have one arguments part and one content part, got 3 part(s)`,
	}, {
		`1337: content part is not valid Base64`,
		"1337;File=foobar:!!!!!",
		`expected content part to be valid Base64`,
	}, {
		`1337: image name is missing`,
		"1337;File=foobar:AA==",
		`name= argument not supplied, required to determine content type`,
	}, {
		`1337: invalid base64 encoding`,
		"1337;File=name=foo.baz:AA==",
		`name= value of "foo.baz" is not valid base64`,
	}, {
		`1337: can't determine content type`,
		"1337;File=name=" + base64Encode("foo.baz") + ":AA==",
		`can't determine content type for "foo.baz"`,
	}, {
		`1337: no image content`,
		"1337;File=name=foo.jpg:",
		`image content missing`,
	}, {
		`1338: url missing`,
		"1338;",
		`url= argument not supplied`,
	},
}

var validCases = []struct {
	name     string
	input    string
	expected *element
}{
	{
		`1337: image with name, content & inline`,
		`1337;File=name=Zm9vLmdpZg==;inline=1:AA==`,
		&element{url: "foo.gif", content: "AA==", contentType: "image/gif", elementType: ELEMENT_ITERM_IMAGE},
	}, {
		`1337: image without inline=1 does not render`,
		`1337;File=name=Zm9vLmdpZg==:AA==`,
		nil,
	}, {
		`1337: adapts content type based on image name`,
		`1337;File=name=` + base64Encode("foo.jpg") + `;inline=1:AA==`,
		&element{url: "foo.jpg", content: "AA==", contentType: "image/jpeg", elementType: ELEMENT_ITERM_IMAGE},
	}, {
		`1337: handles width & height`,
		`1337;File=name=Zm9vLmdpZg==;width=100%;height=50px;inline=1:AA==`,
		&element{url: "foo.gif", content: "AA==", contentType: "image/gif", width: "100%", height: "50px", elementType: ELEMENT_ITERM_IMAGE},
	}, {
		`1337: protects against XSS in image name, width & height by stripping brackets & quotes`,
		`1337;File=name=` + base64Encode(`foo".gif`) + `;width="100%;height='50px>;inline=1:AA==`,
		&element{url: "foo.gif", content: "AA==", contentType: "image/gif", width: "100%", height: "50px", elementType: ELEMENT_ITERM_IMAGE},
	}, {
		`1337: converts width & height without percent or px to em`,
		`1337;File=name=Zm9vLmdpZg==;width=1;height=5;inline=1:AA==`,
		&element{url: "foo.gif", content: "AA==", contentType: "image/gif", width: "1em", height: "5em", elementType: ELEMENT_ITERM_IMAGE},
	}, {
		`1337: malfored arguments are silently ignored`,
		`1337;File=name=Zm9vLmdpZg==;inline=1;sdfsdfs;====ddd;herp=derps:AA==`,
		&element{url: "foo.gif", content: "AA==", contentType: "image/gif", elementType: ELEMENT_ITERM_IMAGE},
	}, {
		`1338: image with filename`,
		"1338;url=tmp/foo.gif",
		&element{url: "tmp/foo.gif", elementType: ELEMENT_IMAGE},
	}, {
		`1338: image with filename containing an escaped ;`,
		"1338;url=tmp/foo\\;bar.gif",
		&element{url: "tmp/foo;bar.gif", elementType: ELEMENT_IMAGE},
	}, {
		`1338: image with filename, width, height & alt tag`,
		"1338;url=foo.gif;width=50px;height=50px;alt=foo gif",
		&element{url: "foo.gif", width: "50px", height: "50px", alt: "foo gif", elementType: ELEMENT_IMAGE},
	}, {
		`1339: link with url only`,
		"1339;url=foo.gif",
		&element{url: "foo.gif", elementType: ELEMENT_LINK},
	}, {
		`1339: link with url and content`,
		"1339;url=foo.gif;content=bar",
		&element{url: "foo.gif", content: "bar", elementType: ELEMENT_LINK},
	},
}

func TestErrorCases(t *testing.T) {
	for _, c := range errorCases {
		elem, err := parseElementSequence(c.input)
		if elem != nil {
			t.Errorf("%s\ninput\t\t%q\nexpected no image, received %+v", c.name, c.input, elem)
		} else if err.Error() != c.expected {
			t.Errorf("%s\ninput\t\t%q\nreceived\t%q\nexpected\t%q", c.name, c.input, err.Error(), c.expected)
		}
	}
}

func TestElementCases(t *testing.T) {
	for _, c := range validCases {
		elem, err := parseElementSequence(c.input)
		if err != nil {
			t.Errorf("%s\ninput\t\t%q\nexpected no error, received %s", c.name, c.input, err.Error())
		} else if !reflect.DeepEqual(elem, c.expected) {
			t.Errorf("%s\ninput\t\t%q\nreceived\t%+v\nexpected\t%+v", c.name, c.input, elem, c.expected)
		}
	}
}
