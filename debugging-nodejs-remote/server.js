require('@google-cloud/debug-agent').start({serviceContext: {enableCanary: false}});
const {URL} = require("url");
const {createServer} = require("http");

// simple HTTP server
createServer((req, res) => {
  const url = new URL(req.url, "http://localhost");

  res.writeHead(200, {'Content-Type': 'text/plain'});
  if (url.searchParams.has("a")) {
    const str = url.searchParams.get("a");
    res.write(`Parsed to ${atoi(str)}\n`);
  } else {
    res.write('Hello!\nSet query parameter a to invoke atoi.\n');
  }
  res.end();
}).listen(process.env.PORT || 8080);


function atoi(str) {
    let result = 0;
    let is_negative = false;
    let index = 0;
    
    while (index < str.length && isspace(str[index])) {
        index++
    }
    
    if (index < str.length) {
        is_negative = str[index] === '-';
        if (str[index] === '-' || str[index] === '+') {
            index++;
        }
    }

	while (index < str.length && isdigit(str[index]))
	{
		result *= 10;
		result -= parsedigit(str[index]);
		index++;
	}

    return is_negative ? result : -result;
}

function isspace(c) {
	if (c == ' ')
		return (true);
	if (c == '\t')
		return (true);
	if (c == '\v')
		return (true);
	if (c == '\f')
		return (true);
	if (c == '\r')
		return (true);
	return (false);
}

function isdigit(c) {
    return c.charCodeAt(0) >= '0'.charCodeAt(0) && c.charCodeAt(0) <= '9'.charCodeAt(0);
}

function parsedigit(c) {
    return c.charCodeAt(0) - '1'.charCodeAt(0);
}